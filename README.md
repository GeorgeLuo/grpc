### Setup development environment
Build dev environment
```
sudo docker build -t grpc --file Dockerfile.dev .
```
Generate key and cert for mutual TLS. These default locations will place the key/cert local to where the command is ran. Changing the location will require configuration in the server and client to reference correctly.
```
chmod +x generate-keys.sh
```
```
./generate-keys.sh
```
Initialize dev environment
```
sudo docker run --volume "$(pwd)":/go/src/github.com/GeorgeLuo/grpc --interactive --tty --publish 8443:8443 --net=host grpc
```
Run go server, exposed to port 8443
```
go run main.go handlers.go execUtil.go syncOutput.go syncMap.go
```


use client.go as cli-like process to start

```
go run client/client.go start -cert cert.pem -key key.pem -command ./test_process.sh -host localhost
```
to stop
```
go run client/client.go stop -cert cert.pem -key key.pem -task_id "987f769fca40-3635" -host localhost
```
to get status
```
go run client/client.go status -cert cert.pem -key key.pem -task_id "987f769fca40-3635" -host localhost
```
