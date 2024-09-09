-- Roles Table: Defines different user roles in the system
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) UNIQUE NOT NULL, -- e.g., 'Admin', 'Vendor', 'Customer', 'Guest'
                       description TEXT
);

-- Permissions Table: Defines permissions associated with roles
CREATE TABLE permissions (
                             id SERIAL PRIMARY KEY,
                             name VARCHAR(50) UNIQUE NOT NULL, -- e.g., 'CREATE_PRODUCT', 'VIEW_PRODUCT'
                             description TEXT
);

-- Role_Permissions Table: Associates permissions with roles
CREATE TABLE role_permissions (
                                  id SERIAL PRIMARY KEY,
                                  role_id INT NOT NULL,
                                  permission_id INT NOT NULL,
                                  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
                                  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- User_Roles Table: Associates users with roles
CREATE TABLE user_roles (
                            id SERIAL PRIMARY KEY,
                            user_id INT NOT NULL,
                            role_id INT NOT NULL,
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
