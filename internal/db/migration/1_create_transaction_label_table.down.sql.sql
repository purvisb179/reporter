-- Drop the transaction_label table first due to foreign key constraints
DROP TABLE IF EXISTS transaction_label;

-- Drop the label table
DROP TABLE IF EXISTS label;

-- Drop the transaction table
DROP TABLE IF EXISTS transaction;

-- Finally, drop the action_type ENUM type
DROP TYPE IF EXISTS action_type;
