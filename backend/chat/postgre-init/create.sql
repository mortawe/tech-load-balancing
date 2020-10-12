CREATE TABLE IF NOT EXISTS users
(
    id         serial                                 NOT NULL
        CONSTRAINT user_pk
            PRIMARY KEY,
    username   varchar                                NOT NULL,
    created_at timestamp WITH TIME ZONE DEFAULT now() NOT NULL
);

ALTER TABLE users
    OWNER TO "user";

CREATE UNIQUE INDEX IF NOT EXISTS user_id_uindex
    ON users (id);

CREATE UNIQUE INDEX IF NOT EXISTS user_username_uindex
    ON users (username);

CREATE TABLE IF NOT EXISTS chats
(
    id         serial                                 NOT NULL
        CONSTRAINT chat_pk
            PRIMARY KEY,
    name       varchar,
    created_at timestamp WITH TIME ZONE DEFAULT now() NOT NULL
);

ALTER TABLE chats
    OWNER TO "user";

CREATE UNIQUE INDEX IF NOT EXISTS chat_id_uindex
    ON chats (id);

CREATE UNIQUE INDEX IF NOT EXISTS chat_name_uindex
    ON chats (name);

CREATE TABLE IF NOT EXISTS chat_users
(
    chat_id integer NOT NULL
        CONSTRAINT chat_users_chats_id_fk
            REFERENCES chats
            ON UPDATE CASCADE ON DELETE CASCADE,
    user_id integer NOT NULL
        CONSTRAINT chat_users_users_id_fk
            REFERENCES users
            ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE UNIQUE INDEX chat_users_chat_id_user_id_uindex
    ON chat_users (chat_id, user_id);


ALTER TABLE chat_users
    OWNER TO "user";

CREATE TABLE IF NOT EXISTS messages
(
    id         serial                                 NOT NULL
        CONSTRAINT messages_pk
            PRIMARY KEY,
    chat_id    integer                                NOT NULL
        CONSTRAINT messages_chat_id_fk
            REFERENCES chats
            ON UPDATE CASCADE ON DELETE CASCADE,
    author_id  integer                                NOT NULL
        CONSTRAINT messages_user_id_fk
            REFERENCES users
            ON UPDATE CASCADE ON DELETE CASCADE,
    text       varchar,
    created_at timestamp WITH TIME ZONE DEFAULT now() NOT NULL,
    CONSTRAINT messages_chat_users_user_id_chat_id_fk
        FOREIGN KEY (author_id, chat_id) REFERENCES chat_users (user_id, chat_id)
);

ALTER TABLE messages
    OWNER TO "user";

CREATE UNIQUE INDEX IF NOT EXISTS messages_id_uindex
    ON messages (id);

CREATE INDEX chat_users_user_id_index
    ON chat_users (user_id);

CREATE INDEX messages_chat_id_index
    ON messages (chat_id);
create index messages_chat_id_created_at_index
	on messages (chat_id, created_at);

create index messages_chat_id_created_at_index_2
	on messages (chat_id asc, created_at desc);
