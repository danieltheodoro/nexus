PRAGMA foreign_keys = ON;

CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	
	username TEXT NOT NULL UNIQUE
		CHECK(length(trim(username)) > 0),
	email TEXT NOT NULL UNIQUE
		CHECK(email LIKE '%@%.%'),
	password_hash TEXT NOT NULL 
		CHECK(length(trim(password_hash)) > 0),
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
	
	name TEXT NOT NULL
		CHECK(length(trim(name)) > 0),
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

	name TEXT NOT NULL
		CHECK(length(trim(name)) > 0),
	
	position INTEGER NOT NULL,

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (project_id)
		REFERENCES projects(id)
		ON DELETE CASCADE,

	UNIQUE (project_id, name),
	UNIQUE (project_id, position)
);

CREATE TABLE tasks (
	id INTEGER PRIMARY KEY,
	list_id INTEGER NOT NULL,
	creator_id INTEGER NOT NULL,

	title TEXT NOT NULL 
		CHECK(length(trim(title)) > 0),
	description TEXT,

	priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),

	position INTEGER NOT NULL,

	due_date TEXT,
	completed_at TEXT,

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (list_id)
		REFERENCES lists(id)
		ON DELETE CASCADE,

	FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	UNIQUE (list_id, position)
);

CREATE TABLE labels (
	id INTEGER PRIMARY KEY,
	project_id INTEGER NOT NULL,

	name TEXT NOT NULL
		CHECK(length(trim(name)) > 0),
	color TEXT NOT NULL
		CHECK(length(trim(color)) > 0),

	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (project_id)
		REFERENCES projects(id)
		ON DELETE CASCADE,

	UNIQUE (project_id, name)
);

CREATE TABLE task_labels (
	task_id INTEGER NOT NULL,
	label_id INTEGER NOT NULL,

	PRIMARY KEY (task_id, label_id),

	FOREIGN KEY (task_id)
		REFERENCES tasks(id)
		ON DELETE CASCADE,

	FOREIGN KEY (label_id)
		REFERENCES labels(id)
		ON DELETE CASCADE
);

CREATE INDEX idx_projects_user_id 
ON projects(user_id);

CREATE INDEX idx_lists_project_id
ON lists(project_id);

CREATE INDEX idx_tasks_list_id
ON tasks(list_id);
CREATE INDEX idx_tasks_creator_id 
ON tasks(creator_id);

CREATE INDEX idx_labels_project_id
ON labels(project_id);

CREATE INDEX idx_task_labels_label_id
ON task_labels(label_id);
