FROM ubuntu:14.04

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates

COPY ./etc/warpdrive.conf /etc/warpdrive.conf
COPY ./bin/warpdrive/warpdrive-linux-amd64 /usr/bin/warpdrive

RUN mkdir /usr/bin/temp

EXPOSE 8221

CMD ["warpdrive", "-config=/etc/warpdrive.conf"]
