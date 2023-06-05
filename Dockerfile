FROM golang:1.19.1 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -o ochacafe main.go && chmod +x ./ochacafe

FROM alpine:3.18
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /workspace/ ./

EXPOSE 8080
CMD ["./ochacafe"]