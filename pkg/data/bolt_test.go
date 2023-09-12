package data

import (
	"testing"
)


func TestAdd(t *testing.T) {
  storage,err := NewStorage("./test1.tdb")

  if err != nil {
    t.Error("Failed to get storage")
  }

  err = storage.SetProperty("namespace", "key", "value", "")
  if err != nil {
    t.Error("Failed to set")
  }

  val, err := storage.GetProperty("namespace", "key")
  if err != nil {
    t.Error("Failed to get property")
  }

  if val != "value"{
    t.Error("Value mismatch")
  }
}

func TestBucketDoesNotExist(t *testing.T) {
  s, err := NewStorage("./test2.tdb")
  if err != nil {
    t.Error("Failed to get storage")
  }

  val, err := s.GetProperty("no_namespace","anything")
  if err != nil {
    t.Error("Error getting property")
  } 

  if val != "" {
    t.Error("Value mismatch")
  }
}

func TestPropertyDoesNotExist(t *testing.T) {
  s, err := NewStorage("./test3.tdb")
  if err != nil {
    t.Error("Failed to get storage")
  }

  err = s.SetProperty("namespace","prop1","val1","")
  if err != nil {
    t.Error("Error getting property")
  } 

  val, err := s.GetProperty("namespace", "prop2")

  if val != "" {
    t.Error("Value mismatch")
  }
}
