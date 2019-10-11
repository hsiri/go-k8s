FROM golang:latest

# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /src

# We copy everything in the root directory
# into our /app directory
ADD . /src

# We specify that we now wish to execute
# any further commands inside our /app
# directory
WORKDIR /src

# Run the dependency
RUN go get -d github.com/gorilla/mux
RUN go get -d github.com/sirupsen/logrus
RUN go get -d github.com/frozentech/go-tools/filesystem
RUN go get -d github.com/frozentech/go-tools/secure

# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Our start command which kicks off
# our newly created binary executable
CMD ["/src/main"]
