CREATE TABLE events(
    id SERIAL NOT NULL,
    title VARCHAR(255) NOT Null,
    description VARCHAR(500) NOT NUll,
    user_id int NOT NULL 
);
ALTER TABLE events ADD PRIMARY KEY(id);
ALTER TABLE events ADD FOREIGN KEY(user_id) REFERENCES users(id)