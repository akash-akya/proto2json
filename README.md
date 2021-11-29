## Proto2Json

Simple utility to convert binary protobuf to json. Useful to pipe the
output to `jq` for further filtering, processing

### Consume binary protobuf form stdin

```sh
./proto2json -root_path ./sample -message sample.Person < ./sample/data.bin
```

### Consume output from [kt](https://github.com/fgeller/kt)

Value must be base64 encoded.

```sh
kt consume -topic sample-topic -encodevalue base64 | ./proto2json -kt_stream -root_path ./sample -message sample.Person
```

Will replace `value` with decoded json and write to stdout
