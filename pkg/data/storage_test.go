package data

import (
	"testing"
)


func TestAdd(t *testing.T) {
  storage,err := NewStorage("./test1.tdb")

  if err != nil {
    t.Error("Failed to get storage")
  }

  err = storage.SetProperty("namespace", "key", "value")
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

  err = s.SetProperty("namespace","prop1","val1")
  if err != nil {
    t.Error("Error getting property")
  } 

  val, err := s.GetProperty("namespace", "prop2")

  if val != "" {
    t.Error("Value mismatch")
  }
}

func TestReadBucket(t *testing.T) {
  s, err := NewStorage("./test_storage.tdb")
  if err != nil {
    t.Error("Failed to get storage")
  }

  errors := make([]error, 3)
  errors[0] = s.SetProperty("n1", "p1", "v1")
  errors[1] = s.SetProperty("n1", "p2", "v2")
  errors[2] = s.SetProperty("n1", "p3", "v3")
  if errors[0] != nil || errors[1] != nil || errors[2] != nil {
    t.Errorf("Failed to set properties: %v\n", errors)
  }

  vals, errs := s.ReadNamespaceData("n1")
  if errs != nil {
    t.Error(err)
  }

  if vals["p1"] != "v1" || vals["p2"] != "v2" || vals["p3"] != "v3" {
    t.Errorf("Got %s, %s, %s, want v1, v2, v3 \n", vals["p1"], vals["p2"], vals["p3"])
  }

}
