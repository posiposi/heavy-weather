-- +goose Up
CREATE TABLE cities (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name       VARCHAR(255)    NOT NULL,
    latitude   DECIMAL(8,5)    NOT NULL,
    longitude  DECIMAL(8,5)    NOT NULL,
    timezone   VARCHAR(64)     NOT NULL DEFAULT 'Asia/Tokyo',
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_cities_name (name),
    UNIQUE KEY uk_cities_coordinates (latitude, longitude)
);

-- +goose Down
DROP TABLE cities;
