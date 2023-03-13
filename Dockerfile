FROM  golang:1.19-alpine as builder


RUN mkdir /app
WORKDIR /app

ADD go.sum .
ADD go.mod .
RUN go mod download -x


ADD . .
# RUN go mod tidy
# RUN go mod download -x
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM --platform=linux/amd64 alpine:latest as runner

WORKDIR  /app

COPY --from=builder /app/talkee /app/talkee

EXPOSE 80