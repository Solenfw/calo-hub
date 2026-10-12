import {
   CreateMartinReportRequest,
   MartinReportListResponse,
   ReportProductResponse,
   UpdateMartinReportRequest,
   ApiErrorResponse
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

/** GET /report/martin */
export const listMartinReports = async (): Promise<MartinReportListResponse[]> => {
    try {
        return await request<MartinReportListResponse[]>('/report/martin');
    } catch (error) {
        console.error('Error fetching Martin reports:', error);
        return [];
    }
};

/** POST /report/martin */
export const createMartinReport = async (name: string): Promise<MartinReportListResponse | null> => {
    try {
        const payload: CreateMartinReportRequest = { name };
        return await request<MartinReportListResponse>('/report/martin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
    } catch (error) {
        console.error('Error creating Martin report:', error);
        return null;
    }
};

/** DELETE /report/{id} */
export const deleteMartinReport = async (reportId: number): Promise<boolean> => {
    try {
        await request<{ status: string }>(`/report/${reportId}`, {
            method: 'DELETE',
        });
        return true;
    } catch (error) {
        console.error('Error deleting Martin report:', error);
        return false;
    }
};

/** PUT /report/{id}/products */
export const updateReportProducts = async (reportId: number, products: ReportProductResponse[]): Promise<boolean> => {
    try {
        const payload: UpdateMartinReportRequest = { products };
        await request<{ status: string }>(`/report/${reportId}/products`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        return true;
    } catch (error) {
        console.error('Error updating Martin report products:', error);
        return false;
    }
};

/** GET /report/{name} */
export const getMartinReportByName = async (name: string): Promise<MartinReportListResponse | null> => {
    try {
        return await request<MartinReportListResponse>(`/report/${encodeURIComponent(name)}`);
    } catch (error) {
        console.error('Error fetching Martin report:', error);
        return null;
    }
};

/** GET /report/all/{report_id} */
export const getReportProducts = async (reportId: number): Promise<ReportProductResponse[]> => {
    try {
        return await request<ReportProductResponse[]>(`/report/all/${reportId}`);
    } catch (error) {
        console.error('Error fetching Martin report products:', error);
        return [];
    }
};

/** GET /report/template */
export const downloadImportTemplate = async (): Promise<Blob> => {
    const response = await fetch(`${API_BASE_URL}/report/template`);
    if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
    }
    return await response.blob();
};

/** GET /report/{name}/pdf */
export const exportMartinReportPdf = async (reportName: string): Promise<Blob> => {
    const response = await fetch(`${API_BASE_URL}/report/${encodeURIComponent(reportName)}/pdf`);
    if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`);
    }
    return await response.blob();
};
