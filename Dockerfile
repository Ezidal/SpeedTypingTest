FROM golang:1.23.5

WORKDIR /app
EXPOSE 8080

COPY main.go .
COPY ./static/ ./static
COPY ./go.mod .
COPY ./go.sum .
RUN go build -o myapp ./main.go

CMD [ "./myapp" ]