package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/solenfw/calo-hub/internal/handler"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/router"
)

type recordingCatalogStore struct {
	called   string
	code     string
	name     string
	reportID int
	terms    []string
}

func (s *recordingCatalogStore) reset() {
	s.called = ""
	s.code = ""
	s.name = ""
	s.reportID = 0
	s.terms = nil
}

func (s *recordingCatalogStore) GetProductByCode(_ context.Context, code string) (models.Products, error) {
	s.called = "GetProductByCode"
	s.code = code
	return models.Products{Code: code}, nil
}

func (s *recordingCatalogStore) GetProductsByTerms(_ context.Context, terms []string) ([]models.Products, error) {
	s.called = "GetProductsByTerms"
	s.terms = append([]string(nil), terms...)
	return []models.Products{}, nil
}

func (s *recordingCatalogStore) GetImagesByCode(_ context.Context, code string) (models.Images, error) {
	s.called = "GetImagesByCode"
	s.code = code
	return models.Images{Code: code, Images: []string{}}, nil
}

func (s *recordingCatalogStore) GetAllMartinReportLists(context.Context) ([]models.MartinReportList, error) {
	s.called = "GetAllMartinReportLists"
	return []models.MartinReportList{}, nil
}

func (s *recordingCatalogStore) GetMartinReportListByName(_ context.Context, name string) (models.MartinReportList, error) {
	s.called = "GetMartinReportListByName"
	s.name = name
	return models.MartinReportList{Name: name}, nil
}

func (s *recordingCatalogStore) GetMartinReportProductsByReportID(_ context.Context, reportID int) ([]models.MartinReportProduct, error) {
	s.called = "GetMartinReportProductsByReportID"
	s.reportID = reportID
	return []models.MartinReportProduct{}, nil
}

func TestNewRoutesDispatchToExpectedHandler(t *testing.T) {
	store := &recordingCatalogStore{}
	catalogHandler := handler.NewCatalogHandler(store)
	catalogReportHandler := handler.NewReportHandler(store)
	r := router.New(catalogHandler, catalogReportHandler)

	tests := []struct {
		name       string
		path       string
		wantCalled string
		wantCode   string
		wantName   string
		wantID     int
		wantTerms  []string
	}{
		{
			name:       "products search route",
			path:       "/catalog/products?q=knife",
			wantCalled: "GetProductsByTerms",
			wantTerms:  []string{"knife"},
		},
		{
			name:       "product code route",
			path:       "/catalog/products/AA804R",
			wantCalled: "GetProductByCode",
			wantCode:   "AA804R",
		},
		{
			name:       "images route",
			path:       "/catalog/images/AA804R",
			wantCalled: "GetImagesByCode",
			wantCode:   "AA804R",
		},
		{
			name:       "Martin report list route",
			path:       "/catalog/report/martin",
			wantCalled: "GetAllMartinReportLists",
		},
		{
			name:       "Martin report name route",
			path:       "/catalog/report/martin/monthly",
			wantCalled: "GetMartinReportListByName",
			wantName:   "monthly",
		},
		{
			name:       "Martin report products route",
			path:       "/catalog/report/martin/all/42",
			wantCalled: "GetMartinReportProductsByReportID",
			wantID:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.reset()
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
			}
			if store.called != tt.wantCalled {
				t.Fatalf("called = %q, want %q", store.called, tt.wantCalled)
			}
			if tt.wantCode != "" && store.code != tt.wantCode {
				t.Fatalf("code = %q, want %q", store.code, tt.wantCode)
			}
			if tt.wantName != "" && store.name != tt.wantName {
				t.Fatalf("name = %q, want %q", store.name, tt.wantName)
			}
			if tt.wantID != 0 && store.reportID != tt.wantID {
				t.Fatalf("report ID = %d, want %d", store.reportID, tt.wantID)
			}
			if tt.wantTerms != nil && !reflect.DeepEqual(store.terms, tt.wantTerms) {
				t.Fatalf("terms = %#v, want %#v", store.terms, tt.wantTerms)
			}
		})
	}
}

func TestNewRoutesRejectsWrongMethod(t *testing.T) {
	store := &recordingCatalogStore{}
	r := router.New(handler.NewCatalogHandler(store))
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/catalog/products?q=knife", nil))

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if store.called != "" {
		t.Fatalf("handler called = %q, want no handler invocation", store.called)
	}
}

func TestNewRoutesReturnsNotFoundForUnknownPath(t *testing.T) {
	store := &recordingCatalogStore{}
	r := router.New(handler.NewCatalogHandler(store))
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/catalog/not-a-route", nil))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
	if store.called != "" {
		t.Fatalf("handler called = %q, want no handler invocation", store.called)
	}
}
