FROM golang:1-alpine AS builder
WORKDIR /src
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o main.exe ./cmd/

FROM alpine:3
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Bangkok
COPY --from=builder /src/configs/ /go-starter-project/configs/
COPY --from=builder /src/main.exe /go-starter-project/main
WORKDIR /go-starter-project
CMD ["./main"]
EXPOSE 8080
