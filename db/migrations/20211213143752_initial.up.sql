CREATE TABLE credentials (
    account_id TEXT primary key,
    value      BLOB
);

--bun:split

CREATE TABLE access_tokens (
    account_id   TEXT primary key,
    access_token TEXT,
    id_token TEXT
);
