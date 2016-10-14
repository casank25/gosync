# gosync

Gosync copies high volume S3 objects from one S3 bucket to another one withing the same region. It uses concurrency to complete the task.
It can read file containing JSON objects (separated by `\n`) or list objects from one bucket and copy them over to another.


## Prerequisites
- Go 1.5
- aws cli configured with your access key id, secret access key and region

## Instalation
```
> git clone https://github.com/casank25/gosync
> cd gosync
> go build
```

## Usage

```
> ./gosync --source souce-bucket-name --destination destination-bucket-name
```

## Options

- `--source` **Required.** Source Bucket
- `--destination` **Required.** Destination Bucket
- `--region` AWS region. Defaults to 'us-west-2'
- `--retries` How many times should aws retry to copy file before it becomes an error. Defaults to 5
- `--start` Specify a key where to start copying from
- `--workers` How many workers should work concurrently at any given time. Defaults to 100
- `--queue` Max queue size. Defaults to 1000
- `--reader` Where to read object keys from. It can be `file` or `s3`. Defaults to `s3`
- `--file-path` If you chose `--reader` as file, this one is required
