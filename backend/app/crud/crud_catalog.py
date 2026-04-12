import re

from sqlalchemy import or_
from sqlalchemy.orm import Session
from ..models.product import KLSProduct, AesculapProduct



def search_kls_products(db: Session, search_term: str, limit: int = 50 ) -> list[KLSProduct]:
    # Check if exact code
    if re.match(r'(\d{2}-\d{3}-\d{2}-\d{2})', search_term):
        return db.query(KLSProduct).filter(KLSProduct.code == search_term).all()
    
    # Otherwise return a list of matching descriptions
    return db.query(KLSProduct).filter(
        or_(
            KLSProduct.viet_desc.ilike(f"%{search_term}%"),
            KLSProduct.eng_desc.ilike(f"%{search_term}%")
        )
    ).limit(limit).all()
    

def search_aesculap_products(db: Session, search_term: str, limit: int = 50 ) -> list[AesculapProduct]:
    return db.query(AesculapProduct).filter(
        or_(
            AesculapProduct.code == search_term,
            AesculapProduct.viet_desc.ilike(f"%{search_term}%"),
            AesculapProduct.eng_desc.ilike(f"%{search_term}%")
        )
    ).limit(limit).all()
    
