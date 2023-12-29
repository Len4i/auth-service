INSERT INTO
    apps (id, name, secret)
VALUES (999, 'test-app', 'test-secret') ON CONFLICT DO NOTHING;