PRAGMA foreign_keys = ON;

CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	first_name TEXT,
	last_name TEXT,
	avatar_url TEXT, 

	is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)), 

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	
	name TEXT NOT NULL,
	description TEXT,
	color TEXT,
	
	is_archived INTEGER NOT NULL DEFAULT 0 CHECK (is_archived IN (0, 1)),

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	UNIQUE (user_id, name)
);

CREATE TABLE lists (
	id INTEGER PRIMARY KEY,
	project_id INTEGER NOT NULL,

	name TEXT NOT NULL,
	position INTEGER NOT NULL,

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (project_id)
		REFERENCES projects(id)
		ON DELETE CASCADE,

	UNIQUE (project_id, name),
	UNIQUE (project_id, position)
);

CREATE INDEX idx_projects_user_id 
ON projects(user_id);

CREATE INDEX idx_lists_project_id
ON lists(project_id);
