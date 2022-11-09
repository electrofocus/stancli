# stancli

## About
Command line client for [STAN](https://docs.nats.io/legacy/stan/intro) messaging.

## Download
Donwload CLI executable for your OS from [latest release assets](https://github.com/electrofocus/stancli/releases).

## Usage
### Publishing message with `--pub`
Flag `--pub` is used to publish message. Run following command in terminal, but with required STAN subject instead of `here.is.some.subject`
```
./stancli --pub --subject=here.is.some.subject
```

after that type or paste your message body. Finally, to publish message hit `Enter`/`Return` and then `Ctrl-D` to interrupt typing.

### Subscribing subject with `--sub`
```
./stancli --sub --subject=here.is.some.subject
```

## Configuration
To connect to specific Nats MQ, you need to specify configuration.

You can specify connection config using `--url` and `--cluster-id` flags. For expample:
```
./stancli
    --sub
    --subject=here.is.some.subject
    --url=nats://our-nats.dev.cloud:1234
    --cluster-id=our-cluster
```
If you omit these flags, default values are used. Default values are `nats://0.0.0.0:4222` for `--url` and `test-cluster` for `cluster-id`.
