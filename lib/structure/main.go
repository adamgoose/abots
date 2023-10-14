package structure

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/defval/di"
	"github.com/nutsdb/nutsdb"
)

// DB wraps the nutsdb.DB interface to add some helper methods
type DB struct {
	di.Inject

	Nuts *nutsdb.DB
}

// View runs the given function in a nutsdb.Tx
func (d *DB) View(fn func(tx *Tx) error) error {
	return d.Nuts.View(func(tx *nutsdb.Tx) error {
		return fn(&Tx{tx})
	})
}

// Update runs the given function in a nutsdb.Tx
func (d *DB) Update(fn func(tx *Tx) error) error {
	return d.Nuts.Update(func(tx *nutsdb.Tx) error {
		return fn(&Tx{tx})
	})
}

// Tx wraps the nutsdb.Tx interface to add some helper methods
type Tx struct {
	*nutsdb.Tx
}

// PutStruct encodes the given value as a gob and stores it in the given bucket
func (tx *Tx) PutStruct(bucket, key string, value any) error {
	enc := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(enc).Encode(value); err != nil {
		return err
	}

	fmt.Printf("Saving Struct %s\n", key)

	return tx.Put(bucket, []byte(key), enc.Bytes(), 0)
}

// GetStruct decodes the value stored in the given bucket as a gob
func (tx *Tx) GetStruct(bucket, key string, value any) error {
	entry, err := tx.Get(bucket, []byte(key))
	if err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewReader(entry.Value)).Decode(value)
}

func (tx *Tx) PutKVP(bucket string, kvp map[string]string) error {
	for k, v := range kvp {
		if err := tx.Put(bucket, []byte(k), []byte(v), 0); err != nil {
			return err
		}
	}

	return nil
}
