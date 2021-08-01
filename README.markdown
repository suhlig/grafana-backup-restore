# Grafana Backup and Restore

This is more or less a productized version of the [Grafana SDK examples](https://github.com/grafana-tools/sdk/blob/master/cmd/README.md).

## Backup

Saves all dashboards and datasources to files; replicating the folder structure used in Grafana.

## Restore

Re-creates all dashboards and datasources from files; rebuilding the folder structure as found in the filesystem.

# Development

Boring Go code.

# Build

```command
$ docker build . -t suhlig/grafana-backup-restore
```

# Acceptance Test

```command
# Launch new Grafana
$ docker run -d -p 3000:3000 grafana/grafana

# Generate API key
$ export GRAFANA_API_TOKEN=$(curl -X POST -H "Content-Type: application/json" -d '{"name":"Acceptance Test", "role": "Admin"}' http://admin:admin@localhost:3000/api/auth/keys | jq --raw-output .key)

# Run
$ go run .  --verbose --url http://localhost:3000 backup dashboards
```
