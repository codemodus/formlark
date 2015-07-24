package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

// DataStores holds available data stores.
type dataStores struct {
	dcbsRsrcs *boltDB
	dcbAsts   *boltBucket
	dcbMrks   *boltBucket
}

type boltDB struct {
	*bolt.DB
}

type boltBucket struct {
	name []byte
	db   *boltDB
}

func getDataCacheLocal(file string) (d *boltDB, err error) {
	dir := "/var/lib/" + path.Base(os.Args[0]) + "/"
	if err = os.MkdirAll(dir, 0775); err != nil {
		return nil, err
	}
	dt, err := bolt.Open(dir+file, 0660, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	d = &boltDB{dt}
	return d, nil
}

func (d *boltDB) getBucket(bucket string) (b *boltBucket, err error) {
	name := []byte(bucket)
	err = d.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	b = &boltBucket{name: name, db: d}
	return b, nil
}

func (b *boltBucket) rebuild() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(b.name)
	})
	if err != nil {
		return err
	}
	err = b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *boltBucket) get(k string) (io.ReadSeeker, error) {
	v, err := b.getBytes(k)
	if err != nil {
		return nil, err
	}
	rs := bytes.NewReader(v)
	return rs, nil
}

func (b *boltBucket) getBytes(k string) ([]byte, error) {
	var v []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(b.name).Get([]byte(k))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (b *boltBucket) set(k string, r io.Reader) {
	bb := &bytes.Buffer{}
	_, _ = bb.ReadFrom(r)
	b.setBytes(k, bb.Bytes())
}

func (b *boltBucket) setBytes(k string, bs []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.name)
		if err := b.Put([]byte(k), bs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *boltBucket) find(k string) error {
	var v []byte
	_ = b.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(b.name).Get([]byte(k))
		return nil
	})
	if v == nil {
		return errors.New("Value for Key is nil.")
	}
	return nil
}

type bolter interface {
	getID() (string, error)
	get() error
	set() error
}

type boltItem struct {
	DS *dataStores `json:"-"`
	ID string      `json:"-"`
}

func (bi *boltItem) getID() (string, error) {
	if bi.ID != "" {
		return bi.ID, nil
	}
	return "", errors.New("no id")
}

func (bi *boltItem) get() error {
	id, err := bi.getID()
	if err != nil {
		return err
	}
	v, err := bi.DS.dcbAsts.getBytes(id)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, bi); err != nil {
		return err
	}
	return nil
}

func (bi *boltItem) set() error {
	v, err := json.Marshal(bi)
	if err != nil {
		return err
	}
	id, err := bi.getID()
	if err != nil {
		return err
	}
	if err = bi.DS.dcbAsts.setBytes(id, v); err != nil {
		return err
	}
	return nil
}
