FROM golang:1.21-alpine as builder

WORKDIR /src/go

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend .

RUN go install ./worker-lamoda


FROM alpine:3.16 AS worker-test

COPY --from=builder /go/bin/worker-lamoda /usr/local/bin/

ENTRYPOINT ["worker-lamoda"]