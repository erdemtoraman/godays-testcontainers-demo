.PHONY: clean-network


test-simple:
	 cd unit && go test ./... -count=500 && cd ..
test-e2e:
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