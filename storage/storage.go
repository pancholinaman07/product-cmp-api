package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"product-cmp-api/types"
)

type Storage interface {
	CreateProduct(*types.Product) error
	DeleteProduct(string, string) error
	GetProduct(string, string) (*types.Product, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=3007 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createProductTable()
}

func (s *PostgresStore) createProductTable() error {
	query := `create table if not exists product (
		id serial primary key,
		brand varchar(100),
		model varchar(100),
    	prod_name varchar(100),
    	color varchar(100),
    	prod_size varchar(100),
    	processor varchar(100),
    	frontMp varchar(100),
    	backMp varchar(100),
    	ram varchar(100),
		storage varchar(100),
    	madeIn varchar(100),
    	ImageLoc varchar(100)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateProduct(prod *types.Product) error {
	query := `insert into product
	(brand, model, prod_name, color, prod_size, processor, frontMp, backMp, ram, storage, madeIn, ImageLoc)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := s.db.Query(
		query,
		prod.Brand,
		prod.Model,
		prod.Name,
		prod.Color,
		prod.Size,
		prod.Processor,
		prod.FrontMp,
		prod.BackMp,
		prod.Ram,
		prod.Storage,
		prod.MadeIn,
		prod.ImageLoc)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteProduct(brand, model string) error {
	// todo correct query
	_, err := s.db.Query("delete from product where brand=$1 and model=$2", brand, model)
	return err
}

func (s *PostgresStore) GetProduct(brand, model string) (*types.Product, error) {
	//todo correct query
	rows, err := s.db.Query("select * from product where brand=$1 and model=$2", brand, model)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, fmt.Errorf("brand [%s] model [%s] not found", brand, model)
}

func scanIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Brand,
		&product.Model,
		&product.Name,
		&product.Color,
		&product.Size,
		&product.Processor,
		&product.FrontMp,
		&product.BackMp,
		&product.Ram,
		&product.Storage,
		&product.MadeIn,
		&product.ImageLoc)

	return product, err
}
