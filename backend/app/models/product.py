from sqlalchemy import Column, String
from ..db.database import Base

    

class AesculapProduct(Base):
    __tablename__ = "aesculap"
    code = Column(String(20), primary_key=True, index=True)
    eng_desc = Column(String(1024), nullable=False)
    viet_desc = Column(String(1024), nullable=False)
    image = Column(String(50), nullable=True)
    alternative_code = Column(String(20), nullable=True)
    brand = Column(String(20), nullable=False, index=True)


class KLSProduct(Base):
    __tablename__ = "kls_martin"
    code = Column(String(20), primary_key=True, index=True)
    eng_desc = Column(String(1024), nullable=False)
    viet_desc = Column(String(1024), nullable=False)
    alternative_code = Column(String(20), nullable=True)
    brand = Column(String(20), nullable=False, index=True)
    

class KLSProductImage(Base):
    __tablename__ = "martin_images"
    code = Column(String(20), primary_key=True, index=True)
    img1_url = Column(String(50), nullable=False)
    img2_url = Column(String(50), nullable=True)
    img3_url = Column(String(50), nullable=True)