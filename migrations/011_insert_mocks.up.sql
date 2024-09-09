-- Insert 50 Mock Users
DO
$$
    BEGIN
        FOR i IN 1..50 LOOP
                INSERT INTO users (username, password, email, is_blocked, login_attempts)
                VALUES (
                           CONCAT('user', i), -- unique username
                           'password123',
                           CONCAT('user', i, '@example.com'),
                           false,
                           0
                       );
            END LOOP;
    END;
$$;

-- Insert Categories
INSERT INTO categories (name, description) VALUES
                                               ('Electronics', 'Electronic items like phones, laptops, etc.'),
                                               ('Clothing', 'Men and Women clothing items'),
                                               ('Books', 'Various types of books'),
                                               ('Home Appliances', 'Home and kitchen appliances');

-- Insert 100 Mock Products
DO
$$
    BEGIN
        FOR i IN 1..100 LOOP
                INSERT INTO products (name, description, price, category_id, user_id)
                VALUES (
                           CONCAT('Product ', i),
                           'Description for product ' || i,
                           (RANDOM() * 100 + 10)::DECIMAL(10, 2), -- Random price between 10 and 110
                           (SELECT id FROM categories ORDER BY RANDOM() LIMIT 1), -- Random category
                           (SELECT id FROM users ORDER BY RANDOM() LIMIT 1) -- Random user
                       );
            END LOOP;
    END;
$$;

-- Insert Mock Carts for Users
DO
$$
    BEGIN
        FOR i IN 1..50 LOOP
                INSERT INTO carts (user_id)
                VALUES (
                           (SELECT id FROM users ORDER BY RANDOM() LIMIT 1) -- Random user
                       );
            END LOOP;
    END;
$$;

-- Insert Mock Cart Items
DO
$$
    BEGIN
        FOR i IN 1..200 LOOP
                INSERT INTO cart_items (cart_id, product_id, quantity)
                VALUES (
                           (SELECT id FROM carts ORDER BY RANDOM() LIMIT 1), -- Random cart
                           (SELECT id FROM products ORDER BY RANDOM() LIMIT 1), -- Random product
                           (FLOOR(RANDOM() * 5 + 1)::INT) -- Random quantity between 1 and 5
                       );
            END LOOP;
    END;
$$;

-- Insert Mock Orders for Users
DO
$$
    BEGIN
        FOR i IN 1..50 LOOP
                INSERT INTO orders (user_id, total_amount, status, payment_status)
                VALUES (
                           (SELECT id FROM users ORDER BY RANDOM() LIMIT 1), -- Random user
                           (RANDOM() * 500 + 50)::DECIMAL(10, 2), -- Random total amount between 50 and 550
                           'PENDING',
                           'UNPAID'
                       );
            END LOOP;
    END;
$$;

-- Insert Mock Order Items
DO
$$
    BEGIN
        FOR i IN 1..200 LOOP
                INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
                VALUES (
                           (SELECT id FROM orders ORDER BY RANDOM() LIMIT 1), -- Random order
                           (SELECT id FROM products ORDER BY RANDOM() LIMIT 1), -- Random product
                           (FLOOR(RANDOM() * 5 + 1)::INT), -- Random quantity between 1 and 5
                           (RANDOM() * 100 + 10)::DECIMAL(10, 2) -- Random price at the time of order
                       );
            END LOOP;
    END;
$$;
