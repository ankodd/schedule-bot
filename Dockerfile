FROM bitnami/golang:1.23.2-debian-12-r0

WORKDIR /app
COPY . .

ENV STORAGE_PATH=storage/storage.db
ENV GOPROXY=https://proxy.golang.org,direct

EXPOSE 8080

RUN go mod tidy
RUN go mod download

RUN mkdir "storage"
RUN mkdir "build"

RUN go build -o ./build/main ./cmd/schedule-bot/main.go
CMD ["./build/main"]