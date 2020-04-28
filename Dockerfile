FROM golang:1.12-alpine as build_base

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/AnthonyNixon/raymond-api

ENV GO111MODULE=on
COPY go.mod .
RUN go mod tidy

FROM build_base as binary_builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /raymond-api cmd/app/main.go

FROM alpine

COPY --from=binary_builder /raymond-api /raymond-api
ENV GIN_MODE=release
ENV PORT=8080
EXPOSE 8080

CMD ["./raymond-api"]