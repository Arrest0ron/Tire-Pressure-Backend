-- Seed data for tire pressure application
-- Run this in Adminer

-- Insert Users
INSERT INTO users (login, password, is_moderator) VALUES
                                                      ('admin', 'admin123', true),
                                                      ('user1', 'user123', false),
                                                      ('user2', 'user123', false)
    ON CONFLICT (login) DO NOTHING;

-- Insert Tires (Services)
INSERT INTO tires (tire_title, tire_material_coefficient, tire_thickness_coefficient, description, photo, video, is_delete) VALUES
                                                                                                                                ('Летняя шина Michelin', 1.2, 8.5, 'Летняя шина для легковых автомобилей, хороший сцепление на сухом асфальте', '', '', false),
                                                                                                                                ('Зимняя шина Nokian', 1.5, 9.0, 'Зимняя шипованная шина для снежных и ледяных дорог', '', '', false),
                                                                                                                                ('Всесезонная шина Bridgestone', 1.3, 8.0, 'Универсальная шина для умеренного климата', '', '', false),
                                                                                                                                ('Грузовая шина Continental', 1.8, 12.0, 'Шина для грузовых автомобилей и фургонов', '', '', false),
                                                                                                                                ('Спортивная шина Pirelli', 1.1, 7.5, 'Высокопроизводительная шина для спортивных автомобилей', '', '', false)
    ON CONFLICT DO NOTHING;

-- Insert TirePressures (Requests)
INSERT INTO tire_pressures (status, date_create, date_formed, date_completed, creator_id, moderator_id, air_temperature, car_weight) VALUES
                                                                                                                                         ('сформирован', '2025-01-15 10:00:00', '2025-02-01 14:30:00', NULL, 2, NULL, 20.0, 1500.0),
                                                                                                                                         ('завершён', '2024-12-01 09:00:00', '2024-12-15 11:00:00', '2025-01-20 16:00:00', 2, 1, -5.0, 1800.0)
    ON CONFLICT DO NOTHING;

-- Insert TirePressureEntries (M-M)
INSERT INTO tire_pressure_entries (tire_pressure_id, tire_id, coating_coefficient, pressure) VALUES
                                                                                                 (1, 1, 1.0, 220.0),
                                                                                                 (1, 2, 1.2, 240.0),
                                                                                                 (2, 3, 1.1, 230.0),
                                                                                                 (2, 4, 1.5, 280.0)
    ON CONFLICT DO NOTHING;