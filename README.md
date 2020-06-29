# resque_exporter

This will expose metrics from resque in a prometheus format. It gathers metrics by reading from the resque backend, redis.

## configuration

All the configuration is done via environment variables. The following env vars are available:
|           variable name          |               purpose           | default value |
|----------------------------------|---------------------------------|---------------|
| `RESQUE_EXPORTER_REDIS_HOST`     | hostname of ip of redis         | localhost     |
| `RESQUE_EXPORTER_REDIS_PORT`     | port where redis is listening   | 6379          |
| `RESQUE_EXPORTER_REDIS_PASSWORD` | password to connect to redis    | ""            |
| `RESQUE_EXPORTER_REDIS_NAMESPACE`| namespace(prefix) of redis keys | resque        |
| `RESQUE_EXPORTER_REDIS_DB`       | database to read                | 0             |

## exposed metrics

The following metrics are exposed:
|           metric name                  |  type   |               help                  |
|----------------------------------------|---------|-------------------------------------|
| `resque_workers`                       | gauge   | Number of workers                   |
| `resque_workers_working`               | gauge   | Number of workers currently working |
| `resque_queue_jobs{queue="queuename"}` | gauge   | Number of jobs in queue             |
| `resque_jobs_processed_total`          | counter | Total number of processed jobs      |
| `resque_jobs_failed_total`             | counter | Total number of failed jobs         |
| `resque_failed_queue`                  | gauge   | Number of jobs in the failed queue  |
