FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN mkdir logs/

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o be-b-impact

ARG GO_ENV=PRODUCTION

EXPOSE 80

ENTRYPOINT ["/app/be-b-impact"]