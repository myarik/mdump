FROM golang:alpine as build

ARG APP_VERSION

WORKDIR /build
# Copy source
COPY ./ ./

# Build the Go app
RUN go build -o mdump -ldflags "-X main.AppVersion=${APP_VERSION} -s -w" ./main.go

# (stage 2) package the binary into an Alpine container
FROM alpine:3.7 as app

WORKDIR /srv
COPY --from=build /build/mdump /usr/bin/mdump

RUN \
    apk add --no-cache --update git bash curl postgresql-client && \
    rm -rf /var/cache/apk/*

CMD ["/usr/bin/mdump"]