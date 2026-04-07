from sqlalchemy.orm import Session

from ..models.product import Product
from ..schemas.product import ProductCreate


def get_products(db: Session, skip: int = 0, limit: int = 50) -> list[Product]:
    return db.query(Product).offset(skip).limit(limit).all()


def get_product(db: Session, product_id: int) -> Product | None:
    return db.query(Product).filter(Product.id == product_id).first()


def create_product(db: Session, product_in: ProductCreate) -> Product:
    product = Product(
        name=product_in.name,
        description=product_in.description,
        price=product_in.price,
    )
    db.add(product)
    db.commit()
    db.refresh(product)
    return product
