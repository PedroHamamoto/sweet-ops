CREATE TABLE sales (
    id UUID PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    self_consumption BOOLEAN NOT NULL DEFAULT FALSE,
    total DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
