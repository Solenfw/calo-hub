import re

from sqlalchemy import or_
from sqlalchemy.orm import Session
from ..models.product import KLSProduct, AesculapProduct


def get_kls_products(db: Session, skip: int = 0, limit: int = 50) -> list[KLSProduct]:
    return db.query(KLSProduct).offset(skip).limit(limit).all()


def get_aesculap_products(db: Session, skip: int = 0, limit: int = 50) -> list[AesculapProduct]:
    return db.query(AesculapProduct).offset(skip).limit(limit).all()


def get_kls_product(db: Session, search_term: str) -> KLSProduct | None:
    if re.match(r'(\d{2}-\d{3}-\d{2}-\d{2})', search_term):
        return db.query(KLSProduct).filter(KLSProduct.code == search_term).first()
    
    return db.query(KLSProduct).filter(
        or_(
            KLSProduct.viet_desc.ilike(f"%{search_term}%"),
            KLSProduct.eng_desc.ilike(f"%{search_term}%")
        )
    ).first()
    

def get_aesculap_product(db: Session, search_term: str) -> AesculapProduct | None:
    return db.query(AesculapProduct).filter(
        or_(
            AesculapProduct.code == search_term,
            AesculapProduct.viet_desc.ilike(f"%{search_term}%"),
            AesculapProduct.eng_desc.ilike(f"%{search_term}%")
        )
    ).first()
    
