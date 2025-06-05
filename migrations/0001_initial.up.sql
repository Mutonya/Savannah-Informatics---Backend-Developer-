-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum type for order status
CREATE TYPE order_status AS ENUM ('pending', 'completed', 'cancelled');

-- Create customers table
CREATE TABLE customers (
                           id SERIAL PRIMARY KEY,
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           deleted_at TIMESTAMP WITH TIME ZONE,
                           first_name VARCHAR(100) NOT NULL,
                           last_name VARCHAR(100) NOT NULL,
                           email VARCHAR(255) NOT NULL UNIQUE,
                           phone VARCHAR(20) NOT NULL,
                           address VARCHAR(255),
                           oauth_id VARCHAR(255) UNIQUE
);

-- Create categories table with self-referential relationship
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMP WITH TIME ZONE,
                            name VARCHAR(100) NOT NULL UNIQUE,
                            parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

-- Create products table
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          deleted_at TIMESTAMP WITH TIME ZONE,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price DECIMAL(10,2) NOT NULL,
                          sku VARCHAR(100) UNIQUE,
                          category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT
);

-- Create orders table
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP WITH TIME ZONE,
                        customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
                        status order_status DEFAULT 'pending',
                        total DECIMAL(10,2) NOT NULL
);

-- Create order_items table
CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY,
                             created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             deleted_at TIMESTAMP WITH TIME ZONE,
                             order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
                             quantity INTEGER NOT NULL,
                             price DECIMAL(10,2) NOT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);