FROM golang:alpine AS builder

LABEL stage=builder

ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata
RUN apk add --no-cache make

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -ldflags="-s -w" -o /app/rest ./cmd/rest/main.go
RUN go build -ldflags="-s -w" -o /app/storage ./cmd/storage/main.go

FROM alpine as rest

LABEL stage=rest
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York
ENV TZ America/New_York

WORKDIR /app

COPY --from=builder /app/rest /app/rest
COPY --from=builder /build/scripts/wait-for.sh /app/scripts/wait-for.sh
COPY --from=builder /build/config/config.yml /app/config/config.yml

CMD ["./rest"]

FROM alpine as storage

LABEL stage=storage
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York
ENV TZ America/New_York

WORKDIR /app

COPY --from=builder /app/storage /app/storage
COPY --from=builder /build/config/config.yml /app/config/config.yml

CMD ["./storage"]