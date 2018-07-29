NAMESPACE  := logicmonitor
REPOSITORY := releasemanager
VERSION    := 0.1.0-alpha.0

all:
	docker build --rm --build-arg VERSION=$(VERSION) --build-arg CI=$(CI) -t $(NAMESPACE)/$(REPOSITORY):latest .
	docker tag $(NAMESPACE)/$(REPOSITORY):latest $(NAMESPACE)/$(REPOSITORY):$(VERSION)
