import re

from sqlalchemy import or_, and_
from sqlalchemy.orm import Session
from ..models.product import KLSProduct, AesculapProduct, KLSProductImage



def search_kls_products(db: Session, search_term: str, limit: int = 50 ) -> list[KLSProduct]:
    # Check if exact code
    if re.match(r'(\d{2}-\d{3}-\d{2}-\d{2})', search_term):
        return db.query(KLSProduct).filter(KLSProduct.code == search_term).all()
    
    # Otherwise return a list of matching descriptions
    terms = search_term.split()
    conditions = [
        or_(
            KLSProduct.viet.ilike(f"%{term}%"),
            KLSProduct.eng.ilike(f"%{term}%")
        ) for term in terms
    ]
    return db.query(KLSProduct).filter(and_(*conditions)).limit(limit).all()


def search_martin_image (db: Session, code: str) -> KLSProductImage | None:
    return db.query(KLSProductImage).filter(KLSProductImage.code == code).first()


def search_aesculap_products(db: Session, search_term: str, limit: int = 50 ) -> list[AesculapProduct]:
    terms = search_term.split()
    conditions = [
        or_(
            AesculapProduct.code == term,
            AesculapProduct.viet.ilike(f"%{term}%"),
            AesculapProduct.eng.ilike(f"%{term}%")
        ) for term in terms
    ]
    return db.query(AesculapProduct).filter(and_(*conditions)).limit(limit).all()
