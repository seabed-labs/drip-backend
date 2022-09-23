# drip-backend

<a href="https://codeclimate.com/repos/62c905f4ef60b5563400002b/maintainability"><img src="https://api.codeclimate.com/v1/badges/252814ca6aba27f4dc3d/maintainability" /></a>
<a href="https://codeclimate.com/repos/62c905f4ef60b5563400002b/test_coverage"><img src="https://api.codeclimate.com/v1/badges/252814ca6aba27f4dc3d/test_coverage" /></a>
[![CI](https://github.com/dcaf-labs/drip-backend/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/dcaf-labs/drip-backend/actions/workflows/ci.yaml)

## Mainnet

[![Better Uptime Badge](https://betteruptime.com/status-badges/v1/monitor/goyh.svg)](https://betteruptime.com/?utm_source=status_badge)

## Devnet

[![Better Uptime Badge](https://betteruptime.com/status-badges/v1/monitor/g7cf.svg)](https://betteruptime.com/?utm_source=status_badge)

## Deploy Process

### Devnet Staging

- merge into `main`

### Devnet Prod

- merge into `devnet`

### Mainnet Prod

- upgrade devnet prod via pipeline

## Setup

- install go 1.18
- install all packages

```bash
# Add to ZSHRC
export GOPRIVATE=github.com/dcaf-labs/solana-go-clients
go get -u ./...
```

Add the following to ~/.gitconfig

```txt
[url "ssh://git@github.com/"]
insteadOf = https://github.com/
```

- setup a `.env` file
- ex:

```env

# Devnet Staging
ENV="STAGING"
NETWORK="DEVNET"
DRIP_PROGRAM_ID="F1NyoZsUhJzcpGyoEqpDNbUMKVvCnSXcCki1nN3ycAeo"
# random wallet, this is the mint auth
DRIP_BACKEND_WALLET="[141,241,173,131,255,186,170,216,65,246,24,196,173,94,39,225,161,108,251,102,177,20,166,223,13,69,103,38,242,107,72,194,177,170,44,204,179,183,235,4,231,51,88,169,156,153,132,247,235,166,41,123,87,219,139,204,95,1,176,98,72,90,51,82]"

# Test Channel
DISCORD_WEBHOOK_ID=1021592812954857492
DISCORD_ACCESS_TOKEN=qPeOyI4e4k6kYah44k9_PXFQDsuLO7lbHcazLrsKcvzqvrQh_lr1PK21kB3GZCSTv2Xg

PORT="8080"
PSQL_USER="dcaf"
PSQL_PASS="drip"
PSQL_DBNAME="drip"
PSQL_PORT="5432"
PSQL_HOST="localhost"
PSQL_SSLMODE=disable
OUTPUT=./internal/pkg/repository/models
GOOGLE_KEY="540992596258-sa2h4lmtelo44tonpu9htsauk5uabdon.apps.googleusercontent.com"
GOOGLE_SECRET="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"
```

- run all tests

```bash
go test ./...
```

- run the server

```bash
ENV=DEVNET go run main.go
```

## Tests

If an interface is changed, the associated mocks need to be re-generated.

```bash
# ex: from inside drip-backend/pkg/clients/solana
./scripts/create-mocks.sh
```

## API Docs

We use [oapi-codegen](https://github.com/deepmap/oapi-codegen) to code gen our api go types.

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```

To update the spec run the following from root:

```bash
 oapi-codegen -package apispec ./docs/swagger.yaml > pkg/api/apispec/generated.go
```

API docs are viewable at `http://localhost:8080/swagger.json`.

## Database

Locally the db can be started with

```bash
docker-compose --file ./build/docker-compose.yaml  --env-file ./.env up
```

and stopped with

```bash
docker-compose --file ./build/docker-compose.yaml  --env-file ./.env down
```

> **_NOTE:_** TODO(mocha) do codegen in a dockerized db that is setup and destroyed automatically in the script.

### Migrations

- located in `internal/database/psql/migrations`
- All migrations that are a number larger then what the db version is will be run automatically on startup or during codegen
- Running the migrations will automatically increment the schema version in the db

```bash
go run cmd/migrate/main.go
```

> **_NOTE:_** The DB must be running prior to running this script.

### Codegen

- Database [models](app/internal/data/psql/generated) are generated using the database schema via [go-gorm/gen](https://github.com/go-gorm/gen)
- Before generating the models, the database needs to be running
- The codegen script will also run migrations if needed

```bash
go run cmd/codegen/main.go
```

> **_NOTE:_** The DB must be running prior to running this script.

### Process for Creating/Updating Database Models

- Create a migration file under `internal/database/psql/migrations`, and name it appropriately (ex. `2_new_migration.up.sql`)
- Run the migration + codegen script `go run cmd/codegen/main.go`

> **_NOTE:_** The DB must be running prior to running this script.

### Backfill DB for Devnet

Run the following script

```bash
go run cmd/backfill/main.go
```

> **_NOTE:_** This backfills based on the content of `devnet.yaml`.
