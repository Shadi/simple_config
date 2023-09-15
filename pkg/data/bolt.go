package data

import (
	"time"

	"github.com/rs/zerolog/log"
	bolt "go.etcd.io/bbolt"
)

const (
  CallbacksBucket = "callbacks";
  LastUpdatedBucket = "timestamps";
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
    _, err := tx.CreateBucketIfNotExists([]byte(CallbacksBucket))
    if err != nil {
      log.Err(err).Msg("Failed to create callbacks bucket")
      return err
    }
    _, err = tx.CreateBucketIfNotExists([]byte(LastUpdatedBucket))
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

func (d *boltDb) SetProperty(namespace, key, value string) error{
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
    tb := tx.Bucket([]byte(LastUpdatedBucket))
    err = tb.Put(id, []byte(time.Now().Format(time.RFC3339)))
    if err != nil {
      log.Err(err).Msg("Failed to set property last updated")
      return err
    }

    return nil
  }) 
}

func (d *boltDb) RegisterCallback(namespace, key, callback string) error {
 return d.db.Update(func(tx *bolt.Tx) error {
    cb := tx.Bucket([]byte(CallbacksBucket))
    err := cb.Put([]byte(namespace+key), []byte(callback))
    if err != nil {
      log.Err(err).Msg("Failed to set callback")
    }
    return err
  })
}

func (d *boltDb) ReadNamespaceData(namespace string) (map[string]string, error) {
  m := make(map[string]string)
  err := d.db.View(func(tx *bolt.Tx) error {
    tx.Bucket([]byte(namespace)).ForEach(func(k, v []byte) error {
      key := make([]byte, len(k))
      value := make([]byte, len(v))
      copy(key, k)
      copy(value, v)
      m[string(key)] = string(value)
      return nil
    })
    return nil
  })
  return m, err
}
