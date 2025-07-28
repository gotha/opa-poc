FROM golang:1.19 AS builder

WORKDIR /app

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server main.go

FROM scratch

COPY --from=builder /app/server /server

EXPOSE 8080
CMD ["/server"]

