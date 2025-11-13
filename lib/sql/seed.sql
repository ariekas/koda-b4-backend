INSERT INTO profile (pic, phone, address) VALUES
('https://i.pravatar.cc/150?img=1', '081234567890', 'Jl. Merdeka 1, Jakarta'),
('https://i.pravatar.cc/150?img=2', '081234567891', 'Jl. Sudirman 2, Jakarta'),
('https://i.pravatar.cc/150?img=3', '081234567892', 'Jl. Thamrin 3, Jakarta'),
('https://i.pravatar.cc/150?img=4', '081234567893', 'Jl. Gatot Subroto 4, Jakarta'),
('https://i.pravatar.cc/150?img=5', '081234567894', 'Jl. Diponegoro 5, Jakarta');

INSERT INTO users (fullname, email, password, role, profile_id) VALUES
('Ari Eka Saputra', 'ari@example.com', 'password123', 'user', 1),
('Rina Putri', 'rina@example.com', 'password123', 'user', 2),
('Budi Santoso', 'budi@example.com', 'password123', 'user', 3),
('Siti Aminah', 'siti@example.com', 'password123', 'user', 4),
('Andi Wijaya', 'andi@example.com', 'password123', 'admin', 5);

INSERT INTO category_products (name) VALUES
('Elektronik'),
('Fashion'),
('Kecantikan'),
('Olahraga'),
('Makanan & Minuman');

INSERT INTO discounts (name, diskon) VALUES
('diskon 5rb', 5000),
('natal', 10000);

INSERT INTO products (name, price, description, stock,price_discounts, discounts_id, is_flashsale, is_favorite_product, category_products_id) VALUES
('Laptop Lenovo', 10000000, 'Laptop performa tinggi', 10,0,1, FALSE, TRUE, 1),
('Sepatu Adidas', 800000, 'Sepatu olahraga nyaman', 20,0,2, TRUE, FALSE, 2),
('Lipstik Wardah', 120000, 'Lipstik tahan lama', 50,0,1, FALSE, TRUE, 3),
('Bola Sepak', 250000, 'Bola resmi liga', 15,0,1, FALSE, FALSE, 4),
('Kopi Luwak', 500000, 'Kopi premium', 30,0,2, TRUE, TRUE, 5);

INSERT INTO product_images (image, products_id) VALUES
('https://picsum.photos/200?1', 1),
('https://picsum.photos/200?2', 2),
('https://picsum.photos/200?3', 3),
('https://picsum.photos/200?4', 4),
('https://picsum.photos/200?5', 5);

INSERT INTO size_products (name, additional_costs) VALUES
('Small', 0),
('Medium', 5000),
('Large', 10000),
('Extra Large', 15000),
('XXL', 20000);

INSERT INTO variant_products (name, additional_costs) VALUES
('Merah', 0),
('Biru', 0),
('Hitam', 10000),
('Putih', 5000),
('Kuning', 5000);

INSERT INTO variant_products (name, additional_costs) VALUES
('Merah', 0),
('Biru', 0),
('Hitam', 10000),
('Putih', 5000),
('Kuning', 5000);

INSERT INTO carts (users_id, products_id, size_products_id, variant_products_id, quantity) VALUES
(1, 1, 2, 3, 1),
(2, 2, 1, 2, 2),
(3, 3, 1, 1, 3),
(4, 4, 3, 4, 1),
(5, 5, 2, 5, 2);

INSERT INTO deliverys (name, price) VALUES
('JNE', 20000),
('TIKI', 15000),
('SiCepat', 10000),
('Gojek Instant', 5000),
('GrabExpress', 7000);

INSERT INTO payment_methods (name, image_payment) VALUES
('BCA', 'https://logo.clearbit.com/bca.co.id'),
('Mandiri', 'https://logo.clearbit.com/mandiri.co.id'),
('OVO', 'https://logo.clearbit.com/ovo.id'),
('Dana', 'https://logo.clearbit.com/dana.id'),
('COD', NULL);

INSERT INTO status_transactions (status) VALUES 
('complete'),
('Pending'),
('Cancel');


INSERT INTO transactions (users_id, deliverys_id, payment_methods_id, name_user, address_user, phone_user, email_user, total, status_transactions_id, invoice_num) VALUES
(1, 1, 1, 'Ari Eka', 'Jl. Merdeka 1, Jakarta', '081234567890', 'ari@example.com', 10500000, 1, 'INV-0001'),
(2, 2, 2, 'Rina Putri', 'Jl. Sudirman 2, Jakarta', '081234567891', 'rina@example.com', 810000, 2, 'INV-0002'),
(3, 3, 3, 'Budi Santoso', 'Jl. Thamrin 3, Jakarta', '081234567892', 'budi@example.com', 120000, 1, 'INV-0003'),
(4, 4, 4, 'Siti Aminah', 'Jl. Gatot Subroto 4, Jakarta', '081234567893', 'siti@example.com', 255000,3, 'INV-0004'),
(5, 5, 5, 'Andi Wijaya', 'Jl. Diponegoro 5, Jakarta', '081234567894', 'andi@example.com', 507000, 1, 'INV-0005');

INSERT INTO transaction_items (transactions_id, products_id, quantity, subtotal, variant_products_id, size_products_id) VALUES
(1, 1, 1, 10500000, 3, 2),
(2, 2, 2, 810000, 2, 1),
(3, 3, 3, 120000, 1, 1),
(4, 4, 1, 255000, 4, 3),
(5, 5, 2, 507000, 5, 2);

INSERT INTO ratings (rating, review) VALUES
(5, 'Laptop sangat cepat dan bagus!'),
(4, 'Sepatunya nyaman dipakai.'),
(5, 'Lipstik tahan lama dan warnanya cantik.'),
(3, 'Bola sepak agak keras.'),
(5, 'Kopi luwak enak dan wangi.');
