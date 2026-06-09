.PHONY: dev prod down logs test test-backend test-frontend

compose:
	docker compose up -d

# Start dev environment with hot reload
dev:
	uv run uvicorn backend.app.main:app --reload

# Start prod environment in detached mode
# Explicitly loads the base file, then applies the prod overrides
prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d

# Stop and remove containers
# Note: If tearing down prod, you technically need the same -f flags used to bring it up.
down:
	docker compose down
	@echo "If you need to tear down Prod, run: docker compose -f docker-compose.yml -f docker-compose.prod.yml down"

# Tail api logs
logs:
	docker compose logs -f api

# Run all tests
test: test-backend test-frontend

# Run backend unit tests
test-backend:
	UV_CACHE_DIR=.uv-cache uv run --project backend python -m unittest discover -s tests/backend -p 'test_*.py'

# Run frontend unit tests
test-frontend:
	node --test tests/frontend/*.test.mjs
