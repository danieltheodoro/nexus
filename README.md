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

- [x] SQLite schema
- [x] JWT Authentication
- [x] User management
- [x] Projects
- [x] Lists
- [x] Tasks
- [x] Comments
- [x] Labels
- [x] Task ↔ Label associations
- [ ] Frontend
- [ ] Docker
- [ ] Automated tests
