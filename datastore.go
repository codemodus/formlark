package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

// DataStores holds available data stores.
type dataStores struct {
	dcbsRsrcs   *boltDB
	dcbUsers    *boltBucket
	dcbIndUsers *boltBucket
	dcbIndCnfrm *boltBucket
	dcbPosts    *boltBucket
}

type boltDB struct {
	*bolt.DB
}

type boltBucket struct {
	name []byte
	*bolt.Bucket
	db *boltDB
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

func (d *boltDB) getBucket(bucket string) (bb *boltBucket, err error) {
	var b *bolt.Bucket
	name := []byte(bucket)
	err = d.Update(func(tx *bolt.Tx) error {
		b, err = tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	bb = &boltBucket{name: name, Bucket: b, db: d}
	return bb, nil
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

func (b *boltBucket) getManyBytes(start, count int) (map[string][]byte, error) {
	m := make(map[string][]byte)
	ct := 0
	err := b.db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket(b.name).ForEach(func(k, v []byte) error {
			if ct >= start && count > 0 {
				m[string(k)] = v
				count--
			}
			ct++
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return m, err
	}
	return m, nil
}

func (b *boltBucket) set(k string, r io.Reader) error {
	bb := &bytes.Buffer{}
	if _, err := bb.ReadFrom(r); err != nil {
		return err
	}
	if err := b.setBytes(k, bb.Bytes()); err != nil {
		return err
	}
	return nil
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
	err := b.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(b.name).Get([]byte(k))
		return nil
	})
	if err != nil {
		return err
	}
	if v == nil {
		return errors.New("Value for Key is nil.")
	}
	return nil
}

type bolter interface {
	get() (bool, error)
	set() error
}

type boltItem struct {
	BB  *boltBucket `json:"-"`
	BBI *boltBucket `json:"-"`
	ID  string      `json:"-"`
}

func (bi *boltItem) get() (bool, error) {
	if bi.ID == "" {
		return false, errors.New("no id")
	}
	b, err := bi.BB.getBytes(bi.ID)
	if len(b) == 0 || err != nil {
		return false, err
	}
	br := bytes.NewReader(b)
	dec := gob.NewDecoder(br)
	if err := dec.Decode(bi); err != nil {
		return false, err
	}
	return true, nil
}

func (bi *boltItem) set() error {
	if bi.ID == "" {
		return errors.New("no id")
	}

	bb := &bytes.Buffer{}
	enc := gob.NewEncoder(bb)
	if err := enc.Encode(bi); err != nil {
		return err
	}
	if err := bi.BB.set(bi.ID, bb); err != nil {
		return err
	}
	return nil
}
