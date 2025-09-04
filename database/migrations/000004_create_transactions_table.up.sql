CREATE TABLE
    IF NOT EXISTS transactions (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        user_id BIGINT NOT NULL,
        transaction_type VARCHAR(50) NOT NULL, -- e.g debit, credit
        reference VARCHAR(250) NOT NULL,
        payment_method VARCHAR(50) NOT NULL, -- e.g card
        status VARCHAR(20) DEFAULT "pending",
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );