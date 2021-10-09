module github.com/HirokiHanada11/go-restful

go 1.17

require (
	github.com/HirokiHanada11/go-restful/websockets v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)

require github.com/gorilla/websocket v1.4.2 // indirect

replace github.com/HirokiHanada11/go-restful/websockets => /app/websockets
