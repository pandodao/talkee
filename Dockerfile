FROM  golang:1.19-alpine as builder
ARG GIT_TAG
ENV APP_VERSION=$GIT_TAG

RUN mkdir /app
WORKDIR /app

ADD go.sum .
ADD go.mod .
RUN go mod download -x


ADD . .
# RUN go mod tidy
# RUN go mod download -x
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$APP_VERSION"

FROM --platform=linux/amd64 alpine:latest as runner

WORKDIR  /app

COPY --from=builder /app/talkee /app/talkee

EXPOSE 80