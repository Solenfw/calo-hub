from sqlalchemy import Column, Integer, String, text
from ..db.database import Base

    

class AesculapProduct(Base):
    __tablename__ = "aesculap"
    code = Column(String(20), primary_key=True, index=True)
    eng = Column(String(1024), nullable=False)
    viet = Column(String(1024), nullable=False)
    image = Column(String(2000), nullable=True)
    alternative = Column(String(20), nullable=True)
    brand = Column(String(20), nullable=False, index=True)


class KLSProduct(Base):
    __tablename__ = "kls_martin"
    code = Column(String(20), primary_key=True, index=True)
    eng = Column(String(1024), nullable=False)
    viet = Column(String(1024), nullable=False)
    alternative = Column(String(20), nullable=True)
    brand = Column(String(20), nullable=False, index=True)
    

class KLSProductImage(Base):
    __tablename__ = "martin_images"
    code = Column(String(20), primary_key=True, index=True)
    img1_url = Column(String(2000), nullable=False)
    img2_url = Column(String(2000), nullable=True)
    img3_url = Column(String(2000), nullable=True)


class KLSProductForReport(Base):
    __tablename__ = "klsmartin_report"
    no = Column(Integer, primary_key=True, index=True)
    code = Column(String(20), nullable=False)
    description = Column(String(1024), nullable=False)
    img_url = Column(String(2000), nullable=True)
    quantity = Column(Integer, nullable=False)