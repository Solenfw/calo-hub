from typing import List

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from .dependencies import get_db
from ..crud.crud_catalog import get_aesculap_product, get_kls_products, get_kls_product, get_aesculap_product
from ..schemas.product import KLSResponse, AesculapResponse

router = APIRouter()


@router.get("/kls/", response_model=List[KLSResponse])
def get_kls_products(skip: int = 0, limit: int = 50, db: Session = Depends(get_db)):
    return get_kls_products(db, skip=skip, limit=limit)

@router.get("/kls/{search_term}", response_model=KLSResponse)
def read_kls_product(search_term: str, db: Session = Depends(get_db)):
    product = get_kls_product(db, search_term)
    if product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return product


@router.get("/aes/", response_model=List[AesculapResponse])
def get_aesculap_products(skip: int = 0, limit: int = 50, db: Session = Depends(get_db)):
    return get_aesculap_products(db, skip=skip, limit=limit)


@router.get("/aes/{search_term}", response_model=AesculapResponse)
def read_aesculap_product(search_term: str, db: Session = Depends(get_db)):
    product = get_aesculap_product(db, search_term)
    if product is None:
        raise HTTPException(status_code=404, detail="Product not found")
    return product