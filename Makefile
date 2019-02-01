NAMESPACE  := logicmonitor
REPOSITORY := releasemanager
VERSION    := 0.1.0-alpha.1

all:
	docker build --rm --build-arg VERSION=$(VERSION) --build-arg CI=$(CI) -t $(NAMESPACE)/$(REPOSITORY):latest .
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(REPOSITORY):latest

linux: | local
darwin: | local
local:
ifneq ($(MAKECMDGOALS), darwin)
ifneq ($(MAKECMDGOALS), linux)
	$(error Valid local build targets are "linux" and "darwin")
endif
endif
	GOOS=$(MAKECMDGOALS) GOARCH=amd64 CGO_ENABLED=0 go build -o ./$(REPOSITORY) -ldflags "-X \"github.com/logicmonitor/k8s-release-manager/pkg/constants.Version=${VERSION}\"" cmd/releasemanager/main.go
