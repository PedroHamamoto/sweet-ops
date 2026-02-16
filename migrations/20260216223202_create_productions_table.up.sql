CREATE  TABLE productions(
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (product_id) REFERENCES products(id)
);