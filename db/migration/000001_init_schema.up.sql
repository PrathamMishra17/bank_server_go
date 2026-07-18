-- CREATE DATABASE BankAccount;


CREATE TABLE accounts (
	id INT AUTO_INCREMENT PRIMARY KEY,
	owner VARCHAR(255) NOT NULL,
	balance INT NOT NULL,
	currency VARCHAR(255) NOT NULL,
	createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE entries(
	id INT AUTO_INCREMENT PRIMARY KEY,
	account_id INT NOT NULL,
	amount INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers(
	id INT AUTO_INCREMENT PRIMARY KEY,
	from_account_id INT NOT NULL,
	to_account_id INT NOT NULL,
	amount INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- SETTING THE INDEXES

CREATE INDEX idx_account ON accounts(owner);
CREATE INDEX idx_account_id ON entries(account_id);
CREATE INDEX id_from_account ON transfers(from_account_id);
CREATE INDEX id_to_account ON transfers(to_account_id);

-- setting the foreign key

ALTER TABLE entries ADD CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE transfers ADD CONSTRAINT fk_from_account FOREIGN KEY (from_account_id) REFERENCES accounts(id);
ALTER TABLE transfers ADD CONSTRAINT fk_to_account FOREIGN KEY (to_account_id) REFERENCES accounts(id);

