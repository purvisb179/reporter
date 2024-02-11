-- Create the transaction table
CREATE TABLE transaction (
                             id UUID PRIMARY KEY,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             amount INTEGER NOT NULL
);

-- Create the label table
CREATE TABLE label (
                       id SERIAL PRIMARY KEY,
                       key VARCHAR(255) NOT NULL,
                       value VARCHAR(255) NOT NULL,
                       UNIQUE(key, value)
);

-- Create the transaction_label table
CREATE TABLE transaction_label (
                                   transaction_id UUID NOT NULL,
                                   label_id INT NOT NULL,
                                   PRIMARY KEY (transaction_id, label_id),
                                   FOREIGN KEY (transaction_id) REFERENCES transaction(id) ON DELETE CASCADE,
                                   FOREIGN KEY (label_id) REFERENCES label(id) ON DELETE CASCADE
);
