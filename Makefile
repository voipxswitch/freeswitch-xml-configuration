default: gotest

gobuild:
	./dev/build

gotest:
	./dev/test-coverage.sh
