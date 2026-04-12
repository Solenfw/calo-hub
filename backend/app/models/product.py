from datetime import datetime

from sqlalchemy import Column, DateTime, Float, Integer, String

from ..db.database import Base

    

class AesculapProduct(Base):
    __tablename__ = "aesculap"
    id = Column(Integer, primary_key=True, index=True)
    code = Column(String(20), nullable=False, index=True)
    eng_desc = Column(String(1024), nullable=False)
    eng_desc = Column(String(1024), nullable=False)
    image = Column(Float, nullable=True)
    alternative_code = Column(String(20), nullable=True)
    created_at = Column(DateTime, default=datetime.now(datetime.timezone.utc))


class KLSProduct(Base):
    __tablename__ = "kls_martin"
    id = Column(Integer, primary_key=True, index=True)
    code = Column(String(20), nullable=False, index=True)
    eng_desc = Column(String(1024), nullable=False)
    viet_desc = Column(String(1024), nullable=False)
    created_at = Column(DateTime, default=datetime.now(datetime.timezone.utc))