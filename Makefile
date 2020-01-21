.PHONY: clean-network clean, test-parallel test-simple test-e2e watch-container watch-network clean-container
clean:
	go clean -testcache


test-parallel: clean
	(echo 'make test-simple'; echo 'make test-e2e') | parallel -j 2

test-simple: clean
	 cd unit && go test ./... && cd ..
test-e2e: clean
	 cd integration/test && go test ./... && cd ../..

clean-network:
	$(shell docker network rm $(shell docker network list -q) > /dev/null 2>&1 )
	exit 0

clean-container:
	$(shell docker rm -f $(shell docker ps -aq))
	exit 0

watch-network:
	watch docker network list

watch-container:
	watch docker ps