FROM golang:latest

COPY ./ ./
RUN go build -o game .
CMD ["./game"]