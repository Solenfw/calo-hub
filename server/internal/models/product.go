// Package models defines the domain records used by the repository and handler
// layers.
//
// These structs correspond closely to database rows and are converted into DTO
// payloads at the HTTP boundary so the application keeps storage details and API
// contracts intentionally separate.
package models

// Products represents a product row from the catalog.
type Products struct {
	Code        string `db:"code" json:"code"`
	Eng         string `db:"eng" json:"eng"`
	Viet        string `db:"viet" json:"viet"`
	Alternative string `db:"alternative" json:"alternative"`
	Brand       string `db:"brand" json:"brand"`
}

// Images represents image URLs associated with a product code.
type Images struct {
	Code   string   `db:"code" json:"code"`
	Images []string `db:"images" json:"images"`
}

// MartinReportList represents a saved Martin report list.
type MartinReportList struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// MartinReportProduct represents a product row within a Martin report.
type MartinReportProduct struct {
	ReportID int     `db:"report_id" json:"report_id"`
	RowNo    int     `db:"row_no" json:"row_no"`
	Code     string  `db:"code" json:"code"`
	Eng      string  `db:"eng" json:"eng"`
	Image    *string `db:"image" json:"image"`
	Quantity int     `db:"quantity" json:"quantity"`
}
