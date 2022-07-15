# builder
FROM golang:1.18.2-alpine3.16 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY manager.go ./
COPY cli/ cli/
COPY cmd/ cmd/
COPY webhook/ webhook/
COPY handler/ handler/
COPY provider/ provider/
COPY util/ util/
RUN go build ./cmd/cloud-secrets-manager

# runner
FROM alpine:3.16.0 AS runner
WORKDIR /usr/bin/app
RUN addgroup --system app && adduser --system --shell /bin/false --ingroup app app
COPY --from=builder /usr/src/app/cloud-secrets-manager .
RUN chown -R app:app /usr/bin/app
USER app
ENTRYPOINT [ "/usr/bin/app/cloud-secrets-manager" ]
