-- ✅ Seed для ЛР 6: все таблицы, без сложных ON CONFLICT

-- 1. Users
INSERT INTO users (login, password, is_moderator) VALUES
                                                      ('admin', 'admin123', true),
                                                      ('user1', 'user123', false),
                                                      ('user2', 'user123', false)
    ON CONFLICT (login) DO NOTHING;

-- 2. Tires
DELETE FROM tires WHERE tire_title IN (
                                       'Tunga Nordway', 'mrl-mr3', 'MITAS_TR-08', 'Michelin_Latitude_Tour'
    );
INSERT INTO tires (tire_title, tire_material_coefficient, tire_thickness_coefficient, description, short_description_en, photo, video, is_delete) VALUES
                                                                                                                                                      ('Tunga Nordway', 1.3, 0.9, 'Зимняя шипованная шина с агрессивным V-образным протектором. Усиленный каркас, отличное сцепление на льду и укатанном снегу. Подходит для легковых автомобилей и кроссоверов.', 'Studded winter tire with aggressive V-tread for ice roads', 'Tunga_Nordway.jpg', '', false),
                                                                                                                                                      ('mrl-mr3', 2.1, 1.8, 'Всесезонная широкопрофильная шина для спецтехники. Усиленный каркас, стойкость к порезам и ударам. Оптимизирована для работы на твёрдых и смешанных покрытиях.', 'Wide-profile all-season tire for heavy trucks and machinery', 'mrl-mr3.jpg', '', false),
                                                                                                                                                      ('MITAS_TR-08', 1.9, 1.6, 'Всесезонная арочная шина для сельскохозяйственной техники. Гибкая боковина, увеличенное пятно контакта, минимальное давление на почву. Для полевых работ.', 'Arched agricultural tire with flexible sidewall for soft soil', 'MITAS_TR-08.jpg', '', false),
                                                                                                                                                      ('Michelin_Latitude_Tour', 1.1, 0.7, 'Летняя шина для кроссоверов и внедорожников. Асимметричный протектор, низкий уровень шума, топливная экономичность. Комфорт на трассе и в городе.', 'Summer highway tire with low noise and fuel efficient design', 'Michelin_Latitude_Tour.jpg', '', false);

-- 3. TirePressures
DELETE FROM tire_pressures WHERE tire_pressure_id IN (1, 2);
INSERT INTO tire_pressures (tire_pressure_id, status, date_create, date_formed, date_completed, creator_id, moderator_id, air_temperature, car_weight) VALUES
                                                                                                                                                           (1, 'сформирован', '2025-01-15 10:00:00', '2025-02-01 14:30:00', NULL, 2, NULL, 20.0, 1500.0),
                                                                                                                                                           (2, 'завершён', '2024-12-01 09:00:00', '2024-12-15 11:00:00', '2025-01-20 16:00:00', 2, 1, -5.0, 1800.0);

-- 4. TirePressureEntries (без ON CONFLICT)
DELETE FROM tire_pressure_entries WHERE tire_pressure_id IN (1, 2);
INSERT INTO tire_pressure_entries (tire_pressure_id, tire_id, coating_coefficient, pressure) VALUES
                                                                                                 (1, 1, 1.0, 220.0),
                                                                                                 (1, 4, 1.1, 210.0),
                                                                                                 (2, 2, 1.3, 240.0),
                                                                                                 (2, 3, 1.2, 230.0);