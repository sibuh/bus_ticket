CREATE TABLE tickets (
    id uuid primary key default gen_random_uuid(),
    trip_id  UUID NOT NULL,
    bus_id UUID NOT NULL,
    user_id uuid NOT NULL,
    status VARCHAR(20)
);
ALTER TABLE tickets ADD CONSTRAINT tickets_trip_id_foreign FOREIGN KEY(trip_id) REFERENCES trips(id);
ALTER TABLE tickets ADD CONSTRAINT tickets_bus_id_foreign FOREIGN KEY(bus_id) REFERENCES buses(id);
ALTER TABLE tickets ADD CONSTRAINT tickets_user_id_foreign FOREIGN KEY(user_id) REFERENCES users(id);