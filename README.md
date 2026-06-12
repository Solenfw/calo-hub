# Calo Hub

Internal medical production suite for managing surgical instrument catalogs, documentation, and clinical workflows.

> ⚠️ Active development — things will break.

## Stack

- **Frontend** — Next.js 16, Tailwind CSS v4, Motion
- **Backend** — FastAPI, SQLAlchemy, PostgreSQL
- **Infra** — Docker Compose

## Getting Started

```bash
# Start the database
docker compose up -d

# Backend
cd backend && uv sync
uv run uvicorn app.main:app --reload

# Frontend
cd frontend && bun install
bun dev
```

Requires a `.env` file at the root. See `.env.example` if available.

## Structure

```
calo-hub/
├── frontend/   # Next.js app
├── backend/    # FastAPI + PostgreSQL
└── docker-compose.yml
```

## Features (WIP)

- Online catalog search (KLS Martin / B. Braun Aesculap)
- Quick Convert — cross-brand product code mapping
- OCR — extract data from clinical documents
- eSelector — filter components by surgical requirements