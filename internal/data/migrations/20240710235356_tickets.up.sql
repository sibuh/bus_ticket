CREATE TABLE tickets (
    id SERIAL NOT NULL,
    trip_id int NOT NULL,
    bus_no int NOT NULL,
    ticket_no int NOT NULL,
    status VARCHAR(20),
    user_id int
);
ALTER TABLE tickets
ADD FOREIGN KEY(user_id) REFERENCES users(id)