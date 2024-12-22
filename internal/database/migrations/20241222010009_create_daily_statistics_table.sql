-- +goose Up
CREATE TABLE IF NOT EXISTS daily_statistics (
                                                id SERIAL PRIMARY KEY,
                                                date DATE NOT NULL UNIQUE,
                                                total_messages INT DEFAULT 0,
                                                total_words INT DEFAULT 0,
                                                top_user_id BIGINT,
                                                top_user_messages INT DEFAULT 0,
                                                top_voice_user_id BIGINT,
                                                top_voice_count INT DEFAULT 0,
                                                top_video_user_id BIGINT,
                                                top_video_count INT DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS daily_statistics;
