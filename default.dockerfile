FROM golang:1.23.2-bullseye AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X 'github.com/thanksduck/alias-api/cfconfig.allowedDomains=domain1.com,domain2.com' -X 'github.com/thanksduck/alias-api/cfconfig.configJSON=$(cat config.default.json | tr -d '\n')'" -o main .

FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/main .
ENV TZ=Asia/Kolkata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
EXPOSE 6777
CMD ["./main"]
