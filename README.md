# drip-backend

## Setup

- install go
- install all pacakges

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
go run main.go
```

## API Docs

We use [oapi-codegen](https://github.com/deepmap/oapi-codegen) to code gen our api go types.

To install the swaggo cli run the following:

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```

To update the spec run the following from root:

```bash
 oapi-codegen ./docs/swagger.yml > pkg/swagger/generated.go
```

API docs are viewable at `http://localhost:8080/swagger/index.html`.
