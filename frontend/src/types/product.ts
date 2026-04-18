
export interface KLSProduct {
    code: string;
    viet_desc: string;
    eng_desc: string;
    brand: string;
}

export interface AesculapProduct {
    code: string;
    viet_desc: string;
    eng_desc: string;
    image: string;
    alternative_code: string;
    brand: string;
}

export interface KLSImageResponse {
    code: string;
    img1_url: string;
    img2_url: string;
    img3_url: string;
}