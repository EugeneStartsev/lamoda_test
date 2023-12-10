FROM golang:1.21-alpine as builder

WORKDIR /src/go

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend .

RUN go install ./


FROM alpine:3.16 AS worker-test

COPY --from=builder /go/bin/lamoda_test /usr/local/bin/

ENTRYPOINT ["lamoda_test"]