ALTER TABLE product_images
ADD CONSTRAINT fk_product_images_product FOREIGN KEY(products_id) REFERENCES products(id);
