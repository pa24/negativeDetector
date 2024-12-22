-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
                                        id SERIAL PRIMARY KEY,
                                        user_id BIGINT NOT NULL,
                                        username TEXT,
                                        content_type TEXT NOT NULL,
                                        content TEXT,
                                        chat_id BIGINT NOT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS messages;
