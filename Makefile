build:
	go build -ldflags "-s -w" -o dist ./...

certs:
	mkdir -p certs
	mkcert -cert-file certs/public.crt -key-file certs/private.key localhost
