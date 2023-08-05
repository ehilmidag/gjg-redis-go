FROM golang:1.19-alpine

WORKDIR /usr/src/app

COPY . .

RUN chmod +x go-api-boilerplate

EXPOSE 8080

CMD ./go-rest-api-boilerplate