INSERT INTO mandarine.public.users (username, email, password, role_id, is_enabled, is_email_verified, is_password_temp,
                                    created_at, updated_at, deleted_at)
    (SELECT 'test_user_' || id,
            'test_user_' || id || '@mandarine.com',
            '$2a$12$DmCF/Mw9t4/wtBsjNtUd6.lCNRIkemlztOfxgNqeRYCQZCSjrquL.', -- test
            (SELECT id FROM mandarine.public.roles WHERE name = 'admin'),
            true,
            true,
            false,
            NOW(),
            NOW(),
            null
    FROM generate_series(0, 2048) as id)
ON CONFLICT DO NOTHING;