-- Active: 1762853226689@@127.0.0.1@5432@dbcoffeshop
INSERT INTO users (fullname, email, password, role, profile_id, created_at, updated_at)
VALUES
('Ari Eka', 'ari@gmail.com', '123', 'admin', 1, NOW(), NOW()),
('Budi Santoso', 'budi@gmail.com', '123', 'user', 2, NOW(), NOW()),
('Citra Lestari', 'citra@gmail.com', '123', 'user', 3, NOW(), NOW()),
('Dewi Ayu', 'dewi@gmail.com', '123', 'user', 4, NOW(), NOW()),
('Eko Prasetyo', 'eko@gmail.com', '123', 'user', 5, NOW(), NOW());

INSERT INTO profile (pic, phone, address, created_at, updated_at)
VALUES
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', '081234567890', 'Jakarta Selatan', NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', '082345678901', 'Bandung', NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', '083456789012', 'Surabaya', NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', '084567890123', 'Yogyakarta', NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', '085678901234', 'Medan', NOW(), NOW());

INSERT INTO products (name, price, description, stock, isFlashSale, isFavorite_product, category_product_id)
VALUES
('Americano', 25000, 'Kopi hitam klasik dengan aroma kuat.', 100, false, true, 1),
('Cappuccino', 30000, 'Perpaduan espresso, susu, dan foam lembut.', 80, true, true, 1),
('Caffe Latte', 32000, 'Espresso dan susu creamy, cocok untuk santai.', 90, false, true, 1),
('Espresso', 22000, 'Kopi pekat untuk pecinta rasa kuat.', 70, true, false, 1),
('Mocha Latte', 35000, 'Kopi dengan campuran cokelat lezat.', 60, false, true, 1),

('Green Tea Latte', 28000, 'Minuman teh hijau lembut dan manis.', 100, false, false, 2),
('Thai Tea', 26000, 'Teh khas Thailand dengan rasa manis kental.', 120, true, true, 2),
('Lemon Tea', 20000, 'Teh dengan perasan lemon segar.', 150, false, false, 2),
('Black Tea', 18000, 'Teh hitam klasik dengan aroma tajam.', 90, false, false, 2),
('Jasmine Tea', 22000, 'Teh melati harum yang menenangkan.', 80, false, true, 2),

('Fresh Milk', 18000, 'Susu murni segar tanpa tambahan rasa.', 100, false, false, 3),
('Chocolate Milk', 24000, 'Susu cokelat manis favorit anak-anak.', 90, true, true, 3),
('Strawberry Milk', 24000, 'Susu rasa stroberi segar dan lembut.', 80, false, true, 3),
('Banana Milk', 25000, 'Susu dengan aroma pisang alami.', 100, false, false, 3),
('Matcha Milk', 27000, 'Susu dengan bubuk matcha premium.', 100, false, true, 3),

('Orange Juice', 23000, 'Jus jeruk segar penuh vitamin C.', 150, false, false, 4),
('Apple Juice', 24000, 'Jus apel murni tanpa gula tambahan.', 140, false, false, 4),
('Mango Juice', 25000, 'Jus mangga manis segar alami.', 130, true, true, 4),
('Avocado Juice', 28000, 'Jus alpukat kental dengan susu cokelat.', 90, false, true, 4),
('Watermelon Juice', 22000, 'Jus semangka segar, cocok untuk panas.', 150, false, false, 4),

('Donut Original', 15000, 'Donat lembut dengan taburan gula halus.', 80, false, false, 5),
('Chocolate Donut', 17000, 'Donat dengan lelehan cokelat di atasnya.', 100, true, true, 5),
('Strawberry Donut', 17000, 'Donat dengan glaze stroberi manis.', 90, false, false, 5),
('Cheese Croissant', 22000, 'Croissant lembut isi keju meleleh.', 70, true, true, 5),
('Butter Croissant', 20000, 'Croissant klasik dengan mentega berkualitas.', 100, false, false, 5),

('Caramel Macchiato', 34000, 'Kopi dengan sirup karamel dan susu.', 85, true, true, 1),
('Vanilla Latte', 33000, 'Latte lembut dengan aroma vanilla.', 80, false, true, 1),
('Hazelnut Coffee', 36000, 'Kopi dengan rasa kacang hazelnut.', 70, false, false, 1),
('Double Espresso', 26000, 'Dua shot espresso untuk tenaga ekstra.', 60, false, true, 1),
('Irish Coffee', 38000, 'Kopi khas dengan sentuhan krim irish.', 50, false, true, 1),

('Lychee Tea', 23000, 'Teh dingin rasa leci menyegarkan.', 130, true, true, 2),
('Peach Tea', 24000, 'Teh rasa buah persik lembut.', 120, false, false, 2),
('Milk Tea', 25000, 'Teh dengan susu creamy khas Taiwan.', 140, true, true, 2),
('Honey Lemon Tea', 26000, 'Teh lemon dengan madu alami.', 100, false, true, 2),
('Earl Grey Tea', 22000, 'Teh hitam premium dengan aroma bergamot.', 90, false, false, 2),

('Vanilla Milk', 25000, 'Susu rasa vanilla lembut dan manis.', 110, false, false, 3),
('Honey Milk', 26000, 'Susu dengan madu alami.', 120, false, true, 3),
('Almond Milk', 27000, 'Susu almond sehat tanpa laktosa.', 100, false, false, 3),
('Coffee Milk', 28000, 'Perpaduan kopi dan susu manis.', 80, true, true, 3),
('Mint Milk', 24000, 'Susu dengan aroma mint segar.', 90, false, false, 3),

