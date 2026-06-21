-- sql/schema.sql

CREATE TYPE order_status AS ENUM ('pending', 'completed', 'cancelled');
CREATE TYPE unit_of_measure AS ENUM ('kg', 'g', 'l', 'ml', 'unit', 'box');

CREATE TABLE raw_materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock DECIMAL(12, 4) NOT NULL DEFAULT 0.0000,
    uom unit_of_measure NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recipes (
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    raw_material_id INT NOT NULL REFERENCES raw_materials(id) ON DELETE RESTRICT,
    quantity DECIMAL(12, 4) NOT NULL,
    PRIMARY KEY (product_id, raw_material_id)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    total_amount DECIMAL(10, 2) NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL
);

CREATE TABLE financial_ledger (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
