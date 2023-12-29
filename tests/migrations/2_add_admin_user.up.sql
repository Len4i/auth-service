INSERT INTO
    users (id, email, pass_hash, is_admin)
VALUES (
        999,
        'admin-user@localhost.com',
        'some-hash',
        1
    ) ON CONFLICT DO NOTHING;