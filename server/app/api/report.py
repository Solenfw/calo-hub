from operator import or_
from typing import List
from io import BytesIO

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from reportlab import canvas
from reportlab.lib.pagesizes import A4

from server.pdf_report_generator import draw_summary_table_header

from .dependencies import get_db
from ..schemas.product import KLSReportResponse
from ..services.pdf_drawer import draw_header, draw_product_page, draw_summary_table_header, draw_summary_row, prefetch_images

router = APIRouter()

@router.get("/kls/report/", response_model=List[KLSReportResponse])
def get_kls_report(db: Session = Depends(get_db)) -> bytes:
   products = List[KLSReportResponse]
   buffer = BytesIO()
   c = canvas.Canvas(buffer, pagesize=A4)
   image_cache = prefetch_images(products)

   print("Drawing product pages...")
   y = draw_header(c, "BTM", "KLS Martin Products")
   for p in products:
      y = draw_product_page(c, p, y, image_cache=image_cache)

   c.showPage()

   print("Drawing summary page...")
   y = draw_header(c, "BTM", "Summary")
   y = draw_summary_table_header(c, y)
   for p in products:
      y = draw_summary_row(c, p, y)

   c.save()
   buffer.seek(0)
   return buffer.getvalue()
