# Binary name
BINARY=webhook
GOBUILD=go build -o ${BINARY}
GOCLEAN=go clean
VERSION=0.0.2

# Build
build:
	$(GOCLEAN)
	$(GOBUILD)


clean:
	$(GOCLEAN)
	$(RMTARGZ)


release:
	# Clean
	$(GOCLEAN)
	# Build for mac
	#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD)
	# Build for arm
#	$(GOCLEAN)
	#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD)
	# Build for linux
#	$(GOCLEAN)

	GOOS=linux GOARCH=amd64 $(GOBUILD)
	# Build for win
#	$(GOCLEAN)
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD).exe
#	$(GOCLEAN)
docker:
	GOOS=linux GOARCH=amd64 $(GOBUILD)
	docker build -t foo/webhook . --tag registry.cn-shanghai.aliyuncs.com/zzf2001/foo-webhook
	docker push registry.cn-shanghai.aliyuncs.com/zzf2001/foo-webhook