-- ==========================================
-- BIKIN INDUK TRANSAKSI (Pakai ID 7, 8, 9)
-- ==========================================
INSERT INTO transactions (id, type, reference_code, status, created_at, updated_at)
VALUES 
(7, 'topup', 'TOP-007', 'success', NOW() - INTERVAL '2 days', NOW()),
(8, 'transfer', 'TRF-008', 'success', NOW() - INTERVAL '1 day', NOW()),
(9, 'transfer', 'TRF-009', 'success', NOW(), NOW());


-- ==========================================
-- DETAIL TRANSAKSI UNTUK DOMPET ID 3
-- ==========================================

-- 1. Mutasi Masuk: Top Up via Mandiri VA
INSERT INTO topups (transaction_id, wallet_id, method_id, amount, admin_fee, total, payment_reference, expired_at, paid_at)
VALUES (7, 3, 2, 300000.00, 2500.00, 302500.00, 'PAY-MANDIRI-777', NOW() + INTERVAL '1 day', NOW());

-- 2. Mutasi Masuk: Transfer dari User 1 ke User 3
-- (Udah pakai struktur tabel transfers yang simpel: cuma amount & description)
INSERT INTO transfers (transaction_id, sender_wallet_id, receiver_wallet_id, amount, description)
VALUES (8, 1, 3, 150000.00, 'Bayar desain UI/UX');

-- 3. Mutasi Keluar: Transfer dari User 3 ke User 2
INSERT INTO transfers (transaction_id, sender_wallet_id, receiver_wallet_id, amount, description)
VALUES (9, 3, 2, 50000.00, 'Ganti uang ngopi');

ALTER TABLE topups ADD COLUMN tax_amount NUMERIC(15,2) DEFAULT 0;