('Pineapple Juice', 25000, 'Jus nanas manis dan segar.', 130, false, false, 4),
('Guava Juice', 26000, 'Jus jambu biji merah penuh vitamin.', 140, false, true, 4),
('Kiwi Juice', 27000, 'Jus kiwi segar rasa asam manis.', 120, true, true, 4),
('Papaya Juice', 23000, 'Jus pepaya lembut dan manis alami.', 150, false, false, 4),
('Coconut Juice', 28000, 'Air kelapa muda alami menyegarkan.', 160, false, true, 4),

('Muffin Chocolate', 18000, 'Muffin lembut rasa cokelat.', 100, false, false, 5),
('Blueberry Muffin', 19000, 'Muffin dengan potongan blueberry asli.', 110, false, true, 5),
('Cinnamon Roll', 20000, 'Roti gulung manis dengan kayu manis.', 90, false, true, 5),
('Brownies', 21000, 'Kue cokelat padat dan manis.', 80, true, true, 5),
('Cheese Cake', 30000, 'Kue lembut dengan lapisan keju.', 70, false, true, 5);

INSERT INTO category_product (name)
VALUES
('Coffee'),
('Tea'),
('Milk'),
('Juice'),
('Snack');

INSERT INTO image_products (image, product_id, created_at, updated_at)
VALUES
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', 51, NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', 53, NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', 54, NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', 52, NOW(), NOW()),
('https://images.pexels.com/photos/33984951/pexels-photo-33984951.jpeg', 53, NOW(), NOW());

INSERT INTO diskons (percentage, name, start_date, end_date, isActive, created_at, updated_at)
VALUES
(10.00, 'Morning Sale', '2025-11-01', '2025-11-30', true, NOW(), NOW()),
(15.00, 'Weekend Special', '2025-11-05', '2025-11-20', true, NOW(), NOW()),
(5.00, 'Member Discount', '2025-11-01', '2025-12-01', true, NOW(), NOW()),
(20.00, 'Flash Deal', '2025-11-10', '2025-11-15', true, NOW(), NOW()),
(25.00, 'Black Friday', '2025-11-25', '2025-11-30', true, NOW(), NOW());

INSERT INTO product_diskon (product_id, diskon_id, created_at, updated_at)
VALUES
(55, 1, NOW(), NOW()),
(52, 2, NOW(), NOW()),
(57, 3, NOW(), NOW()),
(54, 4, NOW(), NOW()),
(56, 5, NOW(), NOW());

INSERT INTO ratings (user_id, product_id, rating, created_at, updated_at)
VALUES
(1, 51, 5, NOW(), NOW()),
(2, 56, 4, NOW(), NOW()),
(3, 53, 5, NOW(), NOW()),
(4, 57, 3, NOW(), NOW()),
(5, 55, 4, NOW(), NOW());

INSERT INTO orders (user_id, payment_method, status, total, created_at, updated_at)
VALUES
(1, 'Credit Card', 'completed', 120000, NOW(), NOW()),
(2, 'Cash', 'pending', 85000, NOW(), NOW()),
(3, 'Debit', 'completed', 95000, NOW(), NOW()),
(4, 'E-Wallet', 'cancelled', 60000, NOW(), NOW()),
(5, 'Credit Card', 'completed', 150000, NOW(), NOW());

INSERT INTO order_items (product_id, quantity, subtotal, order_id, created_at, updated_at)
VALUES
(51, 2, 40000, 6, NOW(), NOW()),
(56, 1, 25000, 7, NOW(), NOW()),
(53, 3, 60000, 8, NOW(), NOW()),
(57, 1, 30000, 9, NOW(), NOW()),
(55, 4, 80000, 10, NOW(), NOW());

INSERT INTO deliverys (order_id, type, fee, created_at, updated_at)
VALUES
(6, 'Regular', 10000, NOW(), NOW()),
(7, 'Express', 15000, NOW(), NOW()),
(8, 'Pickup', 0, NOW(), NOW()),
(9, 'Regular', 10000, NOW(), NOW()),
(10, 'Express', 15000, NOW(), NOW());

INSERT INTO taxs (order_id, name, tax, created_at, updated_at)
VALUES
(6, 'PPN', 2000, NOW(), NOW()),
(7, 'PPN', 2500, NOW(), NOW()),
(8, 'PPN', 3000, NOW(), NOW()),
(9, 'PPN', 1500, NOW(), NOW()),
(10, 'PPN', 3500, NOW(), NOW());

INSERT INTO size_product (name, product_id, additional_costs, created_at, updated_at)
VALUES
  ('Small', 51, 10000, NOW(), NOW()),
  ('Medium', 51, 50000, NOW(), NOW()),
  ('Large', 51, 12000, NOW(), NOW());

INSERT INTO variant (name, product_id, additional_costs, created_at, updated_at)
VALUES
  ('Hot', 51, 32100, NOW(), NOW()),
  ('Iced', 51,300, NOW(), NOW()),
  ('Hot', 52, 3100, NOW(), NOW()),
  ('Cold', 52, 2100, NOW(), NOW()),
  ('Sweet', 53, 3500, NOW(), NOW());
