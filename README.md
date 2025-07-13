# Invoice Microservice

This is a simple **invoice microservice** built in **Go**, using:

- ✅ **Hexagonal architecture** (also known as Ports & Adapters)
- 📨 **NATS** as the communication protocol (no HTTP)
- 🛢️ **PostgreSQL** as the database
- 💾 **GORM** as the ORM
- 🐳 Built with a **multi-stage Dockerfile**, outputting a minimal binary-only container
- 🧪 Easily testable via the NATS CLI

---

## 📂 Project Structure

```bash
.
├── cmd/
│   └── main.go               # Application entry point
├── internal/
│   ├── adapter/
│   │   ├── db/               # DB adapter (GORM)
│   │   │   └── db.go
│   │   └── nats/             # NATS adapter
│   │       └── nats.go
│   ├── application/          # Business logic
│   │   └── invoice.go
│   ├── config/               # Configuration
│   │   └── config.go
│   ├── domain/               # Domain model
│   │   └── invoice.go
├── go.mod
├── Dockerfile
└── README.md
```

---

## 🚀 Features

- Create, read, update and delete invoices via NATS
- UUID-based primary keys
- Auto-migration with GORM
- Error responses in JSON format with status codes
- Environment-based configuration (no hardcoded URLs)

---

## 🛠️ Configuration

Environment variables required:

| Variable       | Description                                                                  |
| -------------- | ---------------------------------------------------------------------------- |
| `NATS_URL`.    | NATS server URL (e.g. `nats://localhost:4222`)                               |
| `DATABASE_URL` | PostgreSQL URL (e.g. `postgresql://postgres:secret@localhost:5432/invoices`) |

You can set them in your shell or via Docker `-e` flags.

---

## 🐳 Running with Docker

Build the binary-only image (multi-stage):

```bash
docker build -t invoice-ms .
```

Run the container:

```bash
docker run -d   --name invoice-service   --network host   -e NATS_URL=nats://localhost:4222   -e DATABASE_URL="postgresql://postgres:secret@localhost:5432/invoices"   invoice-ms
```

---

## 🧪 Testing with NATS CLI

### Create Invoice

```bash
nats request invoice.create '{"amount": 199.99, "customer": "John Doe"}'
```

### Get Invoice by ID

```bash
nats request invoice.get '{"id": "your-uuid-here"}'
```

### Get All Invoices

```bash
nats request invoice.get_all ''
```

### Update Invoice

```bash
nats request invoice.update '{"id": "your-uuid-here", "amount": 299.99, "customer": "Jane Doe"}'
```

### Delete Invoice

```bash
nats request invoice.delete '{"id": "your-uuid-here"}'
```

---

## 🧰 Development

### Run locally

```bash
go run ./cmd/main.go
```

Make sure you have NATS and PostgreSQL running and the environment variables set.

---

## ✅ Requirements

- Go 1.20+
- PostgreSQL 13+
- NATS Server
- Docker (optional for containerized deployment)

---

## ✨ License

MIT License.
