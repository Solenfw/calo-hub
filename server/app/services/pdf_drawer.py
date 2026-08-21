from concurrent.futures import ThreadPoolExecutor
from io import BytesIO
from urllib.parse import urlparse

import requests
from PIL import Image
from reportlab.lib.colors import HexColor, black
from reportlab.lib.pagesizes import A4
from reportlab.lib.units import cm
from reportlab.lib.utils import ImageReader
from reportlab.pdfbase.pdfmetrics import stringWidth
from reportlab.pdfgen import canvas

# ==============================================================================
# CONFIGURATION & LAYOUT CONSTANTS
# ==============================================================================

# --- Page Setup ---
PAGE_W, PAGE_H = A4

MARGIN_L = 1.8 * cm
MARGIN_R = 1.8 * cm
MARGIN_TOP = 1.6 * cm
MARGIN_BOTTOM = 1.8 * cm
CONTENT_W = PAGE_W - MARGIN_L - MARGIN_R

# --- Column Positions & Widths (Detailed Catalog View) ---
COL_NO_X = MARGIN_L
COL_CODE_X = MARGIN_L + 1.0 * cm
COL_DESC_X = MARGIN_L + 3.6 * cm
COL_QTY_X = MARGIN_L + 10.8 * cm
COL_IMG_X = MARGIN_L + 12.6 * cm
COL_IMG_W = CONTENT_W - 12.6 * cm

DESC_COL_W = COL_QTY_X - COL_DESC_X - 0.3 * cm

# --- Typography & Row Heights ---
FONT_BODY, FONT_BODY_SIZE = "Helvetica", 10.5
LINE_HEIGHT = 14

IMG_MAX_W = COL_IMG_W - 0.2 * cm
IMG_MAX_H = 3.6 * cm

ROW_TOP_PAD = 0.55 * cm
ROW_BOTTOM_PAD = 0.55 * cm
MIN_ROW_CONTENT_H = 2.3 * cm

# --- Summary Table Configuration ---
TABLE_HEADER_H = 0.9 * cm
CELL_PAD_X = 0.15 * cm
CELL_PAD_Y = 0.25 * cm
TABLE_BORDER_COLOR = HexColor("#333333")
TABLE_BORDER_WIDTH = 0.6

SUMMARY_COL_BOUNDS = [
    MARGIN_L,
    COL_CODE_X,
    COL_DESC_X,
    COL_QTY_X + 5 * cm,
    PAGE_W - MARGIN_R,
]
SUMMARY_COL_LABELS = ["No", "code", "Description", "Qty"]
SUMMARY_DESC_COL_W = (COL_QTY_X + 5 * cm) - COL_DESC_X - 2 * CELL_PAD_X
SUMMARY_ROW_GAP = 0.35 * cm

# ==============================================================================
# IMAGE FETCHING & PREFETCHING UTILITIES
# ==============================================================================

# Persistent HTTP session for connection pooling and keep-alive optimization
_session = requests.Session()


def _fetch_bytes(url, timeout=6, retries=1):
    """
    Downloads raw bytes for a given URL with a basic retry budget.
    """
    for attempt in range(retries + 1):
        try:
            resp = _session.get(url, timeout=timeout)
            resp.raise_for_status()
            return resp.content
        except Exception as e:
            if attempt == retries:
                print(f"Failed to fetch {url}: {e}")
                return None


def load_image(path_or_url, raw_bytes=None):
    """
    Builds a ReportLab ImageReader instance from:
    - Pre-fetched raw bytes,
    - A local file path, or
    - A remote URL (as fallback if not prefetched).

    Returns None if loading fails.
    """
    if not path_or_url:
        return None
    try:
        if raw_bytes is not None:
            img = Image.open(BytesIO(raw_bytes))
        else:
            parsed = urlparse(path_or_url)
            if parsed.scheme in ("http", "https"):
                data = _fetch_bytes(path_or_url)
                if data is None:
                    return None
                img = Image.open(BytesIO(data))
            else:
                img = Image.open(path_or_url)
        img.load()
    except Exception as e:
        print(f"Failed to load image {path_or_url}: {e}")
        return None

    if img.mode not in ("RGB", "RGBA", "L"):
        img = img.convert("RGB")
    return ImageReader(img)


