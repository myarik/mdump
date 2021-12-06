FROM golang:alpine as build

ARG APP_VERSION

WORKDIR /build
# Copy source
COPY ./ ./

# Build the Go app
RUN go build -o mdump -ldflags "-X main.AppVersion=${APP_VERSION} -s -w" ./main.go

# (stage 2) package the binary into an Alpine container
FROM postgres:13-alpine

WORKDIR /srv
COPY --from=build /build/mdump /usr/bin/mdump

RUN \
    apk add --no-cache --update git bash curl && \
    rm -rf /var/cache/apk/*

CMD ["/usr/bin/mdump"]