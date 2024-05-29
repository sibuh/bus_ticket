CREATE TABLE events(
    id SERIAL NOT NULL,
    title VARCHAR(255) NOT Null,
    description VARCHAR(500) NOT NUll,
    user_id int NOT NULL,
    start_date TIMESTAMPtz NOT NULL,
    end_date TIMESTAMPtz NOT NULL,
    created_at TIMESTAMPtz NOT NULL DEFAULT 'now()',
    updated_at TIMESTAMPtz NOT NULL DEFAULT 'now()',
    deleted_at TIMESTAMPtz NULL

);
ALTER TABLE events ADD PRIMARY KEY(id);
ALTER TABLE events ADD FOREIGN KEY(user_id) REFERENCES users(id)