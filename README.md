# MinURLy

A minimal, fast, and secure URL shortener service built with Go and MongoDB.

MinURLy provides a simple API to create, manage, and resolve short URLs.
The project is designed with clean architecture, modular components, and production-ready patterns.

### 🚀 Features

✨ Create short URLs with unique short codes
🔁 Redirect to original URLs
👤 User authentication with Google OAuth + session management
📦 MongoDB as storage
🧱 Clean, layered architecture (handlers → services → repositories)
🔒 Error handling with custom API errors
📜 Structured logging with Zerolog
⚙️ Config support using environment variables

### 🏗️ Setup & Run

1. Install dependencies

```bash
   go mod tidy
```

2. Run the server

```bash
   go run ./cmd/minurly
```

3. Build binary

```bash
   go build -o bin/minurly ./cmd/minurly
   ./bin/minurly
```

### 🧱 Tech Stack

- Go 1.22+
- MongoDB
- Zerolog for structured logging
- net/http + Gorilla

### 🤝 Contributing

PRs and suggestions are welcome!
Feel free to open issues or contribute features.

### 📄 License

MIT License. Free to use, modify, and distribute.
