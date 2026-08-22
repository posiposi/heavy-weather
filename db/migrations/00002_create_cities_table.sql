-- +goose Up
CREATE TABLE cities (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name       VARCHAR(255)    NOT NULL,
    latitude   DECIMAL(8,5)    NOT NULL,
    longitude  DECIMAL(8,5)    NOT NULL,
    timezone   VARCHAR(64)     NOT NULL DEFAULT 'Asia/Tokyo',
    created_at DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_cities_name (name),
    UNIQUE KEY uk_cities_coordinates (latitude, longitude)
);

-- +goose Down
DROP TABLE cities;
