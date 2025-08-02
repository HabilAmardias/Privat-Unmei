INSERT INTO admins (id)
VALUES
((SELECT id FROM users where email = 'admin1@mail.com'));