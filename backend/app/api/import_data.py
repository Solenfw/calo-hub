from fastapi import APIRouter, Depends, File, HTTPException, UploadFile
from sqlalchemy.orm import Session

from ..api.dependencies import get_db
from ..crud.crud_catalog import create_product
from ..schemas.product import ProductCreate
from ..services.csv_parser import parse_csv

router = APIRouter()


@router.post("/upload")
async def upload_csv(file: UploadFile = File(...), db: Session = Depends(get_db)):
    if file.content_type not in ["text/csv", "application/vnd.ms-excel"]:
        raise HTTPException(status_code=400, detail="Invalid upload format, expected CSV")

    payload = (await file.read()).decode("utf-8")
    rows = parse_csv(payload)
    imported = []
    for row in rows:
        product_in = ProductCreate(**row)
        imported.append(create_product(db, product_in))

    return {"imported": len(imported)}
