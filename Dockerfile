FROM golang:1.23.4-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /websockets-go main.go

FROM scratch

COPY --from=builder /websockets-go /websockets-go

EXPOSE 8082

CMD ["/websockets-go"]

