package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
)

type fakeCatalogReportStore struct {
	getMartinReportListByName         func(context.Context, string) (models.MartinReportList, error)
	getMartinReportProductsByReportID func(context.Context, int) ([]models.MartinReportProduct, error)
}

func (f *fakeCatalogReportStore) GetMartinReportListByName(ctx context.Context, name string) (models.MartinReportList, error) {
	if f.getMartinReportListByName != nil {
		return f.getMartinReportListByName(ctx, name)
	}
	return models.MartinReportList{}, nil
}

func (f *fakeCatalogReportStore) GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	if f.getMartinReportProductsByReportID != nil {
		return f.getMartinReportProductsByReportID(ctx, reportID)
	}
	return nil, nil
}

func TestCatalogHandlerGetMartinReportByName(t *testing.T) {
	tests := []struct {
		name       string
		reportName string
		storeError error
		wantStatus int
		wantError  string
		wantBody   dto.MartinReportListResponse
	}{
		{
			name:       "missing name",
			wantStatus: http.StatusBadRequest,
			wantError:  "report name is required",
		},
		{
			name:       "found report",
			reportName: "Report A",
			wantStatus: http.StatusOK,
			wantBody:   dto.MartinReportListResponse{Name: "Report A"},
		},
		{
			name:       "not found",
			reportName: "Missing",
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "unexpected store error",
			reportName: "Report A",
			storeError: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogReportStore{
				getMartinReportListByName: func(_ context.Context, name string) (models.MartinReportList, error) {
					if tt.reportName != "" && name != tt.reportName {
						t.Errorf("store name = %q, want %q", name, tt.reportName)
					}
					return models.MartinReportList{Name: name}, tt.storeError
				},
			}
			h := NewReportHandler(store)
			recorder := httptest.NewRecorder()
			h.GetMartinReportByName(recorder, requestWithRouteParams(http.MethodGet, "/catalog/report/martin/"+url.PathEscape(tt.reportName), map[string]string{"name": tt.reportName}))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestReportHandlerGetMartinReport(t *testing.T) {
	image := "https://example.test/image.png"
	tests := []struct {
		name       string
		reportName string
		reportID   int
		products   []models.MartinReportProduct
		storeError error
		wantStatus int
		wantError  string
		wantPDF    bool
	}{
		{
			name:       "missing report name",
			wantStatus: http.StatusBadRequest,
			wantError:  "report name is required",
		},
		{
			name:       "not found report metadata",
			reportName: "Missing",
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "not found report products",
			reportName: "Report A",
			reportID:   7,
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "failed to generate PDF",
			reportName: "Report A",
			reportID:   7,
			products:   []models.MartinReportProduct{{ReportID: 7, RowNo: 1, Code: "AA804R", Eng: "Ruler", Image: &image, Quantity: 1}},
			wantStatus: http.StatusOK,
			wantPDF:    true,
		},
		{
			name:       "valid report generates PDF",
			reportName: "Report A",
			reportID:   7,
			products:   []models.MartinReportProduct{{ReportID: 7, RowNo: 1, Code: "AA804R", Eng: "Ruler", Image: &image, Quantity: 1}},
			wantStatus: http.StatusOK,
			wantPDF:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogReportStore{
				getMartinReportListByName: func(_ context.Context, name string) (models.MartinReportList, error) {
					if tt.reportName != "" && name != tt.reportName {
						t.Errorf("store name = %q, want %q", name, tt.reportName)
					}
					if tt.storeError == repository.ErrNotFound && tt.reportName == "Missing" {
						return models.MartinReportList{}, tt.storeError
					}
					return models.MartinReportList{ID: tt.reportID, Name: name}, tt.storeError
				},
				getMartinReportProductsByReportID: func(_ context.Context, reportID int) ([]models.MartinReportProduct, error) {
					if tt.reportID != 0 && reportID != tt.reportID {
						t.Errorf("store report ID = %d, want %d", reportID, tt.reportID)
					}
					if tt.storeError == repository.ErrNotFound && tt.reportName == "Missing" {
						return nil, tt.storeError
					}
					return tt.products, tt.storeError
				},
			}
			h := NewReportHandler(store)
			recorder := httptest.NewRecorder()
			h.GetMartinReport(recorder, requestWithRouteParams(http.MethodGet, "/catalog/report/martin/"+url.PathEscape(tt.reportName)+"/pdf", map[string]string{"name": tt.reportName}))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			if tt.wantPDF {
				if recorder.Code != tt.wantStatus {
					t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
				}
				if recorder.Header().Get("Content-Type") != "application/pdf" {
					t.Fatalf("Content-Type = %q, want application/pdf", recorder.Header().Get("Content-Type"))
				}
				if len(recorder.Body.Bytes()) == 0 {
					t.Fatal("PDF body is empty")
				}
				return
			}
		})
	}
}

func TestCatalogHandlerGetMartinReportProducts(t *testing.T) {
	image := "https://example.test/image.png"
	tests := []struct {
		name       string
		reportID   string
		products   []models.MartinReportProduct
		storeError error
		wantStatus int
		wantError  string
		wantBody   []dto.MartinReportProductResponse
	}{
		{
			name:       "missing report ID",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid report ID",
		},
		{
			name:       "non numeric report ID",
			reportID:   "abc",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid report ID",
		},
		{
			name:       "zero report ID",
			reportID:   "0",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid report ID",
		},
		{
			name:       "negative report ID",
			reportID:   "-1",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid report ID",
		},
		{
			name:     "valid report ID returns DTOs",
			reportID: "42",
			products: []models.MartinReportProduct{
				{ReportID: 42, RowNo: 1, Code: "AA804R", Eng: "Ruler", Image: &image, Quantity: 2},
				{ReportID: 42, RowNo: 2, Code: "AA805R", Eng: "Scale", Image: nil, Quantity: 1},
			},
			wantStatus: http.StatusOK,
			wantBody: []dto.MartinReportProductResponse{
				{RowNo: 1, Code: "AA804R", Eng: "Ruler", Image: &image, Quantity: 2},
				{RowNo: 2, Code: "AA805R", Eng: "Scale", Image: nil, Quantity: 1},
			},
		},
		{
			name:       "zero products become empty array",
			reportID:   "42",
			products:   []models.MartinReportProduct{},
			wantStatus: http.StatusOK,
			wantBody:   []dto.MartinReportProductResponse{},
		},
		{
			name:       "not found",
			reportID:   "42",
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "unexpected store error",
			reportID:   "42",
			storeError: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogReportStore{
				getMartinReportProductsByReportID: func(_ context.Context, reportID int) ([]models.MartinReportProduct, error) {
					if tt.reportID != "" && reportID != 42 {
						t.Errorf("store report ID = %d, want 42", reportID)
					}
					return tt.products, tt.storeError
				},
			}
			h := NewReportHandler(store)
			recorder := httptest.NewRecorder()
			h.GetMartinReportProducts(recorder, requestWithRouteParams(http.MethodGet, "/catalog/report/martin/all/"+tt.reportID, map[string]string{"report_id": tt.reportID}))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}
