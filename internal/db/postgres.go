package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase(DatabaseURl string) (*pgxpool.Pool, error) {

	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	config, err := pgxpool.ParseConfig(DatabaseURl)

	if err != nil {
		log.Printf("Parsing config from Database URL Error : %v", err.Error())
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Println("Database connection Error")
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Println("Pinging not successful!")
		pool.Close()
		return nil, err
	}
	log.Println("Database connected successfullu!")

	return pool, err
}
