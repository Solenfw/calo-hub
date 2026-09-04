// Package dto defines API request and response payload shapes.
package dto

// ProductResponse is the public product payload returned by catalog endpoints.
type ProductResponse struct {
	Code        string `json:"code"`
	Eng         string `json:"eng"`
	Viet        string `json:"viet"`
	Alternative string `json:"alternative"`
	Brand       string `json:"brand"`
}

// ImageResponse is the public image payload for a product.
type ImageResponse struct {
	Code   string   `json:"code"`
	Images []string `json:"images"`
}

// MartinReportListResponse is the public payload for Martin report metadata.
type MartinReportListResponse struct {
	Name string `json:"name"`
}

// MartinReportProductResponse is the public payload for products in a Martin report.
type MartinReportProductResponse struct {
	RowNo    int     `json:"row_no"`
	Code     string  `json:"code"`
	Eng      string  `json:"eng"`
	Image    *string `json:"image"`
	Quantity int     `json:"quantity"`
}
