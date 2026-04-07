from pydantic import BaseModel


class ProductBase(BaseModel):
    name: str
    description: str | None = None
    price: float

    class Config:
        orm_mode = True


class ProductCreate(ProductBase):
    pass


class ProductResponse(ProductBase):
    id: int
