CREATE TABLE datacars (
	vin text PRIMARY KEY,
	brand text NOT NULL,
	model text NOT NULL,
	price int NOT NULL,
	carstatus text NOT NULL,
	odometer int NOT NULL
);