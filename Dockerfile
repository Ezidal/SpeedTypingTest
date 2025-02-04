FROM golang:1.23.5

WORKDIR /app
EXPOSE 8080

COPY ./cmd ./cmd
COPY ./templates/ ./templates
COPY ./go.mod .
COPY ./go.sum .
RUN go build -o myapp ./cmd/main.go

CMD [ "./myapp" ]