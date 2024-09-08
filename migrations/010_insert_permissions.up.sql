-- 1. Insert Users
INSERT INTO users (username, password, email, is_blocked, login_attempts)
VALUES ('admin', '$2a$10$za4y5cGJX5rmTBjAytjvm.Br4g5t4iGivpsoZqvqmAYeP5LWpRVC6', 'admin@example.com', false, 0),
       ('vendor', '$2a$10$za4y5cGJX5rmTBjAytjvm.Br4g5t4iGivpsoZqvqmAYeP5LWpRVC6', 'vendor@example.com', false, 0),
       ('customer', '$2a$10$za4y5cGJX5rmTBjAytjvm.Br4g5t4iGivpsoZqvqmAYeP5LWpRVC6', 'customer@example.com', false, 0),
       ('guest', '$2a$10$za4y5cGJX5rmTBjAytjvm.Br4g5t4iGivpsoZqvqmAYeP5LWpRVC6', 'guest@example.com', false, 0);

-- 2. Insert Roles
INSERT INTO roles (name, description)
VALUES ('Admin', 'Full access to manage users, products, categories, and orders'),
       ('Vendor', 'Manage own products and orders'),
       ('Customer', 'Browse products, manage cart, and place orders'),
       ('Guest', 'Limited access to browse products');

-- 3. Insert Permissions
INSERT INTO permissions (name, description)
VALUES ('CREATE_PRODUCT', 'Permission to create a product'),
       ('VIEW_PRODUCT', 'Permission to view a product'),
       ('DELETE_PRODUCT', 'Permission to delete a product'),
       ('MANAGE_USERS', 'Permission to manage users'),
       ('PLACE_ORDER', 'Permission to place an order');

-- 4. Assign Permissions to Roles
INSERT INTO role_permissions (role_id, permission_id)
VALUES (1, 1), -- Admin can CREATE_PRODUCT
       (1, 2), -- Admin can VIEW_PRODUCT
       (1, 3), -- Admin can DELETE_PRODUCT
       (1, 4), -- Admin can MANAGE_USERS
       (2, 1), -- Vendor can CREATE_PRODUCT
       (2, 2), -- Vendor can VIEW_PRODUCT
       (3, 2), -- Customer can VIEW_PRODUCT
       (3, 5), -- Customer can PLACE_ORDER
       (4, 2);
-- Guest can VIEW_PRODUCT

-- 5. Assign Roles to Users
-- Retrieve user IDs dynamically using subqueries to ensure the correct user_ids
INSERT INTO user_roles (user_id, role_id)
VALUES ((SELECT id FROM users WHERE username = 'admin'), 1),    -- Admin Role
       ((SELECT id FROM users WHERE username = 'vendor'), 2),   -- Vendor Role
       ((SELECT id FROM users WHERE username = 'customer'), 3), -- Customer Role
       ((SELECT id FROM users WHERE username = 'guest'), 4); -- Guest Role
