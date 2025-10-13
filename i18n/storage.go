package i18n

import (
	"errors"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/cryptox"
	bolt "go.etcd.io/bbolt"
)

// TranTx represents a translation cache transaction
type TranTx struct {
	tx *bolt.Tx
}

// Commit commits the transaction
func (p *TranTx) Commit() error {
	err := p.tx.Commit()
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

// Rollback rolls back the transaction
func (p *TranTx) Rollback() error {
	err := p.tx.Rollback()
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

// NeedTran checks if the translation needs to be updated based on source content hash
func (p *TranTx) NeedTran(lang string, key string, value any) (bool, error) {
	bucket, err := p.tx.CreateBucketIfNotExists(candy.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return false, err
	}

	keyBytes := candy.ToBytes(key)
	oldHash := string(bucket.Get(keyBytes))
	if oldHash == "" {
		// First time seeing this key, store the hash
		newHash := cryptox.Md5(candy.ToBytes(value))
		err = bucket.Put(keyBytes, candy.ToBytes(newHash))
		if err != nil {
			log.Errorf("err:%s", err)
		}
		return false, nil
	}

	// Check if source content has changed
	newHash := cryptox.Md5(candy.ToBytes(value))
	return oldHash != newHash, nil
}

// Update updates the translation cache with the new source content hash
func (p *TranTx) Update(lang string, key string, value any) error {
	bucket, err := p.tx.CreateBucketIfNotExists(candy.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	keyBytes := candy.ToBytes(key)
	newHash := cryptox.Md5(candy.ToBytes(value))
	err = bucket.Put(keyBytes, candy.ToBytes(newHash))
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

// TranCache represents a translation cache backed by BoltDB
type TranCache struct {
	db *bolt.DB
}

// NewTranCache creates a new translation cache at the given file path
func NewTranCache(filepath string) (*TranCache, error) {
	if filepath == "" {
		return nil, errors.New("filepath is empty")
	}
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &TranCache{db}, nil
}

// Close syncs and closes the cache database
func (p *TranCache) Close() error {
	_ = p.db.Sync()
	return p.db.Close()
}

// Begin starts a new transaction
func (p *TranCache) Begin() (*TranTx, error) {
	tx, err := p.db.Begin(true)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	return &TranTx{tx: tx}, nil
}
