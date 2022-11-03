package postgres

import (
	"context"
	"restapisrv/storage"

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
       	SELECT $1, $2, $3, $4, $5, $6
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

// SearchSortedProducts возвращает продукты, отсортированные по паттерну имени.
func (s *Storage) SearchSortedProducts(pr storage.SearchItem) ([]storage.SearchItem, error) {
	rows, err := s.db.Query(context.Background(), `
		(SELECT prod_name, MIN(price)
			FROM products
			INNER JOIN products_stores USING(prod_id)
			WHERE lower(prod_name) LIKE '%' || lower($1) || '%'
			GROUP BY prod_name
			ORDER BY position(lower($1) in lower(prod_name)), 1)
        	UNION ALL
		(SELECT prod_name, MIN(price)
			FROM products
			INNER JOIN products_stores USING(prod_id)
			WHERE lower(prod_tr_name) LIKE '%' || lower($1) || '%'
			GROUP BY prod_name, prod_tr_name
			ORDER BY position(lower($1) in lower(prod_tr_name)), 1);
	`,
		pr.Prod_name,
	)
	if err != nil {
		return nil, err
	}
	var products []storage.SearchItem
	for rows.Next() {
		var t storage.SearchItem
		err = rows.Scan(
			&t.Prod_name,
			&t.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, t)

	}
	return products, rows.Err()
}

// SearchSortedProductsWithStore возвращает продукты, отсортированные по паттерну имени мин цену и магазин.
func (s *Storage) SearchSortedProductsWithStore(pr storage.PriceListItem) ([]storage.PriceListItem, error) {
	rows, err := s.db.Query(context.Background(), `
		(WITH y AS (SELECT prod_id, prod_name, MIN(price) AS price
			FROM products
			INNER JOIN products_stores USING(prod_id)
			WHERE lower(prod_name) LIKE '%' || lower($1) || '%'
			GROUP BY prod_id, prod_name
		)
		SELECT products.prod_name, prod_logo, store_name, store_address, store_phone, store_email, store_logo, store_latitude, store_longitude, price
		FROM products
		INNER JOIN y USING(prod_id) INNER JOIN products_stores  USING(price)
		INNER JOIN stores USING(store_id)
		ORDER BY position(lower($1) in lower(products.prod_name)), 1)
		UNION ALL
		(WITH y AS (SELECT prod_id, prod_name, MIN(price) AS price
			FROM products
			INNER JOIN products_stores USING(prod_id)
			WHERE lower(prod_tr_name) LIKE '%' || lower($1) || '%'
			GROUP BY prod_id, prod_name, prod_tr_name
		)
		SELECT products.prod_name, prod_logo, store_name, store_address, store_phone, store_email, store_logo, store_latitude, store_longitude, price
		FROM products
		INNER JOIN y USING(prod_id) INNER JOIN products_stores  USING(price)
		INNER JOIN stores USING(store_id)
		ORDER BY position(lower($1) in lower(prod_tr_name)), 1
		);
	`,
		pr.Prod_name,
	)
	if err != nil {
		return nil, err
	}
	var products []storage.PriceListItem
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
		products = append(products, t)

	}
	return products, rows.Err()
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

// ShortPriceList показывает весь список товаров и минимальных цен.
func (s *Storage) ShortPriceList() ([]storage.SearchItem, error) {
	rows, err := s.db.Query(context.Background(), `
		(SELECT prod_name, MIN(price)
			FROM products
			INNER JOIN products_stores USING(prod_id)
			GROUP BY prod_name
			ORDER BY 1)
	`,
	)
	if err != nil {
		return nil, err
	}
	var products []storage.SearchItem
	for rows.Next() {
		var t storage.SearchItem
		err = rows.Scan(
			&t.Prod_name,
			&t.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, t)

	}
	return products, rows.Err()
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

// Пользователи.
// SignUp добавление нового пользователя.
func (s *Storage) SignUp(user storage.Credentials) error {
	rows, err := s.db.Query(context.Background(), `
		INSERT INTO users (useremail, password)
		VALUES ($1,$2)
	`,
		user.Useremail, user.Password,
	)
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}

// SignIn вход пользователя.
func (s *Storage) SignIn(user storage.CredentialsShort) (storage.Credentials, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT * from users
		WHERE users.useremail = $1;
	`,
		user.Useremail,
	)

	var u storage.Credentials
	if err != nil {
		u.Useremail = user.Useremail
		return u, err
	}
	rowNum := 0

	for rows.Next() {
		err = rows.Scan(
			&u.Useremail,
			&u.Password,
			&u.Userstatus,
			&u.Confirmstring,
			&u.Usernickname,
			&u.Lastlogindate,
			&u.Lastlogindevice,
		)
		if err != nil {
			u.Useremail = user.Useremail
			return u, err
		}
		rowNum++

	}

	if rowNum < 1 {
		u.Useremail = user.Useremail
		return u, storage.ErrUserNotFound
	}
	rows.Close()
	return u, rows.Err()
}

// UserExist существует ли пользователь, вернуть почту и статус.
func (s *Storage) UserExistEmailStatus(user storage.CredentialsUserEmailStatus) error {
	rows, err := s.db.Query(context.Background(), `
		SELECT useremail, userstatus from users
		WHERE useremail = $1;
	`,
		user.Useremail,
	)
	if err != nil {
		return err
	}
	rowNum := 0
	var u storage.CredentialsUserEmailStatus
	for rows.Next() {
		err = rows.Scan(
			&u.Useremail,
			&u.Userstatus,
		)
		if err != nil {
			u.Useremail = user.Useremail
			return err
		}
		rowNum++

	}
	if rowNum < 1 {
		u.Useremail = user.Useremail
		return storage.ErrUserNotFound
	}

	rows.Close()
	return rows.Err()
}

// SetConfirmString добавление строки подтверждения почты.
func (s *Storage) SetConfirmString(user storage.CredentialsConfirm) error {
	rows, err := s.db.Query(context.Background(), `
		UPDATE users
		SET confirmstring = $2 
		WHERE useremail = $1;
	`,
		user.Useremail, user.Confirmstring,
	)
	if err != nil {
		return err
	}
	rows.Close()
	return rows.Err()
}

// ConfirmStringAndStatus получение строки подтверждения почты и статуса.
func (s *Storage) ConfirmStringAndStatus(user storage.CredentialsConfirm) (storage.CredentialsConfirm, error) {
	rows1, err := s.db.Query(context.Background(), `
		SELECT useremail, userstatus, confirmstring from users
		WHERE users.confirmstring = $1;
	`,
		user.Confirmstring,
	)

	var u storage.CredentialsConfirm
	if err != nil {
		u.Useremail = user.Useremail
		return u, err
	}
	rowNum := 0

	for rows1.Next() {
		err = rows1.Scan(
			&u.Useremail,
			&u.Userstatus,
			&u.Confirmstring,
		)
		if err != nil {
			u.Useremail = user.Useremail
			return u, err
		}
		rowNum++

	}

	if u.Userstatus > 0 {
		return u, storage.ErrAlreadyConfirmed
	}

	if rowNum < 1 {
		u.Useremail = user.Useremail
		return u, storage.ErrConfirmStringNotFound
	}

	rows2, err := s.db.Query(context.Background(), `
		UPDATE users
		SET userstatus = 1 
		WHERE confirmstring = $1;
	`,
		user.Confirmstring,
	)
	if err != nil {
		return u, err
	}

	rows2.Close()

	return u, rows2.Err()
}

// SetUserStatus изменение статуса пользователя по useremail.
func (s *Storage) SetUserStatus(user storage.CredentialsConfirm) error {
	rows, err := s.db.Query(context.Background(), `
		UPDATE users
		SET userstatus = $2 
		WHERE confirmstring = $1;
	`,
		user.Useremail, user.Userstatus,
	)
	rows.Close()
	if err != nil {
		return err
	}
	if rows.CommandTag().RowsAffected() < 1 {
		return storage.ErrUserNotFound
	}
	return rows.Err()
}
