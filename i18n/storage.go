package i18n

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
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

func (p *TranTx) Check(lang string, key string, value any) (bool, error) {
	bucket, err := p.tx.CreateBucketIfNotExists(stringx.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return false, err
	}

	return stringx.ToString(bucket.Get(stringx.ToBytes(key))) == cryptox.Md5(anyx.ToBytes(value)), nil
}

func (p *TranTx) Update(lang string, key string, value any) error {
	bucket, err := p.tx.CreateBucketIfNotExists(stringx.ToBytes(lang))
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	err = bucket.Put(anyx.ToBytes(key), stringx.ToBytes(cryptox.Md5(anyx.ToBytes(value))))
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
		filepath = state.Config.I18n.AutoTran.RecordPath
	}
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &TranCache{db}, nil
}

func (p *TranCache) Close() error {
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
