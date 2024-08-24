FROM golang:1.22.5

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o backend .

EXPOSE 3000

CMD ["./backend"]