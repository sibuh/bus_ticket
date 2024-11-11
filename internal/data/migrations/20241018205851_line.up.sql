CREATE TABLE lines (
    id uuid primary key default gen_random_uuid(),
    destination VARCHAR(100) NOT NULL,
    departure VARCHAR(100) NOT NULL,
    price float NOT NULL,
    schedule JSONB NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    update_at timestamptz
);
CREATE TABLE line_trips (
    id uuid primary key default gen_random_uuid(),
    line uuid NOT NULL,
    date timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT 'now()',
    update_at timestamptz
);
ALTER TABLE line_trips
ADD CONSTRAINT lines_id_foreign FOREIGN KEY(line) REFERENCES lines(id);