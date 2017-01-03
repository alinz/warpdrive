FROM ubuntu:14.04

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates

COPY ./bin/warpdrive/warpdrive-linux-amd64 /bin/warpdrive

EXPOSE 8221

CMD ["/bin/warpdrive", "-config=/warpdrive/conf/warpdrive.conf"]
