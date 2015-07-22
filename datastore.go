package main

import (
	"bytes"
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

func (b *boltBucket) get(k string) (rs io.ReadSeeker, err error) {
	var v []byte
	err = b.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(b.name).Get([]byte(k))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("Value for Key is nil.")
	}

	rs = bytes.NewReader(v)
	return rs, nil
}

func (b *boltBucket) set(k string, r io.Reader) {
	bb := &bytes.Buffer{}
	_, _ = bb.ReadFrom(r)
	b.setBytes(k, bb.Bytes())
}

func (b *boltBucket) setBytes(k string, bs []byte) {
	_ = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.name)
		if err := b.Put([]byte(k), bs); err != nil {
			return err
		}
		return nil
	})
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
