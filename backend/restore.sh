#!/bin/bash
# Usage: ./restore.sh <filename_without_extension_or_full_path>
if [ -z "$1" ]; then
  echo "Usage: ./restore.sh <filename>"
  echo "Available backups:"
  ls -1 data/rollback/*.sql
  exit 1
fi

file="$1"
[[ "$file" != *.sql ]] && file="${file}.sql"

if [ ! -f "$file" ]; then
  echo "File not found: $file"
  exit 1
fi

docker exec -i <postgres_container> psql -U <user> <dbname> < "$file"
echo "Restored from: $file"