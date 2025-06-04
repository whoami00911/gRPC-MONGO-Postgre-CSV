FROM golang:1.24-alpine

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go version  \ 
    && go build -o webApp ./cmd/main.go

CMD [ "./webApp" ]