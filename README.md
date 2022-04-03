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

We use [swaggo](https://github.com/swaggo/swag) to code gen our api spec.

To install the swaggo cli run the following:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

To update the spec run the following from root:

```bash
swag init -g internal/server/http/http.go
```

API docs are viewable at `http://localhost:8080/swagger/index.html`.
