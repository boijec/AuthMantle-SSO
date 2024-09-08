FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
RUN apk add --no-cache make gcc musl-dev linux-headers && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates
COPY . .
RUN make
FROM scratch AS runner
COPY --from=builder /app/authmantle-sso .
COPY --from=builder /app/templates ./templates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV HOST="0.0.0.0"
CMD ["./authmantle-sso"]