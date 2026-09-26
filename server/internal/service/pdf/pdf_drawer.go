// Package pdf renders Martin report data into PDF bytes. This package only
// builds PDF output from already-fetched data — no repository or HTTP calls
// beyond fetching the product images themselves, which are already-known
// URLs stored on each product.
//
// Ported from a reportlab (Python) script. ReportLab's canvas has y=0 at the
// page BOTTOM, increasing upward. fpdf has y=0 at the page TOP, increasing
// downward — every y-movement below was re-derived from actual up/down
// intent, not a mechanical sign flip. Two spots where a naive flip would be
// wrong: image anchor placement (ReportLab anchors bottom-left and grows up;
// fpdf anchors top-left and grows down — see drawProductRow), and the
// summary header's baseline centering (stays a subtraction in both).
package pdf

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "embed"
	_ "image/jpeg" // registers jpeg decoding for image.DecodeConfig
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/solenfw/calo-hub/internal/models"
)

// ==============================================================================
// CONFIGURATION & LAYOUT CONSTANTS
// ==============================================================================
//
// The document unit is points ("pt"), matching what ReportLab's canvas
// actually uses internally — `* cm` in the Python script was always just a
// convenience multiplier resolving to a point value at import time, same as
// here via the cm() helper. Bare numeric constants in the Python script
// (LINE_HEIGHT, TABLE_BORDER_WIDTH, font sizes) were already in points and
// are carried over unchanged.

const ptPerCm = 28.346456693

func cm(v float64) float64 { return v * ptPerCm }

var (
	pageW = cm(21.0) // A4
	pageH = cm(29.7)

	marginL      = cm(1.8)
	marginR      = cm(1.8)
	marginTop    = cm(1.6)
	marginBottom = cm(1.8)
	contentW     = pageW - marginL - marginR

	colNoX   = marginL
	colCodeX = marginL + cm(1.0)
	colDescX = marginL + cm(3.6)
	colQtyX  = marginL + cm(10.8)
	colImgX  = marginL + cm(12.6)
	colImgW  = contentW - cm(12.6)

	descColW = colQtyX - colDescX - cm(0.3)

	fontBody     = "Roboto"
	fontBodySize = 10.5
	lineHeight   = 14.0 // already points in the source script — unchanged

	imgMaxW = colImgW - cm(0.2)
	imgMaxH = cm(3.6)

	rowTopPad      = cm(0.55)
	rowBottomPad   = cm(0.55)
	minRowContentH = cm(2.3)

	tableHeaderH = cm(0.9)
	cellPadX     = cm(0.15)
	cellPadY     = cm(0.25)
	tableBorderW = 0.6 // bare in the source too — already points

	summaryColBounds = []float64{
		marginL,
		colCodeX,
		colDescX,
		colQtyX + cm(5),
		pageW - marginR,
	}
	summaryColLabels = []string{"No", "code", "Description", "Qty"}
	summaryDescColW  = (colQtyX + cm(5)) - colDescX - 2*cellPadX
)

const (
	maxImageWorkers   = 10
	imageFetchTimeout = 15 * time.Second
)

var httpClient = &http.Client{Timeout: imageFetchTimeout}

//go:embed assets/Roboto-Regular.ttf
var robotoRegularFontBytes []byte


// ==============================================================================
// ENTRY POINT
// ==============================================================================

// GeneratePDF renders a Martin report as a PDF: a detailed per-product page
// followed by a condensed summary table. Takes already-fetched product rows
// (via repository.GetMartinReportProductsByReportID) — no DB access happens
// in this package.
func GeneratePDF(ctx context.Context, companyName, reportName string, products []models.MartinReportProduct) ([]byte, error) {
	imageCache := prefetchImages(ctx, products)
	if len(robotoRegularFontBytes) == 0 {
		return nil, fmt.Errorf("pdf font bytes missing: embedded Roboto font not loaded")
	}

	doc := fpdf.New("P", "pt", "A4", "")
	doc.SetAutoPageBreak(false, 0) // pagination is handled manually below, same as the Python script's showPage() calls
	doc.AddUTF8FontFromBytes("Roboto", "", robotoRegularFontBytes)
	doc.AddUTF8FontFromBytes("Roboto", "B", robotoRegularFontBytes)
	doc.AddPage()
	y := drawHeader(doc, companyName, reportName)
	for _, p := range products {
		y = drawProductRow(doc, p, y, imageCache)
	}

	// add summary table.
	// doc.AddPage()
	// y = drawHeader(doc, companyName, "Summary")
	// y = drawSummaryTableHeader(doc, y)
	// for _, p := range products {
	// 	y = drawSummaryRow(doc, p, y, companyName)
	// }

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, fmt.Errorf("rendering pdf: %w", err)
	}
	return buf.Bytes(), nil
}

