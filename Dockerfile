FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app

# RUN apk update \
#     && apk --no-cache --update add build-base git

COPY ./bbiot-api/go.mod ./bbiot-api/go.sum ./

RUN go mod download && go mod tidy

COPY ./bbiot-api ./

RUN go build -o main ./cmd

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
CMD [ "/app/main" ]