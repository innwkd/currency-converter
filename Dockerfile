FROM golang AS build

ADD . /go/src/github.com/yddmat/currency-converter
WORKDIR /go/src/github.com/yddmat/currency-converter

RUN CGO_ENABLED=0 GOOS=linux go build -o converter .

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/yddmat/currency-converter/converter /srv

WORKDIR /srv

EXPOSE 8080
EXPOSE 4444

ENTRYPOINT ./converter