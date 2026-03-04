-- +goose Up
-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE
);

-- Создаем таблицу шин
CREATE TABLE IF NOT EXISTS tires (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    photo VARCHAR(255),
    video VARCHAR(255),
    tire_coefficient DECIMAL(10,2) NOT NULL,
    is_delete BOOLEAN DEFAULT FALSE
);

-- Создаем таблицу заявок на расчёт давления
CREATE TABLE IF NOT EXISTS tire_pressures (
    id SERIAL PRIMARY KEY,
    air_temperature FLOAT NOT NULL,
    car_weight FLOAT NOT NULL,
    date_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    creator_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    total_pressure FLOAT DEFAULT 0
);

-- Создаем таблицу записей в заявках на расчёт давления
CREATE TABLE IF NOT EXISTS tire_pressure_entries (
    id SERIAL PRIMARY KEY,
    tire_pressure_id INTEGER NOT NULL REFERENCES tire_pressures(id) ON DELETE CASCADE,
    tire_id INTEGER NOT NULL REFERENCES tires(id) ON DELETE CASCADE,
    coating_coeff FLOAT NOT NULL,
    pressure FLOAT NOT NULL,
    CONSTRAINT unique_tire_pressure_tire UNIQUE (tire_pressure_id, tire_id)
);

-- Добавляем ограничения и вычисляемые поля
ALTER TABLE tire_pressures
ADD CONSTRAINT chk_status CHECK (status IN ('draft', 'deleted', 'formed', 'completed', 'rejected'));


-- Добавляем тестового пользователя
INSERT INTO users (login, password, is_moderator) VALUES
('admin', 'password123', true);

-- Создаем 8 шин
INSERT INTO tires (title, description, photo, video, tire_coefficient, is_delete) VALUES
('Michelin Pilot Sport 4S', 'Высокопроизводительная летняя шина для спортивных автомобилей', 'Michelin Pilot Sport 4S.jpg', 'yellow_black.mp4', 2.2, false),
('Nokian Hakkapeliitta R5', 'Профессиональная зимняя шина для максимального сцепления', 'Nokian Hakkapeliitta R5.jpg', 'yellow_black.mp4', 2.3, false),
('Continental AllSeasonContact', 'Всесезонная шина премиум-класса', 'Continental AllSeasonContact.jpg', 'yellow_black.mp4', 2.2, false),
('Bridgestone Potenza RE003', 'Спортивная летняя шина для мощных седанов', 'Bridgestone Potenza RE003.jpeg', 'yellow_black.mp4', 2.3, false),
('Pirelli P Zero', 'Ультра-высокопроизводительная летняя шина', 'Pirelli Winter Cinturato.jpg', 'yellow_black.mp4', 2.4, false),
('Goodyear Eagle F1', 'Спортивная летняя шина с высоким сцеплением', 'Goodyear UltraGrip Ice 2.jpg', 'yellow_black.mp4', 2.3, false),
('Yokohama Advan Sport', 'Высокопроизводительная шина для спортивных автомобилей', 'Yokohama Advan Neova.jpg', 'yellow_black.mp4', 2.2, false),
('Hankook Ventus S1', 'Универсальная спортивная летняя шина', 'Vredestein Quatrac Pro.jpg', 'yellow_black.mp4', 2.1, false);

-- Создаем 2 заявки на расчёт давления
INSERT INTO tire_pressures (air_temperature, car_weight, status, creator_id) VALUES
(20.0, 1500.0, 'draft', 1),
(15.0, 1800.0, 'completed', 1);

-- Обновляем вторую заявку для расчёта total_pressure
UPDATE tire_pressures
SET total_pressure = 5.75
WHERE id = 2;

-- Добавляем 3 шины в первую заявку (ID=1)
INSERT INTO tire_pressure_entries (tire_pressure_id, tire_id, coating_coeff, pressure) VALUES
(1, 1, 0.8, 0),
(1, 2, 0.9, 0),
(1, 3, 0.7, 0);

-- Добавляем 3 шины во вторую заявку (ID=2)
INSERT INTO tire_pressure_entries (tire_pressure_id, tire_id, coating_coeff, pressure) VALUES
(2, 2, 0.85, 0),
(2, 3, 0.75, 0),
(2, 4, 0.9, 0);
