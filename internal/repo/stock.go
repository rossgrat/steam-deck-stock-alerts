package repo

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("stock_state")

type StockRepo struct {
	db *bolt.DB
}

func NewStockRepo(path string) (*StockRepo, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening bolt database: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("creating bucket: %w", err)
	}

	return &StockRepo{db: db}, nil
}

func (r *StockRepo) GetState(packageID string) (*bool, error) {
	var state *bool

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		v := b.Get([]byte(packageID))
		if v == nil {
			return nil
		}
		val := string(v) == "true"
		state = &val
		return nil
	})

	return state, err
}

func (r *StockRepo) SetState(packageID string, available bool) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		val := "false"
		if available {
			val = "true"
		}
		return b.Put([]byte(packageID), []byte(val))
	})
}

func (r *StockRepo) Close() error {
	return r.db.Close()
}
