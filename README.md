# Nexus

Task manager in development.

## Status

In development

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

## Roadmap

- [x] Database schema
- [ ] Seed data
- [ ] Backend
- [ ] REST API
- [ ] Authentication
- [ ] Frontend
