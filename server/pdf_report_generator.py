import time
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


# ==============================================================================
# DATASET
# ==============================================================================

products = [
    {
        "No.": 1,
        "model": "13-448-00-07",
        "description": "Forceps, acc. to Overholt-Fino, fine, curved, length 21.5 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=1%2F1%2Fd%2Ff%2F11dfe495b07733574e2e74ff110ff5f18748a34e_13_448_00_07_h_50.tif",
    },
    {
        "No.": 2,
        "model": "15-921-32-07",
        "description": "marTract®, retractor, acc. to Kirschner,  78 x 65 mm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=6%2F9%2Ff%2F5%2F69f5333fd8ec87a194f85776b2a37a972f24f4c3_15_921_32_07_n.ai",
    },
    {
        "No.": 3,
        "model": "23-747-25-07",
        "description": "Bone holding forceps, acc. to Ulrich, with thread fixation, curved, length 25 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=b%2F4%2F4%2Fc%2Fb44ce77ac707aa24f4ef279f7d65894999bc5f49_23_747_28_07_h_50.tif",
    },
    {
        "No.": 4,
        "model": "17-500-04-04",
        "description": "Power Plug AU f. MedLED",
        "quantity": 1,
        "image": None,
    },
    {
        "No.": 5,
        "model": "24-978-93-14",
        "description": "Punch, with ejector, detachable, solid black, long stroke, 40° upward cutting, mouth width 4 mm, stroke 16 mm, length 18 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=4%2Ff%2F7%2Fd%2F4f7d3db7a8a90df65b5c9351903b7361ac9a0f1a_24_978_91_14_h_100.tif",
    },
    {
        "No.": 6,
        "model": "11-100-16-07",
        "description": "Operating scissors, blunt/blunt, straight, length 16.5 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=5%2F8%2F4%2F6%2F58463db17edb3dc6016d1975f10b1128adb17aef_11_100_14_07_h_50.tif",
    },
    {
        "No.": 7,
        "model": "24-416-31-07",
        "description": "Aortic clamp, Atrauma, acc. to Kowalski, figure 5, clamp length 65 mm, length 27 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=5%2Fa%2F9%2F6%2F5a964c300be75578cad1f41e9210e225d626d361_24_416_31_07_n_100.ai",
    },
    {
        "No.": 8,
        "model": "15-761-05-07",
        "description": "Spreader, acc. to William, right, 1 x 5 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=6%2F0%2Fd%2F5%2F60d5a6195f5d4e0168b9efecbc07a44cd71567e1_15_761_05_07_h_50.tif",
    },
    {
        "No.": 9,
        "model": "11-107-16-07",
        "description": "Operating scissors, pointed/blunt, curved, length 16.5 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=2%2F4%2F5%2F1%2F2451c0bf1ee9023e6740233dd6b1b4f9b7f4099f_11_107_14_07_h_50.tif",
    },
    {
        "No.": 10,
        "model": "15-920-45-07",
        "description": "marTract®, blade holder, with stop, 40 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=a%2F4%2F8%2F2%2Fa482477d4b6a0dbce4f64106a79505f5237a83c6_15_920_45_07_h_.tif",
    },
    {
        "No.": 11,
        "model": "12-589-18-07",
        "description": "Atraumatic, micro forceps, width 1.2 mm, length 18 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=a%2Fe%2F9%2F6%2Fae965f99c9cd556bb3611e74103f1435d8a853ac_12_589_15_07_n_100.eps",
    },
    {
        "No.": 12,
        "model": "15-922-35-07",
        "description": "marTract®, blade, acc. to Breisky, 90 x 35 mm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=f%2F4%2F2%2F9%2Ff4297db6494574623a33fe3d256e1ac5735ef827_15_922_35_07_n.ai",
    },
    {
        "No.": 13,
        "model": "18-524-62-01",
        "description": "Suction tube, acc. to Fergusson, Ø 2.5 mm, working length 130 mm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=9%2F0%2F5%2Fa%2F905affa5ae4375cdc278cb95bd66a5e3dbcb8f94_18_524_63_01.tif",
    },
    {
        "No.": 14,
        "model": "16-160-14-01",
        "description": "Grooved director, straight, length 14.5 cm",
        "quantity": 1,
        "image": None,
    },
    {
        "No.": 15,
        "model": "24-189-12-04",
        "description": "Atrium retractor, acc. to Cooley, malleable, made of nitinol, blade depth 60 mm x blade width 45 mm, accessories for sternum retractor marGate",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=a%2Ff%2F1%2Fd%2Faf1d5f265f1bde801db0c0707b512edc6a3e7dbf_24_189_11_04_h_100.tif",
    },
    {
        "No.": 16,
        "model": "15-623-30-07",
        "description": "Abdominal spatula, acc. to Kader, width 30 mm, length 28 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=a%2F5%2F4%2F7%2Fa54778463758bea352bfb5a7ddbed8e0e4f7acf8_15_623_30_07_h_50.tif",
    },
    {
        "No.": 17,
        "model": "24-131-26-07",
        "description": "Bone shears, acc. to Sauerbruch, length 26 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=2%2F2%2Fa%2Fd%2F22ad28955e28f3553f77e61762414fe9efe5314d_24_131_26_07_h_50.tif",
    },
    {
        "No.": 18,
        "model": "24-189-20-07",
        "description": "Valve retractor, malleable, made of stainless steel, blade depth 18 mm x blade width 65 mm, accessories for sternum retractor marGate",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=e%2F8%2Fb%2F1%2Fe8b1176ca543dd2ca080d0a9a633ed709d55a9b7_24_189_23_07_h_100.tif",
    },
    {
        "No.": 19,
        "model": "20-023-21-09",
        "description": "Microneedle holder, titanium, curved, without ratchet, length 21 cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=2%2F6%2Fa%2Fd%2F26ade6183dc02b030a538a6d917bdd330c2729b9_20_011_15_07_n_100.eps",
    },
    {
        "No.": 20,
        "model": "23-998-42-07",
        "description": "Spoon, lumbar, angled backwards / curved, figure 2, length 24cm",
        "quantity": 1,
        "image": "https://www.klsmartin.com/shop/?eID=akeneo-asset&mediaCode=2%2Fa%2F8%2F1%2F2a816d552dd83c0793ab75a874dbde4a316bde67_23_998_40_07_n_100.ai",
    },
]

# ==============================================================================
# MAIN EXECUTION
# ==============================================================================

if __name__ == "__main__":
    print("Starting PDF generation report for a list of instruments...")
    start = time.time()

    print("Prefetching images...")
    c = canvas.Canvas("products.pdf", pagesize=A4)
    image_cache = prefetch_images(products)

    print("Drawing product pages...")
    y = draw_header(c, "BTM", "KLS Martin Products")
    for p in products:
        y = draw_product_page(c, p, y, image_cache=image_cache)

    c.showPage()

    print("Drawing summary page...")
    y = draw_header(c, "BTM", "Summary")
    y = draw_summary_table_header(c, y)
    for p in products:
        y = draw_summary_row(c, p, y)

    c.save()

    elapsed = time.time() - start
    print(f"Runtime: {elapsed:.3f} sec ({elapsed*1000:.1f} ms)")