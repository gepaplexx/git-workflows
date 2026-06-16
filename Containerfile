ARG VERSION=latest

FROM golang:1.22-alpine as BUILD
ARG VERSION

RUN apk add --no-cache gcc g++ make

WORKDIR /app/

COPY src/go.mod src/go.sum ./
RUN go mod verify && go mod download
COPY src/ .

RUN GOOOS=linux GOARCH=amd64 go build -o git-workflows -ldflags="-X main.version=$VERSION" .

FROM alpine:3.24.1

RUN apk add --no-cache ca-certificates curl wget bash git openssh

RUN addgroup -g 1000 -S workflow && \
    adduser -u 1000 -S workflow -G workflow

RUN mkdir -p /workflow/

COPY --from=BUILD /app/git-workflows /bin/git-workflows
COPY src/templates/default-descriptor.json /workflow/default-descriptor.json

RUN chgrp -R 0 /workflow  && \
    chmod -R g=u /workflow/ && \
    chgrp -R 0 /bin/git-workflows &&  \
    chmod -R g=u /bin/git-workflows


USER 1000

ENTRYPOINT [ "/bin/git-workflows" ]