// ==============================================================================
// IMAGE FETCHING
// ==============================================================================

// prefetchImages downloads every distinct image URL concurrently (bounded
// worker pool, same pattern as cmd/migrate-images-local's download loop) and
// returns a url -> bytes cache. A failed fetch is cached as an explicit nil
// entry so drawProductRow treats it as "no image" instead of retrying.
func prefetchImages(ctx context.Context, products []models.MartinReportProduct) map[string][]byte {
	urlSet := make(map[string]struct{})
	for _, p := range products {
		if p.Image != nil && *p.Image != "" {
			urlSet[*p.Image] = struct{}{}
		}
	}
	if len(urlSet) == 0 {
		return nil
	}

	cache := make(map[string][]byte, len(urlSet))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxImageWorkers)

	for u := range urlSet {
		wg.Add(1)
		sem <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()

			data, err := fetchImageBytes(ctx, url)
			mu.Lock()
			if err == nil {
				cache[url] = data
			} else {
				cache[url] = nil
			}
			mu.Unlock()
		}(u)
	}
	wg.Wait()
	return cache
}

func fetchImageBytes(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func imageDimensions(data []byte) (w, h float64, err error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, err
	}
	return float64(cfg.Width), float64(cfg.Height), nil
}

// ==============================================================================
// TEXT UTILITIES
// ==============================================================================

