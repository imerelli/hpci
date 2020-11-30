all: bin/hpci
test: lint unit-test

PLATFORM=local 

.PHONY: bin/hpci
bin/hpci:
	@DOCKER_BUILDKIT=1 docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: unit-test
unit-test:
	@DOCKER_BUILDKIT=1 docker build . --target unit-test

.PHONY: unit-test-coverage
unit-test-coverage:
	@DOCKER_BUILDKIT=1 docker build . --target unit-test-coverage \
	--output coverage/
	cat coverage/cover.out

.PHONY: lint
lint:
	@DOCKER_BUILDKIT=1 docker build . --target lint
