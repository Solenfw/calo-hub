FROM python:3.12-slim AS base

# Install uv securely
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

# Tell uv to create the virtual environment OUTSIDE of /app
# This completely prevents Docker Compose bind-mount conflicts!
ENV UV_PROJECT_ENVIRONMENT="/opt/venv"

# Ensure the virtual environment is on the PATH
ENV PATH="/opt/venv/bin:$PATH"

WORKDIR /app
COPY pyproject.toml uv.lock ./

# ==========================================
FROM base AS dev

RUN uv sync --frozen --no-cache
COPY . .
CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]

# ==========================================
FROM base AS prod

# ADDED --no-dev: You don't want dev tools (like pytest/ruff) in your production image
RUN uv sync --frozen --no-dev --no-cache
COPY . .

# Secure permissions for both the app and the virtual environment
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app /opt/venv

USER appuser
CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]