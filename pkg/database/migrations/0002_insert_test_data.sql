INSERT INTO Category (name, product_id) VALUES
    ('Electronic', NULL),
    ('Books', NULL),
    ('Clothing', NULL);

INSERT INTO Products (name, category_id) VALUES
    ('Laptop', 1),
    ('Smartphone', 1),
    ('Novel', 2),
    ('Hoodie', 3);

UPDATE Category SET product_id = (SELECT id FROM Products WHERE name = 'Laptop') WHERE name = 'Electronics';
UPDATE Category SET product_id = (SELECT id FROM Products WHERE name = 'Novel') WHERE name = 'Books';
UPDATE Category SET product_id = (SELECT id FROM Products WHERE name = 'Hoodie') WHERE name = 'Clothing';
