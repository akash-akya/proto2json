## Proto2Json

Simple utility to convert protobuf message to json. Useful for
streaming protobuf messages to `jq` for further processing. Works with [kt](https://github.com/fgeller/kt) utility.

```
Usage of proto2json:
  -kt_stream
        consume from stream produced by kt tool and decode. 'value' must be encoded as base64 using "-encodevalue base64" option
  -root_path string
        root path for proto files (default "."). Path must contain all proto files including needed for decoding the message
  -message string
        Full proto message name including package name
```

#### Convert single protobuf message

```sh
proto2json -root_path ./sample -message sample.Person < ./sample/data.bin
```

#### Streaming output from [kt](https://github.com/fgeller/kt)

```sh
kt consume -broker 127.0.0.1:9091 -topic sample-topic -encodevalue base64 | proto2json -kt_stream -root_path ./sample -message sample.Person
```

Will replace `value` with decoded json and stream the json to stdout


### Docker

Docker image contains kt, proto2json, jq and other few other basic essentials.

`kt_complete.sh` and `Dockerfile` is based on https://github.com/Paxa/kt

```sh
docker run akashh/kt kt topic -brokers kafka-broker:9091
```
