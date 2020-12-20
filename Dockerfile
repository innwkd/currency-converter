FROM golang AS build

ADD . /app/currency-converter
WORKDIR /app/currency-converter

RUN CGO_ENABLED=0 GOOS=linux go build -o converter .



FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=build /app/currency-converter/converter /srv
COPY --from=build /app/currency-converter/.env /srv

WORKDIR /srv

EXPOSE 12345
EXPOSE 4444

ENTRYPOINT ./converter