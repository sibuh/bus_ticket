CREATE TABLE sessions (
    id UUID PRIMARY KEY, 
	ticket_id UUID NOT NULL,     
	payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
	payment_url    VARCHAR(200) NOT NULL,
	cancel_url     VARCHAR(200),
	amount        FLOAT NOT NULL,
	created_at     TIMESTAMPtz NOT NULL
);
ALTER TABLE sessions ADD CONSTRAINT sessions_ticket_id_foreign FOREIGN KEY(ticket_id) REFERENCES tickets(id);