FROM golang:1.16-alpine as builder
ARG VERSION=0.0.0-docker
WORKDIR /app
# Explicitely not updating  as this breaks Docker layer caching. Deps will be updated each time we update the base image.
RUN apk --no-cache add ca-certificates git

# Better dep caching
COPY go.mod .
COPY go.sum .
COPY ./api/v1/model/go.mod api/v1/model/go.mod
COPY ./api/v1/model/go.sum api/v1/model/go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'github.com/riser-platform/riser-server/pkg/util.VersionString=${VERSION}'" -o /go/bin/riser-server

FROM alpine:3.13
RUN apk --no-cache add ca-certificates git openssh-client
# Add git hosts to known hosts file so we can use
# StrickHostKeyChecking with git+ssh
ADD ./docker/known_hosts.sh /known_hosts.sh
RUN sh /known_hosts.sh /etc/ssh/ssh_known_hosts && rm /known_hosts.sh

# SSH key config for optionally ounted key
COPY ./docker/ssh_config /etc/ssh/ssh_config

RUN mkdir -p /etc/riser/ssh

# Riser uses the --author flag with all commits, git still requires something to be set globally
RUN git config --global user.email "riser-server@riser.dev"
RUN git config --global user.name "riser-server"
COPY --from=builder /go/bin/riser-server /riser-server
COPY ./migrations /migrations


EXPOSE 8000
CMD ["/riser-server"]