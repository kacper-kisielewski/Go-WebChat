FROM golang:1.14-alpine

WORKDIR /app

RUN apk add git gcc libc-dev

RUN go get -u -v github.com/cosmtrek/air@v1.27.0

COPY . .

EXPOSE 8000

CMD ["air", "-d"]
