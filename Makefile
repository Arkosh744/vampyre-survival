BINARY_NAME   = vampyre
CERT          = certs/cert.pem
KEY           = certs/key.pem
SIGN_NAME     = Vampyre Survival
CERT_SUBJ     = /CN=Vampyre Survival/O=arktrix games
CERT_DAYS     = 3650

.PHONY: build build-linux build-windows winres sign gen-certs clean test

# Build and sign Windows exe (default)
build: build-windows sign
	@echo "Done: $(BINARY_NAME)-signed.exe"

winres:
	cd cmd/game && $(shell go env GOPATH)/bin/go-winres make

build-windows: winres
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o $(BINARY_NAME).exe ./cmd/game/

build-linux:
	go build -o $(BINARY_NAME) ./cmd/game/

sign:
	osslsigncode sign \
		-certs $(CERT) \
		-key $(KEY) \
		-n "$(SIGN_NAME)" \
		-h sha256 \
		-in $(BINARY_NAME).exe \
		-out $(BINARY_NAME)-signed.exe
	@rm -f $(BINARY_NAME).exe
	@mv $(BINARY_NAME)-signed.exe $(BINARY_NAME).exe
	@echo "Signed: $(BINARY_NAME).exe"

gen-certs:
	@mkdir -p certs
	openssl req -x509 -newkey rsa:4096 \
		-keyout $(KEY) \
		-out $(CERT) \
		-days $(CERT_DAYS) \
		-nodes \
		-subj "$(CERT_SUBJ)" \
		-addext "extendedKeyUsage=codeSigning" \
		-addext "keyUsage=digitalSignature"
	@echo "Generated: $(CERT), $(KEY)"

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe $(BINARY_NAME)-signed.exe

test:
	go test ./...
