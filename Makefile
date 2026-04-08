build:
	go build -o claude-code-statusline -trimpath -ldflags="-s -w" .

install: build
	mv -fv ./claude-code-statusline ~/.claude

lint:
	golangci-lint run ./...
