/**
 * KLS Product Handler
 * AESC Product Handler
 * This file contains the product handlers for the KLS and AESC products. 
 * Output: a product list with detailed properties.
*/


const KLS_BASE_URL = "https://www.klsmartin.com/shop/en/products/?eID=catalog-search";
const AESC_BASE_URL = "https://surgical-instruments.bbraun.com/api/occ/v2/bbraunb2b/materialSearch?viewId=en_01&salesAreaId=01_B2C_API_AESI&customerNumber=aesculap_client_customer&pageSize=50&pageNumber=0&materialTypes=ARTICLE&sortQuery=score-desc";

export async function getKLSProducts(searchQuery: string) {
    try {
        const response = await fetch(KLS_BASE_URL + `&text=${searchQuery}&lang=en-GB&pageUid=2992`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error fetching KLS products:", error);
        return [];
    }
}

export async function getAESCProducts(searchQuery: string) {
    try {
        const response = await fetch(AESC_BASE_URL + `&text=${searchQuery}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error fetching AESC products:", error);
        return [];
    }
}