package data

import (
	"time"

	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
)

const (
  callbacksBucket = "callbacks";
  lastUpdatedBucket = "timestamps";
)

type boltDb struct{
  db *bolt.DB;
}

func NewStorage(path string) (Storage, error) {
  bdb, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 60 * time.Second})
  if err != nil {
    log.Err(err).Msg("Failed to open DB")
    return nil, err
  }
  err = bdb.Update(func (tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists([]byte(callbacksBucket))
    if err != nil {
      log.Err(err).Msg("Failed to create callbacks bucket")
      return err
    }
    _, err = tx.CreateBucketIfNotExists([]byte(lastUpdatedBucket))
    return err
  })

  log.Info().Msg("Database Opened")
  return &boltDb{db: bdb}, err
}

func (d *boltDb) GetProperty(namespace, key string) (string, error) {
  var val []byte
  err := d.db.View(func(tx *bolt.Tx) error{
    b := tx.Bucket([]byte(namespace))
    if b == nil {
      return nil
    }
    v := b.Get([]byte(key))
    val = make([]byte, len(v))
    copy(val,v)
    return nil
  })
  return string(val),err
}
func (d *boltDb) SetProperty(namespace, key, value, callback string) error{
 return d.db.Update(func(tx *bolt.Tx) error {
    b,err := tx.CreateBucketIfNotExists([]byte(namespace))
    if err != nil {
      log.Err(err).Msg("Failed to find namespace")
      return err
    }
    err = b.Put([]byte(key), []byte(value))
    if err != nil {
      log.Err(err).Msg("Property set failed")
      return err
    }

    id := []byte(namespace+key)
    cb := tx.Bucket([]byte(callbacksBucket))
    err = cb.Put(id, []byte(callback))
    if err != nil {
      log.Err(err).Msg("Failed to set callback")
      return err
    }
    tb := tx.Bucket([]byte(lastUpdatedBucket))
    err = tb.Put(id, []byte(time.Now().Format(time.RFC3339)))
    return err
  }) 
}
