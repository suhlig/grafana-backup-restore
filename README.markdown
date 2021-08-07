# Grafana Backup and Restore

This is more or less a productized version of the [Grafana SDK examples](https://github.com/grafana-tools/sdk/blob/master/cmd/README.md).

## Backup

Saves all dashboards and datasources to files; replicating the folder structure used in Grafana.

## Restore

Re-creates all dashboards and data sources from files; rebuilding the folder structure as found in the file system.

# Development

Boring Go code.

# Build

```command
$ docker build . -t suhlig/grafana-backup-restore
```

# Manual Test

```command
$ ./run
```
