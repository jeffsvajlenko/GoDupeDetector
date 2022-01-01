FROM alpine:latest
WORKDIR /
ADD bin /bin/
CMD ["/bin/goduped"]
