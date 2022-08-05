package postgres

import (
	"context"
	"restapisrv/srv/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// New Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Новый RestInterface.
// AddProduct добавляет продукт и проверяет, если продукт с таким name уже существует.
func (s *Storage) AddProduct(p storage.ProductItem) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO products (prod_name, prod_tr_name, prod_desc1, prod_desc2, prod_desc3, prod_logo)
       	SELECT $1, $2, $3, $4, $5
       	WHERE NOT EXISTS (SELECT 1 FROM products WHERE prod_name=$1);
	`,
		p.Prod_name, p.Prod_tr_name, p.Prod_desc1, p.Prod_desc2, p.Prod_desc3, p.Prod_logo,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeleteProduct удаляет продукт по id.
func (s *Storage) DeleteProduct(p storage.ProductItem) error {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM products
		WHERE products.prod_id = $1;
	`,
		p.Prod_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeleteAllProducts удаляет все продукты, очищает таблицу.
func (s *Storage) DeleteAllProducts() error {
	rows, err := s.db.Query(context.Background(), `
		TRUNCATE TABLE products RESTART IDENTITY CASCADE;
	`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// Products возвращает продукты, отсортированные по id.
func (s *Storage) Products() ([]storage.ProductItem, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			*
		FROM products
		ORDER BY prod_id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var products []storage.ProductItem
	for rows.Next() {
		var t storage.ProductItem
		err = rows.Scan(
			&t.Prod_id,
			&t.Prod_name,
			&t.Prod_tr_name,
			&t.Prod_desc1,
			&t.Prod_desc2,
			&t.Prod_desc3,
			&t.Prod_logo,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, t)

	}
	return products, rows.Err()
}

// SearchSortedProducts возвращает продукты, отсортированные по паттерну имени.
func (s *Storage) SearchSortedProducts(p storage.ProductItem) ([]storage.ProductItem, error) {
	rows, err := s.db.Query(context.Background(), `
	(SELECT *FROM products
		WHERE lower(prod_name) LIKE '%' || lower($1) || '%'
		ORDER BY position(lower($1) in lower(prod_name)), prod_name)
        UNION ALL
        (SELECT * FROM products
		WHERE lower(prod_tr_name) LIKE '%' || lower($1) || '%'
		ORDER BY position(lower($1) in lower(prod_tr_name)), prod_tr_name);
	`,
		p.Prod_name,
	)
	if err != nil {
		return nil, err
	}
	var products []storage.ProductItem
	for rows.Next() {
		var t storage.ProductItem
		err = rows.Scan(
			&t.Prod_id,
			&t.Prod_name,
			&t.Prod_tr_name,
			&t.Prod_desc1,
			&t.Prod_desc2,
			&t.Prod_desc3,
			&t.Prod_logo,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, t)

	}
	return products, rows.Err()
}

// Магазины.
// AddStore добавляет магазин и проверяет, если магазин с таким name и address уже существует.
func (s *Storage) AddStore(st storage.StoreItem) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO stores (store_name, store_address, store_phone, store_email,
							store_logo, store_latitude, store_longitude)
       	SELECT $1, $2, $3, $4, $5, $6, $7
       	WHERE NOT EXISTS (SELECT 1 FROM stores WHERE store_address=$2 AND store_name=$1);
	`,
		st.Store_name, st.Store_address, st.Store_phone, st.Store_email, st.Store_logo,
		st.Store_latitude, st.Store_longitude,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeleteStore удаляет магазин по id.
func (s *Storage) DeleteStore(st storage.StoreItem) error {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM stores
		WHERE stores.store_id = $1;
	`,
		st.Store_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeleteAllStores удаляет все магазины, очищает таблицу.
func (s *Storage) DeleteAllStores() error {
	rows, err := s.db.Query(context.Background(), `
		TRUNCATE TABLE stores RESTART IDENTITY CASCADE;
	`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// Stores возвращает магазины, отсортированные по id.
func (s *Storage) Stores() ([]storage.StoreItem, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			*
		FROM stores
		ORDER BY store_id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var stores []storage.StoreItem
	for rows.Next() {
		var t storage.StoreItem
		err = rows.Scan(
			&t.Store_id,
			&t.Store_name,
			&t.Store_address,
			&t.Store_phone,
			&t.Store_email,
			&t.Store_logo,
			&t.Store_latitude,
			&t.Store_longitude,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, t)

	}
	return stores, rows.Err()
}

// Цена.
// AddUpdatePrice добавляет или обновляет цену на продукт в магазине.
func (s *Storage) AddUpdatePrice(price storage.PriceItem) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO products_stores (prod_id, store_id, price)
		VALUES ($1,$2,$3)
		ON CONFLICT (prod_id, store_id) DO UPDATE
		SET price = $3;
	`,
		price.Prod_id, price.Store_id, price.Price,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeletePrice удаляет цену на продукт в магазине по id продукта и id магазина.
func (s *Storage) DeletePrice(st storage.PriceItem) error {
	rows, err := s.db.Query(context.Background(), `
		DELETE FROM products_stores
		WHERE store_id = $1 AND prod_id = $2;
	`,
		st.Store_id, st.Prod_id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// DeleteAllPrices удаляет все цены, очищает таблицу.
func (s *Storage) DeleteAllPrices() error {
	rows, err := s.db.Query(context.Background(), `
		TRUNCATE TABLE products_stores RESTART IDENTITY;
	`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	return rows.Err()
}

// PriceList показывает весь список цен и координат магазинов.
func (s *Storage) PriceList() ([]storage.PriceListItem, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT prod_name, prod_logo, store_name, store_address, store_phone, store_email, store_logo, store_latitude, store_longitude, price
	FROM products INNER JOIN products_stores USING(prod_id)
				  INNER JOIN stores USING(store_id)
	ORDER BY 1, 10, 3;
	`,
	)
	if err != nil {
		return nil, err
	}
	var prices []storage.PriceListItem
	for rows.Next() {
		var t storage.PriceListItem
		err = rows.Scan(
			&t.Prod_name,
			&t.Prod_logo,
			&t.Store_name,
			&t.Store_address,
			&t.Store_phone,
			&t.Store_email,
			&t.Store_logo,
			&t.Store_latitude,
			&t.Store_longitude,
			&t.Price,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, t)

	}
	return prices, rows.Err()
}

// ProductPrice получение всех цен по названию прордукта prod_name.
func (s *Storage) ProductPrice(pr storage.PriceListItem) ([]storage.PriceListItem, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT prod_name, prod_logo, store_name, store_address, store_phone, store_email, store_logo, store_latitude, store_longitude, price
		FROM products INNER JOIN products_stores USING(prod_id)
					  INNER JOIN stores USING(store_id)
		WHERE products.prod_name = $1
		ORDER BY 1, 10, 3;
	`,
		pr.Prod_name,
	)
	if err != nil {
		return nil, err
	}
	var prices []storage.PriceListItem
	for rows.Next() {
		var t storage.PriceListItem
		err = rows.Scan(
			&t.Prod_name,
			&t.Prod_logo,
			&t.Store_name,
			&t.Store_address,
			&t.Store_phone,
			&t.Store_email,
			&t.Store_logo,
			&t.Store_latitude,
			&t.Store_longitude,
			&t.Price,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, t)

	}
	return prices, rows.Err()
}
