### Description
a Simple central config server backed by [bbolt](https://github.com/etcd-io/bbolt) with support of callbacks when a property is update.

### API

Server has 3 APIs to put, get, and register webhooks:


##### Add/Update Property
```
curl -X 'POST' \
  'localhost:8080/v1/put' \
  -H 'Content-Type: application/json' \
  -d '{
  "namespace": "ns1",
  "key": "key1",
  "value": "val1"
}'
```
Would return `{"message": "property updated"}` in case of success, `400` response for missing parameters, `500` in case of a server error

##### Read Property
```
curl "localhost:8080/v1/get?key=key1&namespace=ns1"
```
would return `{"value":"val1"}` in case `key1` is set in namespace `ns1`, `400` response for missing parameters, `500` in case of a server error

##### Register webhook
```
curl -H "Content-Type: application/json" -d '{"namespace":"ns1","key":"key1","callback":"http://localhost:8080/ping"}' -X POST localhost:8080/v1/webhook
```
would return `{"message":"callback registered"}` when callback registers successfully, `400` for missing parameters, `500` in case of internal server error
callbacks can be registered even before a property is set, after registering it whenever the property value changes the callback url is called

### Running the app

There is a make target for running it locally, simple execute `make run` and it would start it with the default port `8080`.
You can also use the docker image `ghcr.io/shadi/simple_config:latest`, it accepts environment variables `Port` where you can set the port the web server listens to and `GIN_MODE` where you can set it to `RELEASE` or `DEBUG`

more to come
