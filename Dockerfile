FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api    ./cmd/api

FROM alpine:3.20

RUN adduser -D -g '' appuser

COPY --from=build /out/api    /usr/local/bin/server
COPY migrations /migrations

USER appuser
EXPOSE 8080

CMD ["server"]
