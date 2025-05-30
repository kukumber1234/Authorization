CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) DEFAULT 'author'
);

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP,
    rejected_at TIMESTAMP
);


INSERT INTO users (email, password_hash, username, role) VALUES
('admin@gmail.com',   '1234', 'AdminUser', 'admin'),
('moderator@gmail.com', '1234', 'ModeratorUser', 'moderator'),
('user@gmail.com',    '1234', 'SimpleUser', 'author'),
('john.doe@mail.com', 'johnpass', 'JohnD', 'author'),
('jane.smith@mail.com', 'janepass', 'JaneS', 'moderator'),
('alex.ivanov@mail.com', 'alexpass', 'AlexI', 'author'),
('kate.petrova@mail.com', 'katepass', 'KateP', 'moderator'),
('mike.brown@mail.com', 'mikepass', 'MikeB', 'author'),
('emma.white@mail.com', 'emmapass', 'EmmaW', 'author'),
('dave.admin@mail.com', 'adminpass', 'DaveA', 'admin');