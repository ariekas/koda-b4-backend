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
('Super Supreme Original', 90000, 'Pizza dengan pepperoni, daging sapi, jamur, paprika, dan keju mozzarella', 100, 0, NULL, FALSE, TRUE, 1),
('Meat Lovers', 96000, 'Pizza penuh daging sapi, pepperoni, ham, dan sosis', 100, 0, NULL, FALSE, TRUE, 1),
('American Favourite', 96000, 'Pizza klasik Amerika dengan pepperoni dan daging sapi', 100, 0, NULL, FALSE, FALSE, 1),
('Veggie Garden', 76000, 'Pizza sayuran segar seperti paprika, jamur, bawang, dan tomat', 100, 0, NULL, FALSE, FALSE, 1),
('Hawaiian Chicken', 76000, 'Pizza dengan ayam schnitzel, nanas, paprika, dan mozzarella', 100, 0, NULL, FALSE, FALSE, 1),
('Tuna Melt', 96000, 'Pizza tuna dengan jagung dan saus mayonaise khas', 100, 0, NULL, FALSE, FALSE, 1),
('Cheesy Galore', 96000, 'Pizza dengan lapisan ekstra keju mozzarella dan saus pizza khas', 100, 0, NULL, FALSE, TRUE, 1),
('Splitza Classic', 112000, 'Splitza dengan dua rasa dalam satu loyang, rasa klasik Pizza Hut', 100, 0, NULL, FALSE, FALSE, 1),
('Splitza Signature', 155000, 'Splitza premium dengan dua topping signature di satu pizza besar', 100, 0, NULL, FALSE, TRUE, 1),
('QU4RTZA Pizza', 147273, 'Pizza dengan 4 topping berbeda dalam satu loyang', 100, 0, NULL, FALSE, TRUE, 1),
('L1MO Pizza', 322727, 'Pizza panjang besar, cocok untuk rame‑rame', 100, 0, NULL, FALSE, FALSE, 1),
('Pizza Heart Pepperoni', 139091, 'Pizza berbentuk hati dengan topping pepperoni', 100, 0, NULL, FALSE, FALSE, 1),
('Sweet Melts Chocolate', 41000, 'Pizza manis dengan topping coklat dan keju', 100, 0, NULL, FALSE, FALSE, 1),
('Melts Pizza Pepperoni', 50910, 'Melts pizza dengan pepperoni, crust tipis dan renyah', 100, 0, NULL, FALSE, FALSE, 1),
('Melts Pizza Meat Lovers', 57000, 'Melts pizza daging lengkap: sapi, pepperoni, ham', 100, 0, NULL, FALSE, TRUE, 1),
('Stuffed Crust Supreme', 124000, 'Pizza Supreme dengan pinggiran roti berisi keju', 100, 0, NULL, FALSE, TRUE, 1),
('Stuffed Crust Meaty', 124000, 'Pizza Meaty dengan crust berisi keju leleh', 100, 0, NULL, FALSE, FALSE, 1),
('Stuffed Crust Cheesy Galore', 124000, 'Keju ekstra di tengah dan pinggiran pizza', 100, 0, NULL, FALSE, TRUE, 1),
('Pan Pepperoni', 110000, 'Pan pizza dengan banyak pepperoni dan keju', 100, 0, NULL, FALSE, FALSE, 1),
('Pan Meat Monsta', 110000, 'Pan pizza jumbo dengan kombinasi daging berlimpah', 100, 0, NULL, FALSE, FALSE, 1),
('Crazy Crust BBQ Chicken', 120000, 'Pizza ayam dengan saus BBQ dan crust spesial Pizza Hut', 100, 0, NULL, FALSE, FALSE, 1),
('Crown Crust Deluxe Cheese', 88000, 'Pizza Cheese dengan crown crust penuh keju', 100, 0, NULL, FALSE, TRUE, 1),
('Crown Crust Vegan Option', 125000, 'Pizza vegan dengan sayuran dan crust crown vegan', 100, 0, NULL, FALSE, FALSE, 1),
('Black Pepper Beef', 96000, 'Pizza dengan irisan daging sapi dan lada hitam pedas', 100, 0, NULL, FALSE, FALSE, 1),
('Black Pepper Chicken', 96000, 'Pizza ayam lada hitam dengan saus khas', 100, 0, NULL, FALSE, FALSE, 1),
('Tuna Supreme', 100000, 'Pizza tuna premium dengan jamur dan paprika', 100, 0, NULL, FALSE, FALSE, 1),
('Chicken Lovers', 96000, 'Pizza penuh ayam dan keju mozzarella', 100, 0, NULL, FALSE, FALSE, 1),
('Meat Lovers Cheesy Mayo', 86000, 'Pizza daging dengan saus mayonaise dan keju meleleh', 100, 0, NULL, FALSE, TRUE, 1),
('Super Supreme Chicken', 86000, 'Super Supreme versi dengan potongan ayam', 100, 0, NULL, FALSE, FALSE, 1),
('Pepperoni Cheese Lover', 96000, 'Pizza pepperoni dengan banyak lapisan keju', 100, 0, NULL, FALSE, TRUE, 1),
('Deluxe Cheese', 95000, 'Pizza klasik extra keju dengan saus Italia', 100, 0, NULL, FALSE, TRUE, 1),
('Mushroom Lovers', 85000, 'Pizza penuh jamur segar dan mozzarella', 100, 0, NULL, FALSE, FALSE, 1),
('Seafood Delight', 110000, 'Pizza dengan udang, cumi, dan dada ayam', 100, 0, NULL, FALSE, FALSE, 1),
('Spicy Italian Sausage', 105000, 'Sausage Italia pedas dengan paprika merah', 100, 0, NULL, FALSE, FALSE, 1),
('Pepperoni Bacon', 108000, 'Pizza pepperoni dengan irisan bacon renyah', 100, 0, NULL, FALSE, TRUE, 1),
('Margarita Classic', 80000, 'Pizza klasik ala Italia dengan tomat, basil dan mozzarella', 100, 0, NULL, FALSE, TRUE, 1),
('Buffalo Chicken', 102000, 'Pizza ayam dengan saus buffalo pedas dan keju', 100, 0, NULL, FALSE, FALSE, 1),
('BBQ Beef Steak', 112000, 'Pizza dengan daging steak sapi dan saus BBQ khas', 100, 0, NULL, FALSE, FALSE, 1),
('Four Cheese', 115000, 'Pizza dengan kombinasi keju mozzarella, cheddar, parmesan, dan ricotta', 100, 0, NULL, FALSE, TRUE, 1),
('Spinach & Feta', 90000, 'Pizza sayuran bayam dengan keju feta dan saus putih', 100, 0, NULL, FALSE, FALSE, 1),
('Mie Supreme', 95000, 'Pizza gaya fusion dengan topping mie instan dan mozzarella', 100, 0, NULL, FALSE, FALSE, 1),
('Classic Margherita Heart', 139000, 'Pizza hati klasik Margherita untuk pasangan atau hadiah', 100, 0, NULL, FALSE, FALSE, 1),
('Chocolate Dessert Pizza', 45000, 'Pizza manis dengan cokelat, marshmallow, dan kacang', 100, 0, NULL, FALSE, TRUE, 1),
('Apple Cinnamon Pizza', 42000, 'Pizza manis dengan apel, kayu manis, dan gula bubuk', 100, 0, NULL, FALSE, FALSE, 1),
('Banana Caramel Pizza', 43000, 'Pizza dessert dengan pisang, saus karamel dan kacang', 100, 0, NULL, FALSE, FALSE, 1),
('Mixed Berry Pizza', 46000, 'Pizza manis dengan blueberry, strawberry, dan keju krim', 100, 0, NULL, FALSE, FALSE, 1),
('Mexican Fiesta', 108000, 'Pizza dengan daging sapi, paprika, bawang, jalapeno, dan saus pedas', 100, 0, NULL, FALSE, FALSE, 1),
('Taco Pizza', 106000, 'Pizza dengan topping taco daging, selada, dan saus taco', 100, 0, NULL, FALSE, FALSE, 1),
('Chicken Bacon Ranch', 104000, 'Pizza ayam, bacon, dan saus ranch creamy', 100, 0, NULL, FALSE, FALSE, 1),
('Supreme Deluxe', 118000, 'Pizza deluxe dengan daging, sayur, jamur, dan keju premium', 100, 0, NULL, FALSE, TRUE, 1),
('Four Seasons', 120000, 'Pizza dengan 4 kuadran topping: daging, sayur, jamur, dan keju', 100, 0, NULL, FALSE, FALSE, 1);

