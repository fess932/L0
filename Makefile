.PHONY: l0 producer

l0:
	go run cmd/l0/main.go

producer:
	go run cmd/producer/main.go