import httpx


async def fetch_brand_data(endpoint: str) -> dict:
    async with httpx.AsyncClient(timeout=10.0) as client:
        response = await client.get(endpoint)
        response.raise_for_status()
        return response.json()
