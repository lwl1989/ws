# C1000k WebSocket Server

### test run

```
mkdir -p  $GOPATH/src/15x97599p8.imwork.net:965/mars

git clone git@15x97599p8.imwork.net:mars/ws.git

cd ws/main/server

if need change port, sed -i 's/:8080/your port/' ./main.go

go build -o ws_server main.go

./ws_server

cd ../client

go build -o ws_client main.go

can coding to test any more benchmark
```