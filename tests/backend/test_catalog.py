import os
import sys
import unittest
from pathlib import Path

os.environ.setdefault("DATABASE_URL", "sqlite+pysqlite:///:memory:")

ROOT = Path(__file__).resolve().parents[2]
BACKEND = ROOT / "backend"
if str(BACKEND) not in sys.path:
    sys.path.insert(0, str(BACKEND))

from sqlalchemy import create_engine
from sqlalchemy.pool import StaticPool
from sqlalchemy.orm import sessionmaker

from backend.app.api.catalog import (
    get_kls_image as route_get_kls_image,
    search_aesculap as route_search_aesculap,
    search_kls as route_search_kls,
)
from backend.app.crud.crud_catalog import search_aesculap_products, search_kls_products
from backend.app.db.database import Base
from backend.app.models.product import AesculapProduct, KLSProduct, KLSProductImage
from backend.app.schemas.product import AesculapResponse, KLSImageResponse


class CatalogTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.engine = create_engine(
            "sqlite+pysqlite:///:memory:",
            connect_args={"check_same_thread": False},
            poolclass=StaticPool,
            future=True,
        )
        cls.SessionLocal = sessionmaker(
            autocommit=False,
            autoflush=False,
            bind=cls.engine,
            future=True,
        )
        Base.metadata.create_all(bind=cls.engine)

    @classmethod
    def tearDownClass(cls) -> None:
        Base.metadata.drop_all(bind=cls.engine)
        cls.engine.dispose()

    def setUp(self) -> None:
        self.db = self.SessionLocal()
        for model in (KLSProductImage, AesculapProduct, KLSProduct):
            self.db.query(model).delete()
        self.db.commit()

    def tearDown(self) -> None:
        self.db.close()

    def add_kls(self, code: str, eng_desc: str, viet_desc: str = "mo ta") -> None:
        self.db.add(
            KLSProduct(
                code=code,
                eng_desc=eng_desc,
                viet_desc=viet_desc,
                brand="Martin",
            )
        )
        self.db.commit()

    def add_aesculap(
        self,
        code: str,
        eng_desc: str,
        viet_desc: str = "mo ta",
        alternative_code: str | None = None,
    ) -> None:
        self.db.add(
            AesculapProduct(
                code=code,
                eng_desc=eng_desc,
                viet_desc=viet_desc,
                image=f"{code}.jpg",
                alternative_code=alternative_code,
                brand="B-Braun",
            )
        )
        self.db.commit()


class KLSCatalogSearchTests(CatalogTestCase):
    def test_exact_kls_code_search_returns_only_the_matching_product(self) -> None:
        self.add_kls("12-345-67-89", "Straight surgical forceps")
        self.add_kls("98-765-43-21", "Curved surgical forceps")

        results = search_kls_products(self.db, "12-345-67-89")

        self.assertEqual(["12-345-67-89"], [product.code for product in results])

    def test_kls_description_search_requires_every_term_and_respects_limit(self) -> None:
        self.add_kls("11-111-11-11", "Straight surgical forceps")
        self.add_kls("22-222-22-22", "Straight surgical clamp")
        self.add_kls("33-333-33-33", "Fine straight forceps")

        results = search_kls_products(self.db, "straight forceps", limit=1)

        self.assertEqual(1, len(results))
        self.assertIn(results[0].code, {"11-111-11-11", "33-333-33-33"})

    def test_kls_api_returns_empty_list_for_no_matches(self) -> None:
        response = route_search_kls(q="missing", db=self.db)

        self.assertEqual([], response)


class AesculapCatalogSearchTests(CatalogTestCase):
    def test_aesculap_search_matches_code_or_description_terms(self) -> None:
        self.add_aesculap("AB123", "Bone holding forceps", alternative_code="ALT-1")
        self.add_aesculap("CD456", "Needle holder")

        code_results = search_aesculap_products(self.db, "AB123")
        description_results = search_aesculap_products(self.db, "needle holder")

        self.assertEqual(["AB123"], [product.code for product in code_results])
        self.assertEqual(["CD456"], [product.code for product in description_results])

    def test_aesculap_api_serializes_optional_product_fields(self) -> None:
        self.add_aesculap("AB123", "Bone holding forceps", alternative_code="ALT-1")

        response = route_search_aesculap(q="bone", db=self.db)
        serialized = [AesculapResponse.model_validate(item).model_dump() for item in response]

        self.assertEqual(
            [
                {
                    "code": "AB123",
                    "eng_desc": "Bone holding forceps",
                    "viet_desc": "mo ta",
                    "brand": "B-Braun",
                    "image": "AB123.jpg",
                    "alternative_code": "ALT-1",
                }
            ],
            serialized,
        )


class KLSImageTests(CatalogTestCase):
    def test_image_endpoint_returns_stored_urls(self) -> None:
        self.db.add(
            KLSProductImage(
                code="12-345-67-89",
                img1_url="https://example.test/one.jpg",
                img2_url="https://example.test/two.jpg",
                img3_url=None,
            )
        )
        self.db.commit()

        response = route_get_kls_image(code="12-345-67-89", db=self.db)
        serialized = KLSImageResponse.model_validate(response).model_dump()

        self.assertEqual(
            {
                "code": "12-345-67-89",
                "img1_url": "https://example.test/one.jpg",
                "img2_url": "https://example.test/two.jpg",
                "img3_url": None,
            },
            serialized,
        )

    def test_image_endpoint_returns_null_urls_for_missing_code(self) -> None:
        response = route_get_kls_image(code="NOPE", db=self.db)

        self.assertEqual(
            {
                "code": "NOPE",
                "img1_url": None,
                "img2_url": None,
                "img3_url": None,
            },
            response,
        )


if __name__ == "__main__":
    unittest.main()
