FROM golang:1.9 as build
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-release-manager
COPY ./ ./
ARG VERSION
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /releasemanager -ldflags "-X \"github.com/logicmonitor/k8s-release-manager/pkg/constants.Version=${VERSION}\"" cmd/releasemanager/main.go

FROM golang:1.9 as test
ARG CI
ENV CI=$CI
WORKDIR $GOPATH/src/github.com/logicmonitor/k8s-release-manager
RUN go get -u github.com/alecthomas/gometalinter
RUN gometalinter --install
COPY --from=build $GOPATH/src/github.com/logicmonitor/k8s-release-manager ./
RUN chmod +x ./scripts/test.sh; sync; ./scripts/test.sh
RUN cp coverage.txt /coverage.txt

FROM alpine:3.6
LABEL maintainer="Jeff Wozniak <jeff.wozniak@logicmonitor.com>"
RUN apk --update add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apk/*
WORKDIR /app
COPY --from=build /releasemanager /bin
COPY --from=test /coverage.txt /coverage.txt

ENTRYPOINT ["releasemanager"]
CMD ["watch"]
