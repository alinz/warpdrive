FROM alpine

RUN apk add --no-cache ca-certificates

# copy all the certificate requires to load the server
RUN mkdir /cert

ADD ./cert/ca-command.crt /cert
ADD ./cert/command.crt /cert
ADD ./cert/command.key /cert

ADD ./cert/ca-query.crt /cert
ADD ./cert/query.crt /cert
ADD ./cert/query.key /cert

# copy the executable to root path
COPY ./bin/server/warpdrive-linux-server /

# running the code
CMD ["/warpdrive-linux-server"]
