from pydantic import BaseModel, ConfigDict


class ProductBase(BaseModel):
    code: str
    eng_desc: str
    viet_desc: str
    brand: str
    
    model_config = ConfigDict(from_attributes=True) 


class KLSResponse(ProductBase):
    pass

class AesculapResponse(ProductBase):
    image: str | None = None
    alternative_code: str | None = None

class KLSImageResponse(BaseModel):
    code: str
    img1_url: str
    img2_url: str | None = None
    img3_url: str | None = None
    
    model_config = ConfigDict(from_attributes=True)