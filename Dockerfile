FROM golang
WORKDIR /app

COPY . .

RUN go install /app
ENTRYPOINT /go/bin/user-api
EXPOSE 9090 9092