func wrapTextByWidth(doc *fpdf.Fpdf, text string, maxWidth float64) []string {
	if text == "" {
		return []string{""}
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, word := range words {
		candidate := strings.TrimSpace(current + " " + word)
		if doc.GetStringWidth(candidate) <= maxWidth {
			current = candidate
		} else {
			if current != "" {
				lines = append(lines, current)
			}
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

// ==============================================================================
// PDF DRAWING FUNCTIONS
// ==============================================================================

// drawHeader draws the page masthead and column titles, returning the y
// where the first row should start.
func drawHeader(doc *fpdf.Fpdf, companyName, categoryTitle string) float64 {
	y := marginTop

	doc.SetTextColor(0, 0, 0)
	doc.SetFont(fontBody, "B", 15)
	doc.Text(marginL, y, companyName)
	y += cm(0.35)
	doc.SetLineWidth(1)
	doc.Line(marginL, y, pageW-marginR, y)

	y += cm(0.65)
	doc.SetFont(fontBody, "B", 12.5)
	doc.Text(marginL, y, categoryTitle)

	y += cm(0.85)
	if categoryTitle != "Summary" {
		doc.SetFont(fontBody, "", 11)
		doc.Text(colNoX, y, "No.")
		doc.Text(colCodeX, y, "Code")
		doc.Text(colDescX, y, "Description")
		doc.Text(colQtyX, y, "Quantity")
		y += cm(0.25)
		doc.Line(marginL, y, pageW-marginR, y)
	}

	return y + rowTopPad
}

// drawProductRow draws one detailed product row at y, handling page breaks
// automatically. Returns the y for the next row.
func drawProductRow(doc *fpdf.Fpdf, product models.MartinReportProduct, y float64, imageCache map[string][]byte) float64 {
	doc.SetFont(fontBody, "", fontBodySize)
	descLines := wrapTextByWidth(doc, product.Description, descColW)
	textH := float64(len(descLines)) * lineHeight

	var imgData []byte
	if product.Image != nil {
		imgData = imageCache[*product.Image]
	}

	haveImage := imgData != nil
	var imgW, imgH float64
	if haveImage {
		iw, ih, err := imageDimensions(imgData)
		if err != nil || iw <= 0 || ih <= 0 {
			haveImage = false
		} else {
			// Same scaling approach as the Python script: pixel dimensions
			// are treated as point-equivalent for the fit calculation,
			// preserved as-is rather than made DPI-aware.
			scale := math.Min(imgMaxW/iw, imgMaxH/ih)
			if scale > 1.0 {
				scale = 1.0
			}
			imgW, imgH = iw*scale, ih*scale
		}
	}

	contentH := math.Max(textH, math.Max(imgH, minRowContentH))
	rowBlockH := contentH + rowBottomPad

	if y+rowBlockH > pageH-marginBottom {
		doc.AddPage()
		y = drawHeader(doc, "", "")
	}

	doc.SetTextColor(0, 0, 0)
	doc.SetFont(fontBody, "", fontBodySize)
	firstBaseline := y + fontBodySize
	doc.Text(colNoX, firstBaseline, fmt.Sprintf("%d", product.RowNo))
	doc.Text(colCodeX, firstBaseline, product.Code)
	doc.Text(colQtyX, firstBaseline, fmt.Sprintf("%d", product.Quantity))

	ty := firstBaseline
	for _, line := range descLines {
		doc.Text(colDescX, ty, line)
		ty += lineHeight
	}

	if haveImage {
		imgX := colImgX + (colImgW-imgW)/2
		// fpdf anchors images top-left and grows downward, so the row's
		// top (y) IS the image's top edge directly — no subtraction, unlike
		// the Python version's `y - img_h` (which computed a bottom-left
		// anchor for ReportLab's upward-growing images).
		imgY := y
		placeImage(doc, fmt.Sprintf("row-%d", product.RowNo), imgData, imgX, imgY, imgW, imgH)
	}

	rowBottom := y + contentH
	lineY := rowBottom + rowBottomPad
	doc.SetLineWidth(tableBorderW)
	doc.SetDrawColor(0x33, 0x33, 0x33)
	doc.Line(marginL, lineY, pageW-marginR, lineY)

	return lineY + rowTopPad
}

// placeImage registers image bytes under a name unique to this document
// build (RowNo is unique within one report's product slice, and a fresh
// *fpdf.Fpdf is created per GeneratePDF call, so no collision risk).
func placeImage(doc *fpdf.Fpdf, name string, data []byte, x, y, w, h float64) {
	opt := fpdf.ImageOptions{ImageType: "JPG"}
	doc.RegisterImageOptionsReader(name, opt, bytes.NewReader(data))
	doc.ImageOptions(name, x, y, w, h, false, opt, 0, "")
}

// --- Summary table ---

func drawRowBorders(doc *fpdf.Fpdf, colBounds []float64, topY, bottomY float64) {
	doc.SetLineWidth(tableBorderW)
	doc.SetDrawColor(0x33, 0x33, 0x33)
	doc.Line(colBounds[0], topY, colBounds[len(colBounds)-1], topY)
	doc.Line(colBounds[0], bottomY, colBounds[len(colBounds)-1], bottomY)
	for _, x := range colBounds {
		doc.Line(x, topY, x, bottomY)
	}
}

func drawSummaryTableHeader(doc *fpdf.Fpdf, y float64) float64 {
	topY, bottomY := y, y+tableHeaderH
	doc.SetFont(fontBody, "B", fontBodySize)
	doc.SetTextColor(0, 0, 0)
	// Centering offset stays a SUBTRACTION here — this is the other spot
	// where a mechanical sign-flip would be wrong. See the package comment.
	baseline := bottomY - (tableHeaderH-fontBodySize)/2 - 2
	for i, label := range summaryColLabels {
		doc.Text(summaryColBounds[i]+cellPadX, baseline, label)
	}
	drawRowBorders(doc, summaryColBounds, topY, bottomY)
	return bottomY
}

func drawSummaryRow(doc *fpdf.Fpdf, product models.MartinReportProduct, y float64, companyName string) float64 {
	doc.SetFont(fontBody, "", fontBodySize)
	descLines := wrapTextByWidth(doc, product.Description, summaryDescColW)
	contentH := math.Max(float64(len(descLines)), 1) * lineHeight
	rowH := contentH + 2*cellPadY

	if y+rowH > pageH-marginBottom {
		doc.AddPage()
		y = drawHeader(doc, companyName, "")
		y = drawSummaryTableHeader(doc, y)
	}

	topY, bottomY := y, y+rowH
	doc.SetTextColor(0, 0, 0)
	doc.SetFont(fontBody, "", fontBodySize)
	firstBaseline := topY + cellPadY + fontBodySize

	doc.Text(summaryColBounds[0]+cellPadX, firstBaseline, fmt.Sprintf("%d", product.RowNo))
	doc.Text(summaryColBounds[1]+cellPadX, firstBaseline, product.Code)
	doc.Text(summaryColBounds[3]+cellPadX, firstBaseline, fmt.Sprintf("%d", product.Quantity))

	ty := firstBaseline
	for _, line := range descLines {
		doc.Text(summaryColBounds[2]+cellPadX, ty, line)
		ty += lineHeight
	}

	drawRowBorders(doc, summaryColBounds, topY, bottomY)
	return bottomY
}