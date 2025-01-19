CREATE TABLE IF NOT EXISTS product_stock (
    product_id INT PRIMARY KEY,
    stock INT NOT NULL CHECK (stock >= 0),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

