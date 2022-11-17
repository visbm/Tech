CREATE TABLE IF NOT EXISTS accounts (
  account_id serial PRIMARY KEY,
  balance real NOT NULL
  CONSTRAINT balance_cannot_be_negative CHECK (balance >= 0)  
); 
CREATE TABLE IF NOT EXISTS transactions (
  transaction_id serial PRIMARY KEY,
  account_id  integer REFERENCES accounts (account_id) NOT NULL,
  amount real NOT NULL,
  date_time timestamp NOT NULL,
  comment text NOT NULL
); 
