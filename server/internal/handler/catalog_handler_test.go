package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/solenfw/calo-hub/internal/dto"
	"github.com/solenfw/calo-hub/internal/models"
	"github.com/solenfw/calo-hub/internal/repository"
)

type fakeCatalogStore struct {
	getProductByCode                  func(context.Context, string) (models.Products, error)
	getProductsByTerms                func(context.Context, []string) ([]models.Products, error)
	getImagesByCode                   func(context.Context, string) (models.Images, error)
	getAllMartinReportLists           func(context.Context) ([]models.MartinReportList, error)
	getMartinReportListByName         func(context.Context, string) (models.MartinReportList, error)
	getMartinReportProductsByReportID func(context.Context, int) ([]models.MartinReportProduct, error)
}

func (f *fakeCatalogStore) GetProductByCode(ctx context.Context, code string) (models.Products, error) {
	if f.getProductByCode != nil {
		return f.getProductByCode(ctx, code)
	}
	return models.Products{}, nil
}

func (f *fakeCatalogStore) GetProductsByTerms(ctx context.Context, terms []string) ([]models.Products, error) {
	if f.getProductsByTerms != nil {
		return f.getProductsByTerms(ctx, terms)
	}
	return nil, nil
}

func (f *fakeCatalogStore) GetImagesByCode(ctx context.Context, code string) (models.Images, error) {
	if f.getImagesByCode != nil {
		return f.getImagesByCode(ctx, code)
	}
	return models.Images{}, nil
}

func (f *fakeCatalogStore) GetAllMartinReportLists(ctx context.Context) ([]models.MartinReportList, error) {
	if f.getAllMartinReportLists != nil {
		return f.getAllMartinReportLists(ctx)
	}
	return nil, nil
}

func (f *fakeCatalogStore) GetMartinReportListByName(ctx context.Context, name string) (models.MartinReportList, error) {
	if f.getMartinReportListByName != nil {
		return f.getMartinReportListByName(ctx, name)
	}
	return models.MartinReportList{}, nil
}

func (f *fakeCatalogStore) GetMartinReportProductsByReportID(ctx context.Context, reportID int) ([]models.MartinReportProduct, error) {
	if f.getMartinReportProductsByReportID != nil {
		return f.getMartinReportProductsByReportID(ctx, reportID)
	}
	return nil, nil
}

func requestWithRouteParams(method, target string, params map[string]string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func assertJSONSuccess(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, want any) {
	t.Helper()

	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d", recorder.Code, wantStatus)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var got any
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("encode expected response: %v", err)
	}
	var wantValue any
	if err := json.Unmarshal(wantJSON, &wantValue); err != nil {
		t.Fatalf("decode expected response: %v", err)
	}
	if !reflect.DeepEqual(got, wantValue) {
		t.Fatalf("body = %s, want %s", recorder.Body.String(), wantJSON)
	}
}

func assertJSONError(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, wantMessage string) {
	t.Helper()

	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d", recorder.Code, wantStatus)
	}
	var got map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	want := map[string]string{"error": wantMessage}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("body = %v, want %v", got, want)
	}
}

