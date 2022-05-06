# drip-backend

<a href="https://codeclimate.com/repos/626475a0ceebd6791a0001bb/maintainability"><img src="https://api.codeclimate.com/v1/badges/b7d1181add6eb7f241ed/maintainability" /></a>
<a href="https://codeclimate.com/repos/626475a0ceebd6791a0001bb/test_coverage"><img src="https://api.codeclimate.com/v1/badges/b7d1181add6eb7f241ed/test_coverage" /></a>
<a href="https://github.com/dcaf-protocol/drip-backend/actions/workflows/deploy.yaml"><img src="https://github.com/dcaf-protocol/drip-backend/actions/workflows/deploy.yaml/badge.svg?branch=main" /></a>
<a href="https://github.com/dcaf-protocol/drip-backend/actions/workflows/build-test.yaml"><img src="https://github.com/dcaf-protocol/drip-backend/actions/workflows/build-test.yaml/badge.svg?branch=main" /></a>
## Setup

- install go
- install all packages

```bash
go get -u ./...
```

- setup `.env` file

ex:

```env
ENV=DEVNET
# must be the mint authority (tokenOwnerKeypair from setupKepperBot.ts)
DRIP_BACKEND_WALLET="[some byte array]"
```

- run the server

```bash
ENV=DEVNET go run main.go
```

## API Docs

We use [oapi-codegen](https://github.com/deepmap/oapi-codegen) to code gen our api go types.

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```

To update the spec run the following from root:

```bash
 oapi-codegen ./docs/swagger.yml > pkg/swagger/generated.go
```

API docs are viewable at `http://localhost:8080/swagger.json`.

## Data: Postgres Database Migrations and Models
- Database [models](app/internal/data/psql/generated) are generated using the database schema via [sqlboiler](https://github.com/volatiletech/sqlboiler)
- sqlboiler introspects the database schema and creates the model files
- Before generating the models, the database needs to be running, and the migrations need to be executed
```bash
docker-compose --file ./build/docker-compose.yaml  --env-file ./.env up -d
./scripts/run-database-migrations.sh
./scripts/generate-database-models.sh
```
- Note: Running `main.go` will automatically run relevant migrations

### Process for Creating New Database Models
- Create a migration file under `data/psql/migrations`, and name it appropriately (ex. `2_new_migration.up.sql`)
- Run the migration script `./scripts/run-database-migrations.sh # run from the root of the repo`
- Generate database models `./scripts/generate-database-models.sh # run from the root of the repo`
