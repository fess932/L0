DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  order_uid VARCHAR(255) NOT NULL,
  track_number VARCHAR(255) NOT NULL,
  entry VARCHAR(255) NOT NULL,
  date_created TIMESTAMP
);