INSERT INTO product_images (image, products_id) VALUES
('https://picsum.photos/300/200?random=1', 11),
('https://picsum.photos/300/200?random=2', 11),
('https://picsum.photos/300/200?random=3', 11),

('https://picsum.photos/300/200?random=4', 12),
('https://picsum.photos/300/200?random=5', 12),
('https://picsum.photos/300/200?random=6', 12),

('https://picsum.photos/300/200?random=7', 13),
('https://picsum.photos/300/200?random=8', 13),
('https://picsum.photos/300/200?random=9', 13),

('https://picsum.photos/300/200?random=10', 14),
('https://picsum.photos/300/200?random=11', 14),
('https://picsum.photos/300/200?random=12', 14),

('https://picsum.photos/300/200?random=13', 15),
('https://picsum.photos/300/200?random=14', 15),
('https://picsum.photos/300/200?random=15', 15),

('https://picsum.photos/300/200?random=16', 16),
('https://picsum.photos/300/200?random=17', 16),
('https://picsum.photos/300/200?random=18', 16),

('https://picsum.photos/300/200?random=19', 17),
('https://picsum.photos/300/200?random=20', 17),
('https://picsum.photos/300/200?random=21', 17),

('https://picsum.photos/300/200?random=22', 18),
('https://picsum.photos/300/200?random=23', 18),
('https://picsum.photos/300/200?random=24', 18),

