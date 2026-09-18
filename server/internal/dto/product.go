// Package dto defines the public API payload shapes sent between the HTTP layer
// and clients.
//
// These models intentionally keep database internals out of the response schema
// and present a stable contract that the frontend can depend on regardless of
// how the repository stores the data internally.
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MartinReportProductResponse is the public payload for products in a Martin report.
type MartinReportProductResponse struct {
	RowNo    		int     `json:"row_no"`
	Code     		string  `json:"code"`
	Description    string  `json:"description"`
	Image    		*string `json:"image"`
	Quantity 		int     `json:"quantity"`
}

// MartinReportProductRequest is the payload used when replacing the instruments
// attached to a Martin report.
type MartinReportProductRequest struct {
	RowNo    		int     `json:"row_no"`
	Code     		string  `json:"code"`
	Description    string  `json:"description"`
	Image    		string  `json:"image"`
	Quantity 		int     `json:"quantity"`
}

// MartinReportProductsRequest contains the full set of rows to save for a report.
type MartinReportProductsRequest struct {
	Products []MartinReportProductRequest `json:"products"`
}

