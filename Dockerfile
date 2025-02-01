FROM golang:1.23.5

WORKDIR /app
EXPOSE 8080

COPY . .
RUN go build -o myapp .

CMD [ "./myapp" ]