from pydantic import BaseModel, ConfigDict


class ProductBase(BaseModel):
    code: str
    eng: str
    viet: str
    brand: str
    
    model_config = ConfigDict(from_attributes=True) 


class KLSResponse(ProductBase):
    alternative: str | None = None
    pass

class AesculapResponse(ProductBase):
    image: str | None = None
    alternative: str | None = None

class KLSImageResponse(BaseModel):
    code: str
    img1_url: str | None = None
    img2_url: str | None = None
    img3_url: str | None = None
    
    model_config = ConfigDict(from_attributes=True)


class KLSReportResponse(BaseModel):
    no: int
    code: str
    description: str
    img_url: str | None = None
    quantity: int
    
    model_config = ConfigDict(from_attributes=True)