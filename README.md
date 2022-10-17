# stancli

## About
Command line client for [STAN](https://docs.nats.io/legacy/stan/intro) messaging.

## Usage

### Publishing message

Run following command in terminal, but with required STAN subject instead of `some.subject`

```
./stancli pub some.subject
```

after that type or paste your message body. Finally, to publish message hit `Enter`/`Return` and then `Ctrl-D`.

### Subscribing subject

```
./stancli sub some.subject
```

## Configuration

To connect to Nats MQ, you need to specify configuration. Configuration file contains following structure in JSON format:
```json
{
  "url": "nats://0.0.0.0:4222",
  "cluster_id": "test-cluster"
}
```
There is an example of a configuration file in repository.

You can specify path to configuration file using `-config` flag. For expample:
```
./stancli sub some.subject -config ./custom-config.json
```

By default, path to configuration file is `./config.json`, so you can **omit** flag if you use a file named `config.json` located next to executable of stancli.
