import os
import dotenv
from typing import List

from pydantic import BaseSettings
dotenv.load_dotenv()

class Settings(BaseSettings):
    database_url: str = os.getenv("DATABASE_URL")
    frontend_origins: List[str] = ["*"]

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
