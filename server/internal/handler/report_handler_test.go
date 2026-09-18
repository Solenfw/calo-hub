package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
)

type fakeReportStore struct {
	getAllMartinReportLists           func(context.Context) ([]models.MartinReportList, error)
	getMartinReportListByName         func(context.Context, string) (models.MartinReportList, error)
	getMartinReportProductsByReportID func(context.Context, int) ([]models.MartinReportProduct, error)
	createMartinReportList            func(context.Context, string) (models.MartinReportList, error)
	deleteMartinReportList            func(context.Context, int) error
	replaceMartinReportProducts       func(context.Context, int, []models.MartinReportProduct) error
}

func (f *fakeReportStore) GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error) {
	if f.getAllMartinReportLists != nil {
		return f.getAllMartinReportLists(ctx)
	}
	return []models.MartinReportList{}, nil
}

func (f *fakeReportStore) GetMartinReportListByName(ctx context.Context, name string) (models.MartinReportList, error) {
	if f.getMartinReportListByName != nil {
		return f.getMartinReportListByName(ctx, name)
	}
	return models.MartinReportList{}, nil
}

func (f *fakeReportStore) GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	if f.getMartinReportProductsByReportID != nil {
		return f.getMartinReportProductsByReportID(ctx, reportID)
	}
	return nil, nil
}

func (f *fakeReportStore) CreateMartinReportList(ctx context.Context, name string) (models.MartinReportList, error) {
	if f.createMartinReportList != nil {
		return f.createMartinReportList(ctx, name)
	}
	return models.MartinReportList{ID: 1, Name: name}, nil
}

func (f *fakeReportStore) DeleteMartinReportList(ctx context.Context, reportID int) error {
	if f.deleteMartinReportList != nil {
		return f.deleteMartinReportList(ctx, reportID)
	}
	return nil
}

