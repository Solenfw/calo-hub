from pydantic import BaseModel


class ProductBase(BaseModel):
    code: str
    eng_desc: str
    viet_desc: str

    class Config:
        orm_mode = True

class KLSResponse(ProductBase):
    pass

class AesculapResponse(ProductBase):
    image: str | None = None
    alternative_code: str | None = None