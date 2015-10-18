FROM golang:1.5.1
WORKDIR /src
# Application config vars
ENV PORT=8080
# Install application dependencies
RUN apt-get update && \
    apt-get -y install sqlite3
# Use Godeps so we don't have to
# reinstall all of the deps on the container
ENV GOPATH=/src/Godeps/_workspace
# Move data into container
ADD . .
# Create database and compile
RUN go build && \
    make fresh
# Run binary file
CMD ["./src"]
