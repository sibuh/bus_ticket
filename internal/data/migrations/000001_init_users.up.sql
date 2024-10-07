
 CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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