package models

type Products struct {
	Code        string `db:"code" json:"code"`
	Eng         string `db:"eng" json:"eng"`
	Viet        string `db:"viet" json:"viet"`
	Alternative string `db:"alternative" json:"alternative"`
	Brand       string `db:"brand" json:"brand"`
}

type Images struct {
	Code  string   `db:"code" json:"code"`
	Images []string `db:"image" json:"image"`	 
}

type MartinReportList struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"` 
}

type MartinReportProduct struct {
	ReportID int     `db:"report_id" json:"report_id"` 
	RowNo    int     `db:"row_no" json:"row_no"`       
	Code     string  `db:"code" json:"code"`
	Eng		string  `db:"eng" json:"eng"`
	Image    *string  `db:"image" json:"image"`
	Quantity int     `db:"quantity" json:"quantity"`
}