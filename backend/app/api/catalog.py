from typing import List

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from ..api.dependencies import get_db
from ..crud.crud_catalog import get_product, get_products
from ..schemas.product import ProductResponse

router = APIRouter()


@router.get("/", response_model=List[ProductResponse])
def list_products(skip: int = 0, limit: int = 50, db: Session = Depends(get_db)):
    return get_products(db, skip=skip, limit=limit)


@router.get("/{product_id}", response_model=ProductResponse)
def read_product(product_id: int, db: Session = Depends(get_db)):
    product = get_product(db, product_id)
    if product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return product
