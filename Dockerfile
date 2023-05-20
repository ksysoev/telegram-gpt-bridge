############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/gptbot/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get
# Build the binary.
RUN go build -o /go/bin/gptbot main.go
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/gptbot /go/bin/gptbot
# Run the hello binary.
ENTRYPOINT ["/go/bin/gptbot"]