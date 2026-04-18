/*
    * This file contains the handlers for the product search functionality in the frontend.
    * Functions will takes in search terms, call to FastAPI 8000 port, and return the product information to be displayed on the frontend.
*/

import { KLSProduct, AesculapProduct, KLSImageResponse } from '@/types';

const API_BASE_URL = 'http://localhost:8000/catalog';

export const searchKLSProduct = async (searchTerm: string, limit: number): Promise<KLSProduct[]> => {
    try {
        const response = await fetch(`${API_BASE_URL}/kls/?q=${encodeURIComponent(searchTerm)}&limit=${limit}`);
        if (!response.ok) return [];
        return response.json();
    } catch (error) {
        console.error('Error searching KLS product:', error);
        return [];
    }
};

export const getKLSImages = async (code: string): Promise<KLSImageResponse | null> => {
    try {
        const response = await fetch(`${API_BASE_URL}/kls/images/${encodeURIComponent(code)}`);
        if (!response.ok) return null;
        const data = await response.json();
        return data as KLSImageResponse; // Returns the full object with img1_url, img2_url, img3_url
    } catch (error) {
        console.error('Error fetching KLS image:', error);
        return null;
    }
};

export const searchAesculapProduct = async (searchTerm: string, limit: number): Promise<AesculapProduct[]> => {
    try {
        const response = await fetch(`${API_BASE_URL}/aes/?q=${encodeURIComponent(searchTerm)}&limit=${limit}`);
        if (!response.ok) return [];
        return response.json();
    } catch (error) {
        console.error('Error searching Aesculap product:', error);
        return [];
    }
};
