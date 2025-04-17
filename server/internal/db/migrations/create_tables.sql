CREATE TYPE user_status AS ENUM ('online', 'offline', 'away', 'banned');

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP,
    status user_status NOT NULL DEFAULT 'offline'
);

CREATE TYPE chat_room_type AS ENUM ('direct', 'group');

CREATE TABLE IF NOT EXISTS chat_rooms (
    id bigserial PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type chat_room_type NOT NULL DEFAULT 'direct',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    creator_id bigserial REFERENCES users(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id bigserial PRIMARY KEY,
    sender_id bigserial REFERENCES users(id) NOT NULL,
    chat_room_id bigserial REFERENCES chat_rooms(id) NOT NULL,
    encrypted_content TEXT NOT NULL,
    reply_to_message_id bigserial REFERENCES messages(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TYPE chat_room_role AS ENUM ('admin', 'moderator', 'member');

CREATE TABLE IF NOT EXISTS chat_room_members (
    user_id bigserial REFERENCES users(id) NOT NULL,
    chat_room_id bigserial REFERENCES chat_rooms(id) NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    role chat_room_role NOT NULL DEFAULT 'member',
    PRIMARY KEY (user_id, chat_room_id)
);
