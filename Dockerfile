FROM golang
WORKDIR /go/src/github.com/suhlig/grafana-backup-restore
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grafana-backup-restore .

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/src/github.com/suhlig/grafana-backup-restore/grafana-backup-restore .
CMD ["./grafana-backup-restore"]
