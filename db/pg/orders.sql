DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
	order_data jsonb NOT NULL
--
-- 	sm_id INTEGER NOT NULL,
-- 	order_uid VARCHAR(255) NOT NULL,
-- 	track_number VARCHAR(255) NOT NULL,
-- 	entry VARCHAR(255) NOT NULL,
-- 	"locale" VARCHAR(255) NOT NULL,
-- 	internal_signature VARCHAR(255) NOT NULL,
-- 	customer_id VARCHAR(255) NOT NULL,
-- 	delivery_service VARCHAR(255) NOT NULL,
-- 	shardkey VARCHAR(255) NOT NULL,
-- 	oof_shard VARCHAR(255) NOT NULL,
--
--
--
-- 	date_created TIMESTAMP
);