('https://picsum.photos/300/200?random=25', 19),
('https://picsum.photos/300/200?random=26', 19),
('https://picsum.photos/300/200?random=27', 19),

('https://picsum.photos/300/200?random=28', 20),
('https://picsum.photos/300/200?random=29', 20),
('https://picsum.photos/300/200?random=30', 20),

('https://picsum.photos/300/200?random=31', 21),
('https://picsum.photos/300/200?random=32', 21),
('https://picsum.photos/300/200?random=33', 21),

('https://picsum.photos/300/200?random=34', 22),
('https://picsum.photos/300/200?random=35', 22),
('https://picsum.photos/300/200?random=36', 22),

('https://picsum.photos/300/200?random=37', 23),
('https://picsum.photos/300/200?random=38', 23),
('https://picsum.photos/300/200?random=39', 23),

('https://picsum.photos/300/200?random=40', 24),
('https://picsum.photos/300/200?random=41', 24),
('https://picsum.photos/300/200?random=42', 24),

('https://picsum.photos/300/200?random=43', 25),
('https://picsum.photos/300/200?random=44', 25),
('https://picsum.photos/300/200?random=45', 25),

('https://picsum.photos/300/200?random=46', 26),
('https://picsum.photos/300/200?random=47', 26),
('https://picsum.photos/300/200?random=48', 26),

('https://picsum.photos/300/200?random=49', 27),
('https://picsum.photos/300/200?random=50', 27),
('https://picsum.photos/300/200?random=51', 27),

('https://picsum.photos/300/200?random=52', 28),
('https://picsum.photos/300/200?random=53', 28),
('https://picsum.photos/300/200?random=54', 28),

('https://picsum.photos/300/200?random=55', 29),
('https://picsum.photos/300/200?random=56', 29),
('https://picsum.photos/300/200?random=57', 29),

('https://picsum.photos/300/200?random=58', 30),
('https://picsum.photos/300/200?random=59', 30),
('https://picsum.photos/300/200?random=60', 30),

('https://picsum.photos/300/200?random=61', 31),
('https://picsum.photos/300/200?random=62', 31),
('https://picsum.photos/300/200?random=63', 31),

('https://picsum.photos/300/200?random=64', 32),
('https://picsum.photos/300/200?random=65', 32),
('https://picsum.photos/300/200?random=66', 32),

('https://picsum.photos/300/200?random=67', 33),
('https://picsum.photos/300/200?random=68', 33),
('https://picsum.photos/300/200?random=69', 33),

('https://picsum.photos/300/200?random=70', 34),
('https://picsum.photos/300/200?random=71', 34),
('https://picsum.photos/300/200?random=72', 34),

('https://picsum.photos/300/200?random=73', 35),
('https://picsum.photos/300/200?random=74', 35),
('https://picsum.photos/300/200?random=75', 35),

('https://picsum.photos/300/200?random=76', 36),
('https://picsum.photos/300/200?random=77', 36),
('https://picsum.photos/300/200?random=78', 36),

