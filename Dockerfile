FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY .env.production .env
# RUN go build -v -o /usr/local/bin/app ./...
RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN go run scripts/generate/main.go && swag init && go build -tags=production -o /usr/local/bin/main .

CMD ["main"]
