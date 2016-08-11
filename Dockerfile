FROM alpine:3.1

# need to mount your $GOPATH/src
VOLUME  /gopath/src
ENV GOPATH /gopath

ENTRYPOINT ["/bin/wag"]

ADD ./bin/wag /bin/wag
