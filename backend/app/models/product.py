from datetime import datetime

from sqlalchemy import Column, DateTime, Float, Integer, String

from ..db.database import Base

    

class AesculapProduct(Base):
    __tablename__ = "aesculap"
    id = Column(Integer, primary_key=True, index=True)
    created_at = Column(DateTime, default=datetime.now)
    code = Column(String(20), nullable=False, index=True)
    eng_desc = Column(String(1024), nullable=False)
    viet_desc = Column(String(1024), nullable=False)
    image = Column(String(50), nullable=True)
    alternative_code = Column(String(20), nullable=True)
    brand = Column(String(20), nullable=False, index=True)


class KLSProduct(Base):
    __tablename__ = "kls_martin"
    id = Column(Integer, primary_key=True, index=True)
    created_at = Column(DateTime, default=datetime.now)
    code = Column(String(20), nullable=False, index=True)
    eng_desc = Column(String(1024), nullable=False)
    viet_desc = Column(String(1024), nullable=False)
    brand = Column(String(20), nullable=False, index=True)