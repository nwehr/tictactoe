run:
	go run cmd/tictactoe/**.go

run-server:
	go run cmd/server/**.go

tictactoe:
	go build -o tictactoe cmd/tictactoe/**.go
