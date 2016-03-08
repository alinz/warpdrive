FROM ubuntu:14.04

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates

COPY ./bin/warpdrive /usr/bin/warpdrive

RUN mkdir /usr/bin/temp

EXPOSE 8221

CMD ["warpdrive", "-config=/etc/warpdrive.conf"]
