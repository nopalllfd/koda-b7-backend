CREATE TABLE transfers (
transaction_id INTEGER PRIMARY KEY,
sender_wallet_id INTEGER NOT NULL,
receiver_wallet_id INTEGER NOT NULL,
amount NUMERIC(15,2) NOT NULL,
description TEXT,

CONSTRAINT fk_transfer_transaction
    FOREIGN KEY (transaction_id)
    REFERENCES transactions(id)
    ON DELETE CASCADE,

CONSTRAINT fk_transfer_sender_wallet
    FOREIGN KEY (sender_wallet_id)
    REFERENCES wallets(id),

CONSTRAINT fk_transfer_receiver_wallet
    FOREIGN KEY (receiver_wallet_id)
    REFERENCES wallets(id)

);