"""
   This script will generate a report in XLSX format based on the equivalent CSV file,
   and update the Supabase database with the latest data.
"""
import os
import sys
import csv
import subprocess
from openpyxl import Workbook
from sqlalchemy import create_engine, MetaData, Table, inspect
from sqlalchemy.dialects.postgresql import insert
from dotenv import load_dotenv

# Load local environment variables
load_dotenv()

# SQLAlchemy connection URI
SUPABASE_DATABASE_URL = os.environ.get("SUPABASE_DATABASE_URL")

# Create SQLAlchemy connection engine
engine = create_engine(SUPABASE_DATABASE_URL) if SUPABASE_DATABASE_URL else None
metadata = MetaData()

def get_table_name(csv_path):
    """Derives database table name from file name."""
    filename = os.path.basename(csv_path)
    table_name, _ = os.path.splitext(filename)
    return table_name

def get_changed_csv_files():
    """Queries git for modified/created CSVs."""
    try:
        modified = subprocess.check_output(["git", "diff", "--name-only"]).decode("utf-8").splitlines()
        status = subprocess.check_output(["git", "status", "--porcelain"]).decode("utf-8").splitlines()
        untracked = [line[3:] for line in status if line.startswith("?? ")]
        
        if not modified and not untracked:
            try:
                modified = subprocess.check_output(["git", "diff", "--name-only", "HEAD~1", "HEAD"]).decode("utf-8").splitlines()
            except subprocess.CalledProcessError:
                pass

        all_changed = list(set(modified + untracked))
        return [f for f in all_changed if f.endswith(".csv") and os.path.exists(f)]
    except Exception as e:
        print(f"Error querying Git: {e}")
        return []

def write_xlsx_report(records, file_object):
    """Writes list of dictionaries to Excel workbook directly using openpyxl."""
    wb = Workbook()
    ws = wb.active
    ws.title = "Report"
    
    if not records:
        wb.save(file_object)
        return

    # Extract headers from the keys of the first row dictionary
    headers = list(records[0].keys())
    ws.append(headers)

    # Append row values matching headers
    for row in records:
        ws.append([row[h] for h in headers])

    wb.save(file_object)

def sync_single_file(csv_path):
    table_name = get_table_name(csv_path)
    xlsx_path = f"C:/Users/Admin/Documents/{table_name}.xlsx"
    
    print(f"\n--- Syncing: {csv_path} ---")
    
    # 1. Read CSV using native 'with open()' and DictReader
    records = []
    with open(csv_path, mode='r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            # Map empty CSV values to None (converts to SQL NULL)
            cleaned_row = {k: (None if v == "" else v) for k, v in row.items()}
            records.append(cleaned_row)
            
    if not records:
        print("CSV is empty. Skipping.")
        return

    # 2. Write Excel File using with open() and openpyxl
    print(f"Writing Excel: {xlsx_path}")
    os.makedirs("reports", exist_ok=True)
    with open(xlsx_path, mode='wb') as f:
        write_xlsx_report(records, f)
    print("✔ Excel report updated.")
        
    # 3. Database Upsert using SQLAlchemy
    print(f"Upserting to database table: '{table_name}'")
    try:
        # Reflect table structure from database
        table = Table(table_name, metadata, autoload_with=engine)
        
        # Get primary keys for the upsert target
        primary_keys = [key.name for key in inspect(table).primary_key]
        if not primary_keys:
            raise ValueError(f"Table '{table_name}' has no primary key.")

        # Build insert and conflict handling statements
        stmt = insert(table)
        update_dict = {
            c.name: c for c in stmt.excluded if not c.primary_key
        }
        upsert_stmt = stmt.on_conflict_do_update(
            index_elements=primary_keys,
            set_=update_dict
        )

        # Execute transaction
        with engine.begin() as conn:
            conn.execute(upsert_stmt, records)
        print(f"✔ Database sync completed for {len(records)} rows.")
        
    except Exception as e:
        print(f"✘ SQLAlchemy Database Error: {e}")

def main():
    if not SUPABASE_DATABASE_URL:
        print("Error: SUPABASE_DATABASE_URL environment variable is missing.")
        return
        
    if len(sys.argv) > 1:
        files_to_sync = [f for f in sys.argv[1:] if f.endswith(".csv") and os.path.exists(f)]
    else:
        print("Auto-detecting changed CSV files...")
        files_to_sync = get_changed_csv_files()
        
    if not files_to_sync:
        print("No changes found to sync.")
        return
        
    for file in files_to_sync:
        sync_single_file(file)

if __name__ == "__main__":
    main()