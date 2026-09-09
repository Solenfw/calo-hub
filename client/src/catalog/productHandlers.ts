import {
    ApiErrorResponse,
    ImageResponse,
    MartinReportListResponse,
    MartinReportProductResponse,
    ProductResponse,
} from '@/types';

const API_BASE_URL = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080').replace(/\/+$/, '');

async function request<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`);
    const payload: unknown = await response.json().catch(() => null);

    if (!response.ok) {
        const errorPayload = payload as Partial<ApiErrorResponse> | null;
        throw new Error(errorPayload?.error || `Request failed with status ${response.status}`);
    }

    return payload as T;
}

/** GET /catalog/products?q=... */
export const searchProducts = async (searchTerm: string): Promise<ProductResponse[]> => {
    try {
        return await request<ProductResponse[]>(`/catalog/products?q=${encodeURIComponent(searchTerm)}`);
    } catch (error) {
        console.error('Error searching products:', error);
        return [];
    }
};

/** GET /catalog/products/{code} */
export const getProductByCode = async (code: string): Promise<ProductResponse | null> => {
    try {
        return await request<ProductResponse>(`/catalog/products/${encodeURIComponent(code)}`);
    } catch (error) {
        console.error('Error fetching product:', error);
        return null;
    }
};

/** GET /catalog/images/{code} */
export const getImages = async (code: string): Promise<ImageResponse | null> => {
    try {
        return await request<ImageResponse>(`/catalog/images/${encodeURIComponent(code)}`);
    } catch (error) {
        console.error('Error fetching product images:', error);
        return null;
    }
};

/** GET /catalog/report/martin */
export const listMartinReports = async (): Promise<MartinReportListResponse[]> => {
    try {
        return await request<MartinReportListResponse[]>('/catalog/report/martin');
    } catch (error) {
        console.error('Error fetching Martin reports:', error);
        return [];
    }
};

/** GET /catalog/report/martin/{name} */
export const getMartinReportByName = async (name: string): Promise<MartinReportListResponse | null> => {
    try {
        return await request<MartinReportListResponse>(`/catalog/report/martin/${encodeURIComponent(name)}`);
    } catch (error) {
        console.error('Error fetching Martin report:', error);
        return null;
    }
};

/** GET /catalog/report/martin/all/{report_id} */
export const getMartinReportProducts = async (reportId: number): Promise<MartinReportProductResponse[]> => {
    try {
        return await request<MartinReportProductResponse[]>(`/catalog/report/martin/all/${reportId}`);
    } catch (error) {
        console.error('Error fetching Martin report products:', error);
        return [];
    }
};
