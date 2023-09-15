package data

type Storage interface{
  GetProperty(namespace, key string) (string, error) 
  SetProperty(namespace, key, value string) error;
  ReadNamespaceData(namespace string) (map[string]string, error);
  RegisterCallback(namespace, key, callback string) error;
}
