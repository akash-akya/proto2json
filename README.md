## Proto2Json

Simple utility to convert protobuf message to json. Useful for streaming protobuf messages to `jq` for further processing. Works with [kt](https://github.com/fgeller/kt) utility.

```
Usage of ./proto2json:
  -proto_path string
        root path for proto files. Path must contain all proto files including needed for decoding the message
  -type string
        complete proto type including package name
  -kt_stream
        consume from stream produced by kt tool and decode. 'value' must be encoded as base64 using "-encodevalue base64" option
  -skip_unpopulated
        when set skips unpopulated fields such as null, empty-string, empty-list
```

#### Convert single protobuf message

```sh
proto2json -proto_path ./sample -type sample.Person < ./sample/data.bin
```

#### Streaming output from [kt](https://github.com/fgeller/kt)

```sh
kt consume -broker 127.0.0.1:9091 -topic sample-topic -encodevalue base64 \
  | proto2json -kt_stream -proto_path ./sample -type sample.Person
```

It will replace `value` field with decoded json and stream it to stdout


### Docker

Docker image contains kt, proto2json, jq and other few other basic essentials.

`kt_complete.sh` and `Dockerfile` is based on https://github.com/Paxa/kt

```sh
docker run akashh/kt kt topic -brokers kafka-broker:9091
```
