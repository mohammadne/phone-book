CREATE TABLE IF NOT EXISTS contacts(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	phones VARCHAR(15)[],
	description VARCHAR(255),
	user_id INTEGER REFERENCES users (id)
);

CREATE INDEX user_id_idx ON contacts (user_id);
