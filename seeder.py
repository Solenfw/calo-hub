import os
import sqlalchemy
import csv
from dotenv import load_dotenv

load_dotenv()


def get_val(row, index):
    return row[index] if len(row) > index and row[index] != '' else None


def main():
    db_url = os.getenv("DATABASE_URL")
    engine = sqlalchemy.create_engine(db_url)

    with engine.connect() as connection:
        with connection.begin():

            # -------- 1. kls_martin (PARENT FIRST) --------
            with open('data/klsmartin.csv', 'r', newline='') as f:
                reader = csv.reader(f)
                next(reader, None)

                for row in reader:
                    if not row or not row[0]:
                        continue

                    connection.execute(
                        sqlalchemy.text("""
                            INSERT INTO kls_martin (code, eng_desc, viet_desc, alternative_code)
                            VALUES (:code, :eng, :viet, :alternative)
                            ON CONFLICT (code) DO NOTHING
                        """),
                        {
                            "code": row[0],
                            "eng": row[1],
                            "viet": row[2],
                            "alternative": row[3]
                        }
                    )

            # -------- 2. martin_images (CHILD) --------
            with open('data/martin_images.csv', 'r', newline='') as f:
                reader = csv.reader(f)

                for row in reader:
                    if not row or not row[0]:
                        continue

                    connection.execute(
                        sqlalchemy.text("""
                            INSERT INTO martin_images 
                            (code, img1_url, img2_url, img3_url)
                            VALUES (:code, :img1_url, :img2_url, :img3_url)
                            ON CONFLICT (code) DO UPDATE SET
                                img1_url = EXCLUDED.img1_url,
                                img2_url = EXCLUDED.img2_url,
                                img3_url = EXCLUDED.img3_url
                        """),
                        {
                            "code": row[0],
                            "img1_url": get_val(row, 1),
                            "img2_url": get_val(row, 2),
                            "img3_url": get_val(row, 3),
                        }
                    )

            # -------- 3. aesculap --------
            with open('data/aesculap.csv', 'r', newline='') as f:
                reader = csv.reader(f)
                next(reader, None)

                for row in reader:
                    if not row or not row[0]:
                        continue

                    connection.execute(
                        sqlalchemy.text("""
                            INSERT INTO aesculap 
                            (code, eng_desc, viet_desc, image, alternative_code)
                            VALUES (:code, :eng, :viet, :img, :alt)
                            ON CONFLICT (code) DO NOTHING
                        """),
                        {
                            "code": row[0],
                            "eng": row[1],
                            "viet": row[2],
                            "img": get_val(row, 3),
                            "alt": get_val(row, 4),
                        }
                    )


if __name__ == "__main__":
    main()