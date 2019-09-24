# sudo docker build -t grpc .
# sudo docker run --volume /keys:/keys --interactive --tty --publish 8443:8443 grpc

FROM golang

ADD . /keys

# Install deps
RUN go get github.com/GeorgeLuo/grpc

WORKDIR /keys

# Set default run command
ENTRYPOINT grpc

# Expose port 8443.
EXPOSE 8443
