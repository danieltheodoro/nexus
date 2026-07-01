INSERT INTO users (
    username,
    email,
    password_hash,
    first_name,
    last_name
)
VALUES (
    'daniel',
    'daniel@example.com',
    'CHANGE_ME',
    'Daniel',
    'Theodoro'
);

INSERT INTO projects (
    user_id,
    name,
    description,
    color
)
VALUES (
    1,
    'Nexus',
    'Development',
    '#3B82F6'
);

INSERT INTO lists (project_id, name, position)
VALUES
(1, 'Backlog', 1),
(1, 'To Do', 2),
(1, 'In Progress', 3),
(1, 'Review', 4),
(1, 'Done', 5);

INSERT INTO labels (project_id, name, color)
VALUES
(1, 'Bug', '#EF4444'),
(1, 'Feature', '#22C55E'),
(1, 'Documentation', '#3B82F6'),
(1, 'High Priority', '#F59E0B');

INSERT INTO tasks (
    list_id,
    creator_id,
    title,
    priority,
    position
)
VALUES
(2, 1, 'Implement authentication', 'high', 1),
(2, 1, 'Write README', 'medium', 2);

INSERT INTO task_labels
VALUES
(1, 2),
(1, 4),
(2, 3);

INSERT INTO comments (
    task_id,
    author_id,
    content
)
VALUES (
    1,
    1,
    'Start with authentication.'
);
