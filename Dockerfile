FROM golang:1.8

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/oscarpfernandez/idgen

RUN go install github.com/oscarpfernandez/idgen

# Run the outyet command by default when the container starts.
CMD ["/go/bin/idgen","server"]

EXPOSE 8080