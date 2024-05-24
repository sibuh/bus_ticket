
 CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPtz NOT NULL DEFAULT 'now()',
    updated_at TIMESTAMPtz NULL,
    deleted_at TIMESTAMPtz NULL 
);

ALTER TABLE users ADD PRIMARY KEY(id);