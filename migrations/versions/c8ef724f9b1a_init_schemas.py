"""init schemas

Revision ID: c8ef724f9b1a
Revises: 
Create Date: 2026-05-29 04:20:19.206747

"""
from typing import Sequence, Union
from pathlib import Path

from alembic import op


# revision identifiers, used by Alembic.
revision: str = 'c8ef724f9b1a'
down_revision: Union[str, Sequence[str], None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None

def upgrade() -> None:
    sql = (Path(__file__).parent / "schema.sql").read_text()
    op.execute(sql)

def downgrade() -> None:
    """Downgrade schema."""
    op.execute("DROP TABLE IF EXISTS martin_images, kls_martin, aesculap CASCADE;")