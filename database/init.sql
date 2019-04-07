CREATE DATABASE auroradb;
\c auroradb;
CREATE EXTENSION pgcrypto;
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.modified = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TABLE userprofile (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    first VARCHAR(255) NOT NULL,
    last VARCHAR(255) NOT NULL,
    alias VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(64),
    birthday TIMESTAMPTZ,
    state VARCHAR(64),
    notes TEXT,
    picture TEXT
);
CREATE TABLE userlogin (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    userprofile_id INTEGER NOT NULL,
    FOREIGN KEY (userprofile_id) REFERENCES userprofile (id),
    username VARCHAR(255) NOT NULL,
    password TEXT NOT NULL
);
CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES userprofile (id),
    restricted BOOLEAN NOT NULL,
    status VARCHAR(255),
    title TEXT,
    body TEXT
);
CREATE TABLE connected (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES userprofile (id),
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task (id),
    userprofile_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES userprofile (id)
);
CREATE TRIGGER task_timestamp
BEFORE UPDATE ON task
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TABLE usertask (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task (id),
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES userprofile (id)
);
CREATE TABLE timesheets (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    amount NUMERIC,
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task (id)
);
CREATE TABLE comment (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES userprofile (id),
    task_id INTEGER NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task (id),
    body TEXT
);
CREATE TRIGGER comment_timestamp
BEFORE UPDATE ON comment
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
CREATE TABLE tag (
    id SERIAL PRIMARY KEY,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_DATE,
    owner_id INTEGER NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES userprofile (id),
    comment_id INTEGER NOT NULL,
    FOREIGN KEY (comment_id) REFERENCES comment (id),
    body TEXT
);