def prefetch_images(products, max_workers=10):
    """
    Concurrently downloads all remote product images using ThreadPoolExecutor.
    Returns a mapping dict of {url: ImageReader or None}.
    """
    urls = sorted(
        {
            p.get("image")
            for p in products
            if p.get("image") and urlparse(p["image"]).scheme in ("http", "https")
        }
    )
    if not urls:
        return {}

    with ThreadPoolExecutor(max_workers=max_workers) as ex:
        raw_results = list(ex.map(_fetch_bytes, urls))

    cache = {}
    for url, raw in zip(urls, raw_results):
        cache[url] = load_image(url, raw_bytes=raw) if raw is not None else None
    return cache


# ==============================================================================
# TEXT UTILITIES
# ==============================================================================


def wrap_text_by_width(text, font_name, font_size, max_width):
    """
    Word-wraps text based on exact string pixel/point width rather than standard character count.
    """
    if not text:
        return [""]
    words, lines, current = text.split(), [], ""
    for word in words:
        candidate = f"{current} {word}".strip()
        if stringWidth(candidate, font_name, font_size) <= max_width:
            current = candidate
        else:
            if current:
                lines.append(current)
            current = word
    if current:
        lines.append(current)
    return lines or [""]


# ==============================================================================
# PDF DRAWING FUNCTIONS
# ==============================================================================


def draw_header(c, company_name, category_title):
    """
    Draws the main page masthead and column titles.
    Returns the y-coordinate where the first item row should start.
    """
    y = PAGE_H - MARGIN_TOP

    c.setFillColor(black)
    c.setFont("Helvetica-Bold", 15)
    c.drawString(MARGIN_L, y, company_name)
    y -= 0.35 * cm
    c.setLineWidth(1)
    c.line(MARGIN_L, y, PAGE_W - MARGIN_R, y)

    y -= 0.65 * cm
    c.setFont("Helvetica-Bold", 12.5)
    c.drawString(MARGIN_L, y, category_title)

    y -= 0.85 * cm
    if category_title != "Summary":
        c.setFont("Helvetica", 11)
        c.drawString(COL_NO_X, y, "No.")
        c.drawString(COL_CODE_X, y, "Code")
        c.drawString(COL_DESC_X, y, "Description")
        c.drawString(COL_QTY_X, y, "Quantity")
        y -= 0.25 * cm
        c.line(MARGIN_L, y, PAGE_W - MARGIN_R, y)

    return y - ROW_TOP_PAD


def draw_product_page(
    c,
    product,
    y,
    image_cache=None,
    company_name="BTM",
    category_title="KLS Martin Products",
):
    """
    Draws a single detailed product table row at the given y position.
    Automatically handles page breaking and header redrawing if space is insufficient.
    Returns the updated y-coordinate for the next row.
    """
    no = str(product.get("No.", ""))
    article = str(product.get("model", ""))
    qty = str(product.get("quantity", ""))
    description = product.get("description", "") or ""

    desc_lines = wrap_text_by_width(description, FONT_BODY, FONT_BODY_SIZE, DESC_COL_W)
    text_h = len(desc_lines) * LINE_HEIGHT

    img_url = product.get("image")
    if image_cache is not None and img_url in image_cache:
        img_reader = image_cache[img_url]
    else:
        img_reader = load_image(img_url)

    img_w = img_h = 0
    if img_reader is not None:
        iw, ih = img_reader.getSize()
        scale = min(IMG_MAX_W / iw, IMG_MAX_H / ih, 1.0)
        img_w, img_h = iw * scale, ih * scale

    content_h = max(text_h, img_h, MIN_ROW_CONTENT_H)
    row_block_h = content_h + ROW_BOTTOM_PAD

    if y - row_block_h < MARGIN_BOTTOM:
        c.showPage()
        y = draw_header(c, company_name, category_title)

    c.setFillColor(black)
    c.setFont(FONT_BODY, FONT_BODY_SIZE)
    first_baseline = y - FONT_BODY_SIZE
    c.drawString(COL_NO_X, first_baseline, no)
    c.drawString(COL_CODE_X, first_baseline, article)
    c.drawString(COL_QTY_X, first_baseline, qty)

    ty = first_baseline
    for line in desc_lines:
        c.drawString(COL_DESC_X, ty, line)
        ty -= LINE_HEIGHT

    if img_reader is not None:
        img_x = COL_IMG_X + (COL_IMG_W - img_w) / 2
        img_y = y - img_h
        c.drawImage(
            img_reader,
            img_x,
            img_y,
            width=img_w,
            height=img_h,
            preserveAspectRatio=True,
            mask="auto",
        )

    row_bottom = y - content_h
    line_y = row_bottom - ROW_BOTTOM_PAD
    c.setLineWidth(0.6)
    c.setStrokeColor(HexColor("#333333"))
    c.line(MARGIN_L, line_y, PAGE_W - MARGIN_R, line_y)

    return line_y - ROW_TOP_PAD


