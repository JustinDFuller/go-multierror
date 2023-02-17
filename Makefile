test-watch:
	@reflex -s -- sh -c 'clear && $(MAKE) test';

test:
	@go test -race ./...;

