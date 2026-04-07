import csv
from io import StringIO


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
