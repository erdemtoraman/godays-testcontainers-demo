.PHONY: clean-network clean, test-parallel test-simple test-e2e watch-container watch-network clean-container
clean:
	go clean -testcache


test-parallel: clean
	make test-simple & make test-e2e

test-simple: clean
	 cd demo1 && go test ./... && cd ..
test-e2e: clean
	 cd demo2/test && go test ./... && cd ../..

clean-network:
	$(shell docker network rm $(shell docker network list -q) > /dev/null 2>&1 )
	exit 0

clean-container:
	$(shell docker rm -f $(shell docker ps -aq))
	exit 0

watch-network:
	watch docker network list

watch-container:
	watch docker ps --format \"{{.Image}} : {{.Ports}}\"

Watch-labels:
	watch docker ps --format \"{{.Image}} : {{.Labels}}\"