CREATE TABLE products (
    id UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    flavor VARCHAR(255) NOT NULL,
    production_price DECIMAL(10, 2) NOT NULL,
    selling_price DECIMAL(10, 2) NOT NULL,
    markup_margin DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);