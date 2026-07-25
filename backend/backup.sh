#!/bin/bash
# Usage: ./backup.sh [optional_name]
name=${1:-backup_$(date +%Y%m%d_%H%M%S)}
docker exec -t <postgres_container> pg_dump -U <user> <dbname> > "data/rollback/${name}.sql"
echo "Saved: data/rollback/${name}.sql"