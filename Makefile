.PHONY: clean-network


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