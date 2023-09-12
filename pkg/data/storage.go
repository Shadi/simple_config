package data

type Property struct {
  Key string;
  Value string;
  Namespace string;
  Callback string;
}

type Storage interface{
  GetProperty(namespace, key string) (string, error) 
  SetProperty(namespace, key, value, callback string) error;
}
