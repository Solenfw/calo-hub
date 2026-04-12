from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .api import catalog
from .core.config import settings

app = FastAPI(title="Calo Hub Backend")

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.frontend_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(catalog.router, prefix="/catalog", tags=["catalog"])
