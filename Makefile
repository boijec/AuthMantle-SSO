APPLICATION = authmantle-sso
# HC_NAME = hc
# HC_ARGS = "health/main.go"

.PHONY: default clean run

default: $(APPLICATION)

clean:
	@-rm -f $(APPLICATION)
#	@-rm -f $(HC_NAME)

run: $(APPLICATION)
	./$(APPLICATION)

test:
	go test -v ./...

bench:
	go test -v ./... -bench=. -benchmem -count 5

# run_hc: $(HC_NAME)
# 	./$(HC_NAME)

$(APPLICATION):
	CGO_ENABLED=0 go build -v -o $(APPLICATION)
#	CGO_ENABLED=0 go build -v -o $(HC_NAME) $(HC_ARGS)