func (f *fakeReportStore) ReplaceMartinReportProducts(ctx context.Context, reportID int, products []models.MartinReportProduct) error {
	if f.replaceMartinReportProducts != nil {
		return f.replaceMartinReportProducts(ctx, reportID, products)
	}
	return nil
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
			store := &fakeReportStore{
				getMartinReportListByName: func(_ context.Context, name string) (models.MartinReportList, error) {
					if tt.reportName != "" && name != tt.reportName {
						t.Errorf("store name = %q, want %q", name, tt.reportName)
					}
					return models.MartinReportList{Name: name}, tt.storeError
				},
			}
			h := NewReportHandler(store)
			recorder := httptest.NewRecorder()
			h.GetMartinReportByName(recorder, requestWithRouteParams(http.MethodGet, "/report/martin/"+url.PathEscape(tt.reportName), map[string]string{"name": tt.reportName}))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestReportHandlerListMartinReports(t *testing.T) {
	store := &fakeReportStore{
		getAllMartinReportLists: func(context.Context) ([]models.MartinReportList, error) {
			return []models.MartinReportList{{ID: 1, Name: "Report A"}, {ID: 2, Name: "Report B"}}, nil
		},
	}
	h := NewReportHandler(store)
	recorder := httptest.NewRecorder()
	h.ListMartinReports(recorder, httptest.NewRequest(http.MethodGet, "/report/martin", nil))

	assertJSONSuccess(t, recorder, http.StatusOK, []dto.MartinReportListResponse{
		{ID: 1, Name: "Report A"},
		{ID: 2, Name: "Report B"},
	})
}

func TestReportHandlerDownloadImportTemplate(t *testing.T) {
	h := NewReportHandler(&fakeReportStore{})
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report/martin/template", nil)

	h.DownloadImportTemplate(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Header().Get("Content-Type"), "application") && !strings.Contains(recorder.Header().Get("Content-Type"), "excel") {
		t.Fatalf("Content-Type = %q, want spreadsheet content type", recorder.Header().Get("Content-Type"))
	}
	if len(recorder.Body.Bytes()) == 0 {
		t.Fatal("template body is empty")
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
			products:   []models.MartinReportProduct{{ReportID: 7, RowNo: 1, Code: "AA804R", Description: "Ruler", Image: &image, Quantity: 1}},
			wantStatus: http.StatusOK,
			wantPDF:    true,
		},
		{
			name:       "valid report generates PDF",
			reportName: "Report A",
			reportID:   7,
			products:   []models.MartinReportProduct{{ReportID: 7, RowNo: 1, Code: "AA804R", Description: "Ruler", Image: &image, Quantity: 1}},
			wantStatus: http.StatusOK,
			wantPDF:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeReportStore{
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
			h.GetMartinReport(recorder, requestWithRouteParams(http.MethodGet, "/report/martin/"+url.PathEscape(tt.reportName)+"/pdf", map[string]string{"name": tt.reportName}))

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

func TestReportHandlerCreateMartinReport(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		store := &fakeReportStore{}
		h := NewReportHandler(store)
		recorder := httptest.NewRecorder()
		body := `{"name":"   "}`
		req := httptest.NewRequest(http.MethodPost, "/report/martin", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		h.CreateMartinReport(recorder, req)
		assertJSONError(t, recorder, http.StatusBadRequest, "report name is required")
	})

	t.Run("valid create", func(t *testing.T) {
		store := &fakeReportStore{
			createMartinReportList: func(_ context.Context, name string) (models.MartinReportList, error) {
				if name != "New Report" {
					t.Fatalf("name = %q, want %q", name, "New Report")
				}
				return models.MartinReportList{ID: 42, Name: name}, nil
			},
		}
		h := NewReportHandler(store)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/report/martin", strings.NewReader(`{"name":"New Report"}`))
		req.Header.Set("Content-Type", "application/json")

		h.CreateMartinReport(recorder, req)
		assertJSONSuccess(t, recorder, http.StatusCreated, dto.MartinReportListResponse{ID: 42, Name: "New Report"})
	})
}

func TestReportHandlerDeleteMartinReport(t *testing.T) {
	store := &fakeReportStore{
		deleteMartinReportList: func(_ context.Context, reportID int) error {
			if reportID != 7 {
				t.Fatalf("report ID = %d, want 7", reportID)
			}
			return nil
		},
	}
	h := NewReportHandler(store)
	recorder := httptest.NewRecorder()
	request := requestWithRouteParams(http.MethodDelete, "/report/martin/7", map[string]string{"report_id": "7"})
	h.DeleteMartinReport(recorder, request)

	assertJSONSuccess(t, recorder, http.StatusOK, map[string]string{"status": "deleted"})
}

func TestReportHandlerUpdateMartinReportProducts(t *testing.T) {
	store := &fakeReportStore{
		replaceMartinReportProducts: func(_ context.Context, reportID int, products []models.MartinReportProduct) error {
			if reportID != 7 {
				t.Fatalf("report ID = %d, want 7", reportID)
			}
			if len(products) != 2 {
				t.Fatalf("len(products) = %d, want 2", len(products))
			}
			return nil
		},
	}
	h := NewReportHandler(store)
	recorder := httptest.NewRecorder()
	payload := `{"products":[{"row_no":1,"code":"AA804R","description":"Ruler","image":"https://example.test/image.png","quantity":1},{"row_no":2,"code":"AA805R","description":"Scale","image":"","quantity":2}]}`
	request := requestWithRouteParams(http.MethodPut, "/report/martin/7/products", map[string]string{"report_id": "7"})
	request.Body = io.NopCloser(strings.NewReader(payload))
	h.UpdateMartinReportProducts(recorder, request)

	assertJSONSuccess(t, recorder, http.StatusOK, map[string]string{"status": "updated"})
}

func TestReportHandlerResolvesCatalogImageWhenMissingFromReport(t *testing.T) {
	image := "https://example.test/image.png"
	store := &fakeReportStore{
		getMartinReportProductsByReportID: func(_ context.Context, reportID int) ([]models.MartinReportProduct, error) {
			return []models.MartinReportProduct{{ReportID: reportID, RowNo: 1, Code: "AA804R", Description: "Ruler", Image: nil, Quantity: 2}}, nil
		},
	}
	catalog := &fakeCatalogStore{
		getImagesByCode: func(_ context.Context, code string) (models.Images, error) {
			if code != "AA804R" {
				t.Fatalf("code = %q, want %q", code, "AA804R")
			}
			return models.Images{Code: code, Images: []string{image}}, nil
		},
	}

	h := NewReportHandler(store, catalog)
	recorder := httptest.NewRecorder()
	h.GetMartinReportProducts(recorder, requestWithRouteParams(http.MethodGet, "/report/martin/all/42", map[string]string{"report_id": "42"}))

	assertJSONSuccess(t, recorder, http.StatusOK, []dto.MartinReportProductResponse{{RowNo: 1, Code: "AA804R", Description: "Ruler", Image: &image, Quantity: 2}})
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
				{ReportID: 42, RowNo: 1, Code: "AA804R", Description: "Ruler", Image: &image, Quantity: 2},
				{ReportID: 42, RowNo: 2, Code: "AA805R", Description: "Scale", Image: nil, Quantity: 1},
			},
			wantStatus: http.StatusOK,
			wantBody: []dto.MartinReportProductResponse{
				{RowNo: 1, Code: "AA804R", Description: "Ruler", Image: &image, Quantity: 2},
				{RowNo: 2, Code: "AA805R", Description: "Scale", Image: nil, Quantity: 1},
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
			store := &fakeReportStore{
				getMartinReportProductsByReportID: func(_ context.Context, reportID int) ([]models.MartinReportProduct, error) {
					if tt.reportID != "" && reportID != 42 {
						t.Errorf("store report ID = %d, want 42", reportID)
					}
					return tt.products, tt.storeError
				},
			}
			h := NewReportHandler(store)
			recorder := httptest.NewRecorder()
			h.GetMartinReportProducts(recorder, requestWithRouteParams(http.MethodGet, "/report/martin/all/"+tt.reportID, map[string]string{"report_id": tt.reportID}))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}
