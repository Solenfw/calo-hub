from sqlalchemy import create_engine
from sqlalchemy.orm import declarative_base, sessionmaker

from ..core.config import settings

engine = create_engine(settings.database_url, future=True, connect_args={"sslmode": "require"} if "localhost" not in settings.database_url else {})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine, future=True)
Base = declarative_base()
