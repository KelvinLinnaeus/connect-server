# Univyn Connect Server

A comprehensive social networking and learning environment backend built with Go, Gin, and PostgreSQL.

##  Features

- ✅ **User Management**: Registration, authentication, profile management
- ✅ **Authentication**: Secure PASETO token-based auth with refresh tokens
- ✅ **Sessions**: Session management, multi-device support, security tracking
- ✅ **Posts & Social**: Complete feed system, comments, likes, reposts, search
- ✅ **Communities**: Create and manage communities with moderation, search, and discovery
- ✅ **Groups**: Project/study/social groups with roles, applications, and member management
- ✅ **Messaging**: Direct messages, group conversations, channels, reactions, read receipts
- ✅ **Notifications**: Multi-type notifications with priorities and action flags
- ✅ **Events**: Campus events, registrations, attendance tracking, co-organizers
- ✅ **Announcements**: Space-wide announcements with targeting and scheduling
- ✅ **Rate Limiting**: Configurable per-endpoint rate limiting
- ✅ **Middleware**: Authentication, RBAC, CORS, logging, recovery
- 🚧 **Mentorship & Tutoring**: Connect students with mentors and tutors (Coming in Phase 6+)
- 🚧 **Analytics**: System metrics and engagement tracking (Coming in Phase 6+)

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 15+
- golang-migrate CLI
- sqlc (for code generation)

## 🛠️ Installation

### 1. Clone the repository

```bash
git clone https://github.com/connect-univyn/connect_server.git
cd connect_server
```

## 🏃 Running the Server

### Development Mode

```bash
make server
```

The server will start on `http://localhost:8080` (or your configured `SERVER_ADDRESS`).

### Production Build

```bash
make build
./bin/connect
```

## 🧪 Testing

### Run all tests

```bash
make test
```

### Run integration tests

```bash
make test-integration

```
### Access Help

```bash
make help
```