('https://picsum.photos/300/200?random=79', 37),
('https://picsum.photos/300/200?random=80', 37),
('https://picsum.photos/300/200?random=81', 37),

('https://picsum.photos/300/200?random=82', 38),
('https://picsum.photos/300/200?random=83', 38),
('https://picsum.photos/300/200?random=84', 38),

('https://picsum.photos/300/200?random=85', 39),
('https://picsum.photos/300/200?random=86', 39),
('https://picsum.photos/300/200?random=87', 39),

('https://picsum.photos/300/200?random=88', 40),
('https://picsum.photos/300/200?random=89', 40),
('https://picsum.photos/300/200?random=90', 40),

('https://picsum.photos/300/200?random=91', 41),
('https://picsum.photos/300/200?random=92', 41),
('https://picsum.photos/300/200?random=93', 41),

('https://picsum.photos/300/200?random=94', 42),
('https://picsum.photos/300/200?random=95', 42),
('https://picsum.photos/300/200?random=96', 42),

('https://picsum.photos/300/200?random=97', 43),
('https://picsum.photos/300/200?random=98', 43),
('https://picsum.photos/300/200?random=99', 43),

('https://picsum.photos/300/200?random=100', 44),
('https://picsum.photos/300/200?random=101', 44),
('https://picsum.photos/300/200?random=102', 44),

('https://picsum.photos/300/200?random=103', 45),
('https://picsum.photos/300/200?random=104', 45),
('https://picsum.photos/300/200?random=105', 45),

('https://picsum.photos/300/200?random=106', 46),
('https://picsum.photos/300/200?random=107', 46),
('https://picsum.photos/300/200?random=108', 46),

('https://picsum.photos/300/200?random=109', 47),
('https://picsum.photos/300/200?random=110', 47),
('https://picsum.photos/300/200?random=111', 47),

('https://picsum.photos/300/200?random=112', 48),
('https://picsum.photos/300/200?random=113', 48),
('https://picsum.photos/300/200?random=114', 48),

('https://picsum.photos/300/200?random=115', 49),
('https://picsum.photos/300/200?random=116', 49),
('https://picsum.photos/300/200?random=117', 49),

('https://picsum.photos/300/200?random=118', 50),
('https://picsum.photos/300/200?random=119', 50),
('https://picsum.photos/300/200?random=120', 50),

('https://picsum.photos/300/200?random=121', 51),
('https://picsum.photos/300/200?random=122', 51),
('https://picsum.photos/300/200?random=123', 51),

('https://picsum.photos/300/200?random=124', 52),
('https://picsum.photos/300/200?random=125', 52),
('https://picsum.photos/300/200?random=126', 52),

('https://picsum.photos/300/200?random=127', 53),
('https://picsum.photos/300/200?random=128', 53),
('https://picsum.photos/300/200?random=129', 53),

('https://picsum.photos/300/200?random=130', 54),
('https://picsum.photos/300/200?random=131', 54),
('https://picsum.photos/300/200?random=132', 54),

('https://picsum.photos/300/200?random=133', 55),
('https://picsum.photos/300/200?random=134', 55),
('https://picsum.photos/300/200?random=135', 55),

('https://picsum.photos/300/200?random=136', 56),
('https://picsum.photos/300/200?random=137', 56),
('https://picsum.photos/300/200?random=138', 56),

('https://picsum.photos/300/200?random=139', 57),
('https://picsum.photos/300/200?random=140', 57),
('https://picsum.photos/300/200?random=141', 57),

('https://picsum.photos/300/200?random=142', 58),
('https://picsum.photos/300/200?random=143', 58),
('https://picsum.photos/300/200?random=144', 58),

('https://picsum.photos/300/200?random=145', 59),
('https://picsum.photos/300/200?random=146', 59),
('https://picsum.photos/300/200?random=147', 59),

('https://picsum.photos/300/200?random=148', 60),
('https://picsum.photos/300/200?random=149', 60),
('https://picsum.photos/300/200?random=150', 60);

INSERT INTO size_products (name, additional_costs) VALUES
('Small', 0),
('Medium', 5000),
('Large', 10000),
('Extra Large', 15000),
('XXL', 20000);

INSERT INTO variant_products (name, additional_costs) VALUES
('Hot', 0),
('Cold', 0),
('Spicy', 10000),
('Less Sugar', 5000),
('Sweet', 5000);

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
