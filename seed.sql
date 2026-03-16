-- Seed data for dendrochronology application
-- Run this in Adminer

-- Insert Users
INSERT INTO users (login, password, is_moderator) VALUES 
    ('admin', 'admin123', true),
    ('user1', 'user123', false),
    ('user2', 'user123', false)
ON CONFLICT (login) DO NOTHING;

-- Insert Constructions (Services)
INSERT INTO constructions (construction_title, use_life, description, image_url, video_url, is_delete) VALUES 
    ('Бревенчатый сруб', '50-100 лет', 'Бревенчатый сруб из сосны или ели для жилых домов', '', '', false),
    ('Балка перекрытия', '30-50 лет', 'Деревянная балка для межэтажного перекрытия', '', '', false),
    ('Стропильная система', '20-40 лет', 'Стропила для кровли деревянного дома', '', '', false),
    ('Колотый бревно', '40-80 лет', 'Колотые бревна для строительства бань', '', '', false),
    ('Дощатый настил', '15-30 лет', 'Дощатый пол на лагах', '', '', false)
ON CONFLICT DO NOTHING;

-- Insert Dendrochronology (Requests)
INSERT INTO dendrochronologies (status, date_create, date_formed, date_completed, creator_id, moderator_id, total_samples, build_date) VALUES 
    ('сформирован', '2025-01-15 10:00:00', '2025-02-01 14:30:00', NULL, 2, NULL, NULL, NULL),
    ('завершён', '2024-12-01 09:00:00', '2024-12-15 11:00:00', '2025-01-20 16:00:00', 2, 1, 5, 1870)
ON CONFLICT DO NOTHING;

-- Insert Dendrochronology-Construction (M-M)
INSERT INTO dendrochronology_constructions (dendrochronology_id, construction_id, samples_count, cutting_date, date_correction) VALUES 
    (1, 1, 3, '1850', '10'),
    (1, 2, 2, '1860', '5')
ON CONFLICT DO NOTHING;