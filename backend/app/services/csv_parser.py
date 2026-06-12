import csv
from io import StringIO
import sys, os
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))


"""
CSV parsing service for data ingestion. Expects CSVs with the following columns:
    [code, eng, viet, image, compatible_kls_code] (Aesculap).
    [code, eng, viet] (KLS Martin).
"""
def parse_csv(payload: str) -> list[dict[str, str | float]]:
    reader = csv.DictReader(StringIO(payload))
    rows: list[dict[str, str | float]] = []
    for row in reader:
        if not row.get("name") or not row.get("price"):
            continue
        rows.append(
            {
                "name": row["name"].strip(),
                "description": row.get("description", "").strip(),
                "price": float(row["price"]),
            }
        )
    return rows
