INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator'),
        ('user', 'User')
ON CONFLICT DO NOTHING;