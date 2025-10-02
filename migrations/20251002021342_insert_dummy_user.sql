-- +goose Up
INSERT INTO roles (id, name) VALUES (1, 'superadmin');
INSERT INTO users (username, email, fullname, role_id, password) VALUES ('admin', 'admin@domain.com', 'Super Admin', 1, 'hashedpassword');

-- +goose Down
DELETE FROM users WHERE username = 'admin';
DELETE FROM roles WHERE name = 'superadmin';
