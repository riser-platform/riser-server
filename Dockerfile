FROM golang:1.14-alpine as builder
WORKDIR /app
RUN apk add --update --no-cache ca-certificates git

# Better dep caching (other option is to vendor deps)
COPY go.mod .
COPY go.sum .

# Each module must be copied before running `go mod download`.
# When this repo is public we can just use the github reference in go.mod and remove this.
COPY ./api/v1/model/go.mod api/v1/model/go.mod
COPY ./api/v1/model/go.sum api/v1/model/go.sum

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/riser-server

FROM alpine
RUN apk update
RUN apk add git
# Riser uses the --author flag with all commits, git still requires something to be set globally
RUN git config --global user.email "riser-server@riser.dev"
RUN git config --global user.name "riser-server"
COPY --from=builder /go/bin/riser-server /riser-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./migrations /migrations
EXPOSE 8000
CMD ["/riser-server"]