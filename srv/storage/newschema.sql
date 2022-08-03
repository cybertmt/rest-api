DROP TABLE IF EXISTS products_stores;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS stores;

CREATE TABLE products
(
    prod_id    SERIAL PRIMARY KEY,
    prod_name  TEXT   NOT NULL,
    prod_desc1 TEXT,
    prod_desc2 TEXT,
    prod_desc3 TEXT,
    prod_logo  TEXT
);

CREATE TABLE stores
(
    store_id       SERIAL PRIMARY KEY,
    store_name     TEXT   NOT NULL,
    store_address  TEXT   NOT NULL,
    store_phone    TEXT   NOT NULL,
    store_email    TEXT   NOT NULL,
    store_logo     TEXT,
    store_latitude  FLOAT NOT NULL,
    store_longitude FLOAT NOT NULL
);

CREATE TABLE products_stores
(
    prod_id    INT,
    store_id   INT,
    price NUMERIC CHECK (price >= 0),
    FOREIGN KEY (prod_id) REFERENCES products (prod_id) ON DELETE CASCADE,
    FOREIGN KEY (store_id) REFERENCES stores (store_id) ON DELETE CASCADE
);

INSERT INTO products (prod_name, prod_desc1) VALUES ('Асперин', 'Асперин: параметры');
INSERT INTO products (prod_name, prod_desc1) VALUES ('Панадол', 'Панадол: параметры');
INSERT INTO products (prod_name, prod_desc1) VALUES ('Парацетамол', 'Парацетамол: параметры');

INSERT INTO stores (store_name, store_address, store_email, store_phone, store_latitude, store_longitude) 
VALUES ('Ригла', 'Гончарный пр., 6, стр. 1, Москва', 'info@rigla.ru', '8 (800) 777-03-03', 55.739399, 37.649848);
INSERT INTO stores (store_name, store_address, store_email, store_phone, store_latitude, store_longitude) 
VALUES ('Здоров.ру', 'ул. Шаболовка, 34, стр. 3, Москва', 'info@zdorov.ru', '+7 (495) 363-35-00', 55.718311, 37.607876);
INSERT INTO stores (store_name, store_address, store_email, store_phone, store_latitude, store_longitude) 
VALUES ('Горздрав', 'Большая Переяславская ул., 11, Москва', 'info@gorzdrav.ru', '+7 (499) 653-62-77', 55.784470, 37.641093);

INSERT INTO products_stores (prod_id, store_id, price)
VALUES (1, 1, 10.50);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (1, 2, 11.30);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (1, 3, 10.00);

INSERT INTO products_stores (prod_id, store_id, price)
VALUES (2, 1, 120.50);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (2, 2, 140.30);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (2, 3, 135.00);

INSERT INTO products_stores (prod_id, store_id, price)
VALUES (3, 1, 40.50);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (3, 2, 45.30);
INSERT INTO products_stores (prod_id, store_id, price)
VALUES (3, 3, 39.00);
