#!/bin/bash

set -e

echo "🧱 Running database migration..."
go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrate postgres "$DATABASE_URL" up

echo "🛠️  Generating SQLC models..."
sqlc generate

echo "✅ Done!"