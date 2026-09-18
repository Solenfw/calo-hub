import {
    ApiErrorResponse,
    ImageResponse,
    ProductResponse,
} from '@/types';

const API_BASE_URL = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080').replace(/\/+$/, '');

async function request<T>(path: string, init?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`, init);
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
