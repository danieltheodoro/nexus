# Nexus

> Management made easy.

Nexus is a modern management platform built to simplify everyday operations for teams and organizations.

## Status

Active development

## Tech Stack

### Frontend

- React Native
- Expo
- TypeScript

### Backend

- Go
- SQLite
- JWT Authentication

## Current Features

- [x] SQLite schema
- [x] JWT Authentication
- [x] User management
- [x] Projects
- [x] Lists
- [x] Tasks
- [x] Comments
- [x] Labels
- [x] Task ↔ Label associations
- [x] Authentication UI

## Roadmap

- [ ] Dashboard
- [ ] Employee Management
- [ ] Customer Management
- [ ] Inventory Management
- [ ] Reports
- [ ] Settings

## Database Model

```text
users
│
├── projects
│   ├── lists
│   │   └── tasks
│   │        ├── comments
│   │        └── task_labels
│   │             ▲
│   └── labels ───┘
│
└──────────────► comments
```

## Getting Started

### Clone the repository

```bash
git clone https://github.com/danieltheodoro/nexus.git
```

### Install dependencies

```bash
cd nexus/frontend
npm install
```

### Start the development server

```bash
npx expo start
```

## License

This project is licensed under the MIT License.
