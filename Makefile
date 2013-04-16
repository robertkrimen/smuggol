.PHONY: test install release

test:
	go test -i
	go test -tags author

install: test
	go install

release: test
	godocdown -signature . > README.markdown
