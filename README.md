# sample-apiserver

基于`apimaster`工程所创建的apiserver.

1. 首先我们需要创建CA证书

``` shell
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt
```

2. 创建给客户端使用的证书


``` shell
openssl req -out client.csr -new -newkey rsa:4096 -nodes -keyout client.key -subj "/CN=development/O=system:masters"
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt
```

3. 创建`curl`需要的证书

``` shell
openssl pkcs12 -export -in ./client.crt -inkey ./client.key -out client.p12 -passout pass:password
```

4. 放置前面生成的证书在合适的位置，然后启动服务

```shell
go run cmd/sample-apiserver/main.go --insecure-bind-address="127.0.0.1" --insecure-port=8080 \
  --secure-port=8443 --tls-private-key-file="~/keys/ca.key" --tls-cert-file="~/keys/ca.crt" \
  --enable-swagger-ui=true  --swagger-ui-file-path="./third_party/swagger-ui"
```

至此，你已经创建了一个`k8s`风格的API Server, 访问`127.0.0.1:8080`