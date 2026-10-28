-- Order Management System — Database Initialization
-- Run: psql -U postgres -f scripts/init.sql

CREATE DATABASE orderdb;

\c orderdb;

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    product_id VARCHAR(50) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED'
        CHECK (status IN ('CREATED', 'CONFIRMED', 'CANCELLED')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- Seed data for testing
INSERT INTO orders (customer_id, product_id, quantity, price, total_amount, status) VALUES
(101, 'P001', 2, 29.99, 59.98, 'CREATED'),
(102, 'P002', 1, 49.99, 49.99, 'CONFIRMED'),
(101, 'P003', 3, 19.99, 59.97, 'CREATED');
