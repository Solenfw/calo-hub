from operator import or_
import re
from typing import List

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from .dependencies import get_db
from ..crud.crud_catalog import search_aesculap_products, search_kls_products
from ..schemas.product import KLSResponse, AesculapResponse

router = APIRouter()



@router.get("/kls/", response_model=List[KLSResponse])
def search_kls(q: str, db: Session = Depends(get_db), limit: int = 50):
    products = search_kls_products(db, search_term=q, limit=limit)
    if not products:
        # Returning an empty list is better API design than a 404 for a search query
        return [] 
    return products


@router.get("/aes/", response_model=List[AesculapResponse])
def search_aesculap(q: str, db: Session = Depends(get_db), limit: int = 50):
    products = search_aesculap_products(db, search_term=q, limit=limit)
    if not products:
        # Returning an empty list is better API design than a 404 for a search query
        return []
    return products