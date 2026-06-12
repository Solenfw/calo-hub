
export interface KLSProduct {
    code: string;
    viet: string;
    eng: string;
    alternative: string | null;
    brand: string;
}

export interface AesculapProduct {
    code: string;
    viet: string;
    eng: string;
    image: string;
    alternative: string | null;
    brand: string;
}

export interface KLSImageResponse {
    code: string;
    img1_url: string;
    img2_url: string;
    img3_url: string;
}