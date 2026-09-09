
/** Matches server/internal/dto.ProductResponse. */
export interface ProductResponse {
    code: string;
    viet: string;
    eng: string;
    alternative: string;
    brand: string;
}

/** Matches server/internal/dto.ImageResponse. */
export interface ImageResponse {
    code: string;
    images: string[];
}

/** Matches server/internal/dto.MartinReportListResponse. */
export interface MartinReportListResponse {
    name: string;
}

/** Matches server/internal/dto.MartinReportProductResponse. */
export interface MartinReportProductResponse {
    row_no: number;
    code: string;
    eng: string;
    image: string | null;
    quantity: number;
}

/** Matches the error payload written by the server handlers. */
export interface ApiErrorResponse {
    error: string;
}
