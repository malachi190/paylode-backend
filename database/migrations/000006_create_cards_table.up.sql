CREATE TABLE
    IF NOT EXISTS cards (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        user_id BIGINT NOT NULL,
        brand VARCHAR(20),
        last4 CHAR(4),
        token VARCHAR(64) UNIQUE,
        exp_month TINYINT,
        exp_year SMALLINT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );