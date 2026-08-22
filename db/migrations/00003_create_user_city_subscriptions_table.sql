-- +goose Up
CREATE TABLE user_city_subscriptions (
    user_id    BIGINT UNSIGNED NOT NULL,
    city_id    BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(6)     NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (user_id, city_id),
    KEY idx_user_city_subscriptions_city_id (city_id),
    CONSTRAINT fk_user_city_subscriptions_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_city_subscriptions_city_id
        FOREIGN KEY (city_id) REFERENCES cities (id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE user_city_subscriptions;
