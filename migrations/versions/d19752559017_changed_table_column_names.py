"""changed table column names.

Revision ID: d19752559017
Revises: 1d26d9f6e97c
Create Date: 2026-06-12 05:21:49.851202

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = 'd19752559017'
down_revision: Union[str, Sequence[str], None] = '1d26d9f6e97c'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.alter_column('aesculap', 'eng_desc', new_column_name='eng')
    op.alter_column('aesculap', 'viet_desc', new_column_name='viet')
    op.alter_column('aesculap', 'alternative_code', new_column_name='alternative')

    op.alter_column('kls_martin', 'eng_desc', new_column_name='eng')
    op.alter_column('kls_martin', 'viet_desc', new_column_name='viet')
    op.alter_column('kls_martin', 'alternative_code', new_column_name='alternative')


def downgrade() -> None:
    """Downgrade schema."""
    op.alter_column('kls_martin', 'alternative', new_column_name='alternative_code')
    op.alter_column('kls_martin', 'viet', new_column_name='viet_desc')
    op.alter_column('kls_martin', 'eng', new_column_name='eng_desc')

    op.alter_column('aesculap', 'alternative', new_column_name='alternative_code')
    op.alter_column('aesculap', 'viet', new_column_name='viet_desc')
    op.alter_column('aesculap', 'eng', new_column_name='eng_desc')