# --- Summary Table Component Helpers ---


def _draw_row_borders(c, col_bounds, top_y, bottom_y):
    """Draws grid lines around and between columns for table rows."""
    c.setLineWidth(TABLE_BORDER_WIDTH)
    c.setStrokeColor(TABLE_BORDER_COLOR)
    c.line(col_bounds[0], top_y, col_bounds[-1], top_y)
    c.line(col_bounds[0], bottom_y, col_bounds[-1], bottom_y)
    for x in col_bounds:
        c.line(x, top_y, x, bottom_y)


def draw_summary_table_header(
    c, y, col_bounds=SUMMARY_COL_BOUNDS, labels=SUMMARY_COL_LABELS
):
    """Draws the header row for the summary table."""
    top_y, bottom_y = y, y - TABLE_HEADER_H
    c.setFont("Helvetica-Bold", FONT_BODY_SIZE)
    c.setFillColor(black)
    baseline = bottom_y + (TABLE_HEADER_H - FONT_BODY_SIZE) / 2 + 2
    for i, label in enumerate(labels):
        c.drawString(col_bounds[i] + CELL_PAD_X, baseline, label)
    _draw_row_borders(c, col_bounds, top_y, bottom_y)
    return bottom_y


def draw_summary_row(
    c,
    product,
    y,
    company_name="BTM",
    category_title="Summary",
    col_bounds=SUMMARY_COL_BOUNDS,
):
    """
    Draws a single condensed row in the summary table view with automatic page break handling.
    """
    no = str(product.get("No.", ""))
    article = str(product.get("model", ""))
    qty = str(product.get("quantity", ""))
    description = product.get("description", "") or ""

    desc_lines = wrap_text_by_width(
        description, FONT_BODY, FONT_BODY_SIZE, SUMMARY_DESC_COL_W
    )
    content_h = max(len(desc_lines), 1) * LINE_HEIGHT
    row_h = content_h + 2 * CELL_PAD_Y

    if y - row_h < MARGIN_BOTTOM:
        c.showPage()
        y = draw_header(c, company_name, category_title)
        y = draw_summary_table_header(c, y, col_bounds)

    top_y, bottom_y = y, y - row_h
    c.setFillColor(black)
    c.setFont(FONT_BODY, FONT_BODY_SIZE)
    first_baseline = top_y - CELL_PAD_Y - FONT_BODY_SIZE

    c.drawString(col_bounds[0] + CELL_PAD_X, first_baseline, no)
    c.drawString(col_bounds[1] + CELL_PAD_X, first_baseline, article)
    c.drawString(col_bounds[3] + CELL_PAD_X, first_baseline, qty)

    ty = first_baseline
    for line in desc_lines:
        c.drawString(col_bounds[2] + CELL_PAD_X, ty, line)
        ty -= LINE_HEIGHT

    _draw_row_borders(c, col_bounds, top_y, bottom_y)
    return bottom_y