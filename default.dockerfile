FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Accept the base64 encoded config JSON as a build argument from environment
ARG CONFIG_JSON_STRING

# Build the application with the decoded config JSON
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X 'github.com/thanksduck/alias-api/cfconfig.allowedDomains=1as.in,20032003.xyz,n3y.in,00004000.xyz,90123456.xyz' -X 'github.com/thanksduck/alias-api/cfconfig.configJSON=$(echo ${CONFIG_JSON_STRING} | base64 -d)'" \
    -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

RUN apk update && apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Kolkata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 6777
CMD ["./main"]