func TestCatalogHandlerGetProductByCode(t *testing.T) {
	wantProduct := dto.ProductResponse{
		Code:        "AA804R",
		Eng:         "Ruler",
		Viet:        "Thước đo",
		Alternative: "17-412-15-07",
		Brand:       "B-Braun",
	}

	tests := []struct {
		name         string
		code         string
		product      models.Products
		storeError   error
		wantStatus   int
		wantError    string
		wantBody     any
		wantEmptyBody bool
	}{
		{
			name:       "missing code",
			wantStatus: http.StatusBadRequest,
			wantError:  "empty input argument",
		},
		{
			name:          "invalid code",
			code:          "invalid-code",
			wantStatus:    http.StatusOK,
			wantEmptyBody: true,
		},
		{
			name: "valid code returns DTO",
			code: "AA804R",
			product: models.Products{
				Code:        wantProduct.Code,
				Eng:         wantProduct.Eng,
				Viet:        wantProduct.Viet,
				Alternative: wantProduct.Alternative,
				Brand:       wantProduct.Brand,
			},
			wantStatus: http.StatusOK,
			wantBody:   wantProduct,
		},
		{
			name:       "not found",
			code:       "AA804R",
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "unexpected store error",
			code:       "AA804R",
			storeError: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogStore{
				getProductByCode: func(_ context.Context, code string) (models.Products, error) {
					if code != tt.code && tt.code != "" {
						t.Errorf("store code = %q, want %q", code, tt.code)
					}
					return tt.product, tt.storeError
				},
			}
			h := NewCatalogHandler(store)
			recorder := httptest.NewRecorder()
			h.GetProductByCode(recorder, requestWithRouteParams(http.MethodGet, "/catalog/products/"+tt.code, map[string]string{"code": tt.code}))

			if tt.wantEmptyBody {
				if recorder.Code != tt.wantStatus {
					t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
				}
				if recorder.Body.Len() != 0 {
					t.Fatalf("body = %q, want empty response", recorder.Body.String())
				}
				return
			}
			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestCatalogHandlerSearchProducts(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		store      []models.Products
		storeError error
		wantTerms  []string
		wantStatus int
		wantError  string
		wantBody   []dto.ProductResponse
	}{
		{
			name:       "q parameter",
			query:      "/catalog/products?q=knife",
			wantTerms:  []string{"knife"},
			store:      []models.Products{{Code: "AA804R", Eng: "Knife"}},
			wantStatus: http.StatusOK,
			wantBody:   []dto.ProductResponse{{Code: "AA804R", Eng: "Knife"}},
		},
		{
			name:       "search fallback",
			query:      "/catalog/products?q=&search=ruler",
			wantTerms:  []string{"ruler"},
			store:      []models.Products{{Code: "AA805R", Eng: "Ruler"}},
			wantStatus: http.StatusOK,
			wantBody:   []dto.ProductResponse{{Code: "AA805R", Eng: "Ruler"}},
		},
		{
			name:       "both empty",
			query:      "/catalog/products?q=&search=",
			wantStatus: http.StatusBadRequest,
			wantError:  "empty search query",
		},
		{
			name:       "multi word terms",
			query:      "/catalog/products?q=%20%20heart%20%20valve%20%20",
			wantTerms:  []string{"heart", "valve"},
			store:      []models.Products{{Code: "AA806R", Eng: "Heart valve"}},
			wantStatus: http.StatusOK,
			wantBody:   []dto.ProductResponse{{Code: "AA806R", Eng: "Heart valve"}},
		},
		{
			name:       "zero matches",
			query:      "/catalog/products?q=missing",
			store:      []models.Products{},
			wantTerms:  []string{"missing"},
			wantStatus: http.StatusOK,
			wantBody:   []dto.ProductResponse{},
		},
		{
			name:       "not found error",
			query:      "/catalog/products?q=missing",
			storeError: repository.ErrNotFound,
			wantTerms:  []string{"missing"},
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "unexpected store error",
			query:      "/catalog/products?q=knife",
			storeError: errors.New("database unavailable"),
			wantTerms:  []string{"knife"},
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogStore{
				getProductsByTerms: func(_ context.Context, terms []string) ([]models.Products, error) {
					if !reflect.DeepEqual(terms, tt.wantTerms) {
						t.Errorf("terms = %#v, want %#v", terms, tt.wantTerms)
					}
					return tt.store, tt.storeError
				},
			}
			h := NewCatalogHandler(store)
			recorder := httptest.NewRecorder()
			h.SearchProducts(recorder, httptest.NewRequest(http.MethodGet, tt.query, nil))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestCatalogHandlerGetImages(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		images       []string
		storeError   error
		wantStatus   int
		wantError    string
		wantBody     dto.ImageResponse
		wantEmptyBody bool
	}{
		{
			name:       "missing code",
			wantStatus: http.StatusBadRequest,
			wantError:  "product code is required",
		},
		{
			name:       "images returned",
			code:       "AA804R",
			images:     []string{"https://example.test/one", "https://example.test/two"},
			wantStatus: http.StatusOK,
			wantBody:   dto.ImageResponse{Code: "AA804R", Images: []string{"https://example.test/one", "https://example.test/two"}},
		},
		{
			name:       "nil images become empty array",
			code:       "AA804R",
			images:     nil,
			wantStatus: http.StatusOK,
			wantBody:   dto.ImageResponse{Code: "AA804R", Images: []string{}},
		},
		{
			name:          "not found",
			code:          "AA804R",
			storeError:    repository.ErrNotFound,
			wantStatus:    http.StatusOK,
			wantEmptyBody: true,
		},
		{
			name:          "unexpected store error",
			code:          "AA804R",
			storeError:    errors.New("database unavailable"),
			wantStatus:    http.StatusOK,
			wantEmptyBody: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogStore{
				getImagesByCode: func(_ context.Context, code string) (models.Images, error) {
					if tt.code != "" && code != tt.code {
						t.Errorf("store code = %q, want %q", code, tt.code)
					}
					return models.Images{Code: code, Images: tt.images}, tt.storeError
				},
			}
			h := NewCatalogHandler(store)
			recorder := httptest.NewRecorder()
			h.GetImages(recorder, requestWithRouteParams(http.MethodGet, "/catalog/images/"+tt.code, map[string]string{"code": tt.code}))

			if tt.wantEmptyBody {
				if recorder.Code != tt.wantStatus {
					t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
				}
				if recorder.Body.Len() != 0 {
					t.Fatalf("body = %q, want empty response", recorder.Body.String())
				}
				return
			}
			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestCatalogHandlerListMartinReports(t *testing.T) {
	tests := []struct {
		name       string
		reports    []models.MartinReportList
		storeError error
		wantStatus int
		wantError  string
		wantBody   []dto.MartinReportListResponse
	}{
		{
			name:       "multiple reports preserve order",
			reports:    []models.MartinReportList{{Name: "Report A"}, {Name: "Report B"}},
			wantStatus: http.StatusOK,
			wantBody:   []dto.MartinReportListResponse{{Name: "Report A"}, {Name: "Report B"}},
		},
		{
			name:       "zero reports become empty array",
			reports:    []models.MartinReportList{},
			wantStatus: http.StatusOK,
			wantBody:   []dto.MartinReportListResponse{},
		},
		{
			name:       "not found error",
			storeError: repository.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "not found",
		},
		{
			name:       "unexpected store error",
			storeError: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeCatalogStore{
				getAllMartinReportLists: func(_ context.Context) ([]models.MartinReportList, error) {
					return tt.reports, tt.storeError
				},
			}
			h := NewCatalogHandler(store)
			recorder := httptest.NewRecorder()
			h.ListMartinReports(recorder, httptest.NewRequest(http.MethodGet, "/catalog/report/martin", nil))

			if tt.wantError != "" {
				assertJSONError(t, recorder, tt.wantStatus, tt.wantError)
				return
			}
			assertJSONSuccess(t, recorder, tt.wantStatus, tt.wantBody)
		})
	}
}

