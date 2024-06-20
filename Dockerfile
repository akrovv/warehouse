FROM golang:alpine AS builder

RUN mkdir /warehouse
WORKDIR /warehouse

COPY . .

RUN go mod download && go build cmd/warehouse/main.go

FROM alpine

COPY --from=builder /warehouse/main .
COPY --from=builder /warehouse/config.yml .

CMD [ "./main" ]