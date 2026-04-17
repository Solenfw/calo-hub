.PHONY: dev prod down logs

# Start dev environment with hot reload
# Automatically merges docker-compose.yml AND docker-compose.override.yml
dev:
	docker compose up --build

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