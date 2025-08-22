# Makefile for Mahdaad  Test
GO := go
PKG := ./...

# ------------------------
# Run each challenge
# ------------------------
run-ch1:
	go run cmd/challenge1/main.go

run-ch2:
	go run cmd/challenge2/main.go

run-ch3:
	go run cmd/challenge3/main.go

run-ch4:
	go run cmd/challenge4/main.go


# ------------------------
# Test
# ------------------------
test:
	go test -v $(PKG)
