build:
	go build -o bin/network-edr-agent ./agent
	go build -o bin/network-edr-server ./server

test:
	go test ./...

web:
	cd web && npm ci && npm run build

clean:
	rm -rf bin web/dist
