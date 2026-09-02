package dto

// request/response structs, replaces Pydantic schemas
type ProductResponse struct {
	Code        string `json:"code"`
	Eng         string `json:"eng"`
	Viet        string `json:"viet"`
	Alternative string `json:"alternative"`
	Brand       string `json:"brand"`
}

type ImageResponse struct {
	Code  	string `json:"code"`
	Images []string `json:"images"`
}

type MartinReportListResponse struct {
	Name string `json:"name"`
}

type MartinReportProductResponse struct {
	RowNo       int    `json:"row_no"`
	Code        string `json:"code"`
	Eng 			string `json:"eng"`
	Image      	*string `json:"image"`
	Quantity    int    `json:"quantity"`
}

