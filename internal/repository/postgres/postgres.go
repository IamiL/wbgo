package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"time"
	"wbnats/internal/services/order/models"
)

type Storage struct {
	db    *pgxpool.Pool
	cache *cache.Cache
}

func New(host string,
	port string,
	dbName string,
	user string,
	pass string) (*Storage, error) {
	const op = "repository.postgres.New"
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS item(
		chrt_id BIGINT,
	    order_uid VARCHAR(200),
		track_number TEXT,
		price BIGINT,
		rid TEXT,
		name TEXT,
		sale INTEGER,
		size TEXT,
		total_price BIGINT,
		nm_id BIGINT,
		brand TEXT,
		status BIGINT
	                               );

	CREATE TABLE IF NOT EXISTS delivery(
		order_uid VARCHAR(200) PRIMARY KEY,
		name TEXT,
		phone VARCHAR(20),
		zip TEXT,
		city TEXT,
		adress TEXT,
		region TEXT,
		email VARCHAR(330)
	                                   );

	CREATE TABLE IF NOT EXISTS payment(
		order_uid VARCHAR(200) PRIMARY KEY,
		transaction VARCHAR(200),
		request_id TEXT,
		currency VARCHAR(20),
		provider VARCHAR(100),
		amount BIGINT,
		payment_dt BIGINT,
		bank VARCHAR(100),
		delivery_cost BIGINT,
		goods_total BIGINT,
		custom_fee BIGINT
	                                  );

	CREATE TABLE IF NOT EXISTS orders(
		order_uid VARCHAR(200) PRIMARY KEY,
		track_number VARCHAR(200),
		entry VARCHAR(200),
		locale VARCHAR(30),
		internal_signature TEXT,
		customer_id TEXT,
		delivery_service TEXT,
		shardkey VARCHAR(30),
		sm_id BIGINT,
		date_created timestamp,
		oof_shard VARCHAR(30)
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ch := cache.New(5*time.Minute, 10*time.Minute)

	return &Storage{db: pool, cache: ch}, nil
}

func (s *Storage) RestoreCache() {
	orders, err := s.GetOrders(context.Background())
	if err != nil {
		return
	}
	for _, ord := range orders {
		s.cache.Set(ord.UID, &ord, cache.NoExpiration)
	}
}

func (s *Storage) GetOrders(ctx context.Context) ([]models.Order, error) {
	query := `SELECT orders.order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, name, phone, zip, city, adress, region, email  
				FROM orders
    			JOIN payment ON orders.order_uid = payment.order_uid
    			JOIN delivery ON orders.order_uid = delivery.order_uid`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query orders: %w", err)
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.UID, &order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		items, err := s.GetItems(ctx, order.UID)
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) GetItems(ctx context.Context, orderUID string) ([]models.Item, error) {
	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM item WHERE order_uid = $1`

	rows, err := s.db.Query(ctx, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("unable to query orders: %w", err)
	}
	defer rows.Close()

	items := []models.Item{}
	for rows.Next() {
		item := models.Item{}
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (s *Storage) SaveOrder(order *models.Order) (err error) {

	batch := &pgx.Batch{}
	orderQuery := `INSERT INTO orders (
                   order_uid,
                   track_number,
                   entry,
                   locale,
                   internal_signature,
                   customer_id,
                   delivery_service,
                   shardkey,
                   sm_id,
                   date_created,
                   oof_shard
                   ) VALUES (
							@orderUID,
							@trackNumber,
							@entry,
							@locale,
							@internalSignature,
							@customerID,
							@deliveryService,
							@shardkey,
							@smID,
							@dateCreated,
							@oofShard)`
	orderArgs := pgx.NamedArgs{
		"orderUID":          order.UID,
		"trackNumber":       order.TrackNumber,
		"entry":             order.Entry,
		"locale":            order.Locale,
		"internalSignature": order.InternalSignature,
		"customerID":        order.CustomerID,
		"deliveryService":   order.DeliveryService,
		"shardkey":          order.Shardkey,
		"smID":              order.SmID,
		"dateCreated":       order.DateCreated,
		"oofShard":          order.OofShard,
	}
	batch.Queue(orderQuery, orderArgs)

	deliveryQuery := `INSERT INTO delivery (
						order_uid,
						name,
						phone,
						zip,
						city,
						adress,
						region,
						email
						 ) VALUES (
								@orderUID,
							   @name,
							   @phone,
							   @zip,
							   @city,
							   @adress,
							   @region,
							   @email
							   )`
	deliveryArgs := pgx.NamedArgs{
		"orderUID": order.UID,
		"name":     order.Delivery.Name,
		"phone":    order.Delivery.Phone,
		"zip":      order.Delivery.Zip,
		"city":     order.Delivery.City,
		"adress":   order.Delivery.Address,
		"region":   order.Delivery.Region,
		"email":    order.Delivery.Email,
	}
	batch.Queue(deliveryQuery, deliveryArgs)

	paymentQuery := `INSERT INTO payment (
	               order_uid,
		transaction,
		request_id,
		currency,
		provider,
		amount,
		payment_dt,
		bank,
		delivery_cost,
		goods_total,
		custom_fee
						 ) VALUES (
						           @orderUID,
								@transaction,
							   @requestID,
							   @currency,
							   @provider,
							   @amount,
							   @paymentDT,
						    	@bank,
							   @deliveryCost,
							   @goodsTotal,
								@customFee
							   )`
	paymentArgs := pgx.NamedArgs{
		"orderUID":     order.UID,
		"transaction":  order.Payment.Transaction,
		"requestID":    order.Payment.RequestID,
		"currency":     order.Payment.Currency,
		"provider":     order.Payment.Provider,
		"amount":       order.Payment.Amount,
		"paymentDT":    order.Payment.PaymentDT,
		"bank":         order.Payment.Bank,
		"deliveryCost": order.Payment.DeliveryCost,
		"goodsTotal":   order.Payment.GoodsTotal,
		"customFee":    order.Payment.CustomFee,
	}

	batch.Queue(paymentQuery, paymentArgs)

	itemQuery := `INSERT INTO item (
		chrt_id,
	   order_uid,
		track_number,
		price,
		rid,
		name,
		sale,
		size,
		total_price,
		nm_id,
		brand,
		status
						 ) VALUES (
							@chrtID,
							@orderUID,
							@trackNumber,
							@price,
							@rID,
							@name,
							@sale,
							@size,
							@totalPrice,
							@nmID,
							@brand,
							@status
							   )`
	for _, item := range order.Items {
		itemArgs := pgx.NamedArgs{
			"chrtID":      item.ChrtID,
			"orderUID":    order.UID,
			"trackNumber": item.TrackNumber,
			"price":       item.Price,
			"rID":         item.RID,
			"name":        item.Name,
			"sale":        item.Sale,
			"size":        item.Size,
			"totalPrice":  item.TotalPrice,
			"nmID":        item.NmID,
			"brand":       item.Brand,
			"status":      item.Status,
		}
		batch.Queue(itemQuery, itemArgs)
	}
	results := s.db.SendBatch(context.Background(), batch)

	if results.Close() != nil {
		return err
	}
	s.cache.Set(order.UID, order, cache.NoExpiration)
	return nil
}

func (s *Storage) Order(ctx context.Context, uid string) (models.Order, error) {

	if x, found := s.cache.Get(uid); found {
		ord := x.(*models.Order)
		return *ord, nil
	}
	return models.Order{}, fmt.Errorf("no data")
}
