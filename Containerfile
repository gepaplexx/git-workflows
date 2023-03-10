ARG VERSION=latest

FROM golang:1.20-alpine as BUILD
ARG VERSION

RUN apk add --no-cache gcc g++ make

WORKDIR /app/

COPY src/go.mod src/go.sum ./
RUN go mod verify && go mod download
COPY src/ .

RUN GOOOS=linux GOARCH=amd64 go build -o git-workflows -ldflags="-X main.version=$VERSION" .

FROM alpine:3.15.6

RUN apk add --no-cache ca-certificates curl wget bash git openssh

COPY --from=BUILD /app/git-workflows /bin/git-workflows

RUN chgrp -R 0 /bin/git-workflows && chmod -R g=u /bin/git-workflows

ENTRYPOINT [ "/bin/git-workflows" ]