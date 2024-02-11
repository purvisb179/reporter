CREATE TYPE action_type AS ENUM ('add', 'subtract', 'adjust');

CREATE TABLE transaction (
                             id SERIAL PRIMARY KEY,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             amount DECIMAL(10, 2) NOT NULL,
                             action ENUM('add', 'subtract', 'adjust') NOT NULL
);

CREATE TABLE label (
                       id SERIAL PRIMARY KEY,
                       key VARCHAR(255) NOT NULL,
                       value VARCHAR(255) NOT NULL,
                       UNIQUE(key, value) -- Assuming a combination of key and value should be unique.
);

CREATE TABLE transaction_label (
                                   transaction_id INT NOT NULL,
                                   label_id INT NOT NULL,
                                   PRIMARY KEY (transaction_id, label_id),
                                   FOREIGN KEY (transaction_id) REFERENCES transaction(id) ON DELETE CASCADE,
                                   FOREIGN KEY (label_id) REFERENCES label(id) ON DELETE CASCADE
);
