# Nexus

Task manager in development.

## Status

Database and backend in development

## Database Model

```text
users
│
├── projects
│   ├── lists
│   │   └── tasks
│   │        ├── comments
│   │        │
│   │        └── task_labels
│   │             ▲
│   └── labels ───┘
│
└──────────────► comments
```

## Tech Stack

- SQLite
- Go

## Roadmap

- [x] Database schema
- [x] Seed data
- [x] HTTP server
- [x] Database initialization
- [ ] REST API
- [ ] Authentication
- [ ] Frontend
