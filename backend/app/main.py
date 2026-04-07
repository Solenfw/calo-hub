from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .api import catalog, import_data, users
from .core.config import settings

app = FastAPI(title="Calo Hub Backend")

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.frontend_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(catalog.router, prefix="/api/catalog", tags=["catalog"])
app.include_router(import_data.router, prefix="/api/import", tags=["import"])
app.include_router(users.router, prefix="/api/users", tags=["users"])
