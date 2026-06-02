CREATE TABLE topups (
transaction_id INTEGER PRIMARY KEY,
wallet_id INTEGER NOT NULL,
method_id INTEGER NOT NULL,
amount NUMERIC(15,2) NOT NULL,
admin_fee NUMERIC(15,2),
total NUMERIC(15,2) NOT NULL,
payment_reference VARCHAR(100),

CONSTRAINT topups_transaction_id_fkey
    FOREIGN KEY (transaction_id)
    REFERENCES transactions(id)
    ON DELETE CASCADE,

CONSTRAINT topups_wallet_id_fkey
    FOREIGN KEY (wallet_id)
    REFERENCES wallets(id)
    ON DELETE CASCADE,

CONSTRAINT topups_method_id_fkey
    FOREIGN KEY (method_id)
    REFERENCES payment_methods(id)
    ON DELETE CASCADE

);