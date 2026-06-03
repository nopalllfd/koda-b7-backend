-- db/seeds/payment_methods.sql

INSERT INTO payment_methods (name, logo)
VALUES
('BCA','bca.png'),
('BNI', 'bni.png'),
('BRI', 'bri.png'),
('Mandiri', 'mandiri.png'),
('DANA', 'dana.png'),
('OVO', 'ovo.png')
ON CONFLICT (name) DO NOTHING;