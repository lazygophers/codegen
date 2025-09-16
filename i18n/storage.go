package i18n

import (
	"errors"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/cryptox"
	"github.com/lazygophers/utils/stringx"
	bolt "go.etcd.io/bbolt"
)

type TranTx struct {
	tx *bolt.Tx
}

func (p *TranTx) Commit() error {
	err := p.tx.Commit()
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func (p *TranTx) Rollback() error {
	err := p.tx.Rollback()
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func (p *TranTx) NeedTran(lang string, key string, value any) (bool, error) {
	bucket, err := p.tx.CreateBucketIfNotExists(stringx.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return false, err
	}

	old := stringx.ToString(bucket.Get(stringx.ToBytes(key)))
	if old == "" {
		err = bucket.Put(candy.ToBytes(key), stringx.ToBytes(cryptox.Md5(candy.ToBytes(value))))
		if err != nil {
			log.Errorf("err:%s", err)
		}
		return false, nil
	}

	return old != cryptox.Md5(candy.ToBytes(value)), nil
}

func (p *TranTx) Update(lang string, key string, value any) error {
	bucket, err := p.tx.CreateBucketIfNotExists(stringx.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	err = bucket.Put(candy.ToBytes(key), stringx.ToBytes(cryptox.Md5(candy.ToBytes(value))))
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

type TranCache struct {
	db *bolt.DB
}

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

func (p *TranCache) Close() error {
	_ = p.db.Sync()
	return p.db.Close()
}

func (p *TranCache) Begin() (*TranTx, error) {
	tx, err := p.db.Begin(true)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	return &TranTx{tx: tx}, nil
}
