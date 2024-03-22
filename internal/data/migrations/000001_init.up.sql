
 CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    email VARCHAR(255) NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    payment_status VARCHAR(25) NOT NULL DEFAULT 'pending',
    session_id VARCHAR(255) UNIQUE NOT NULL,
    check_in VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPtz NOT NULL DEFAULT 'now()',
    updated_at TIMESTAMPtz NOT NULL DEFAULT 'now'
);

ALTER TABLE users ADD PRIMARY KEY(id);