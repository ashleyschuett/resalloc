CREATE TABLE IF NOT EXISTS users (
  username VARCHAR(40) PRIMARY KEY,
  salt VARCHAR(40),
  password VARCHAR(100),
  account_type VARCHAR(3),
  activated BOOLEAN
);
