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

To connect to Nats MQ, you need to specify the configuration. The configuration file contains the following structure in JSON format:
```json
{
  "url": "nats://0.0.0.0:4222",
  "cluster_id": "test-cluster"
}
```
There is an example of a configuration file in the repository.

You can specify the path to the configuration file using the `-config` flag. For expample:
```
./stancli sub some.subject -config ./custom-config.json
```

By default, the path to the configuration file is `./config.json`, so you can **omit** the flag if you use a file named `config.json` located next to the executable of stancli.
