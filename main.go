package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const version = "0.1.0"
const usage = "Usage: proto2json -proto_path <path> -type <proto-typename>"

var protoRootPath = flag.String("proto_path", "", "root path for proto files. Path must contain all proto files including needed for decoding the message")
var protoType = flag.String("type", "", "complete proto type including package name")
var ktStream = flag.Bool("kt_stream", false, "consume from stream produced by kt tool and decode. 'value' must be encoded as base64 using \"-encodevalue base64\" option")
var skipUnpopulated = flag.Bool("skip_unpopulated", false, "whether to emit unpopulated fields such as null, empty-string, empty-list")

func main() {
	flag.Parse()

	if protoType == nil || *protoType == "" {
		printUsageDie()
	}

	if protoRootPath == nil || *protoRootPath == "" {
		printUsageDie()
	}

	marshalOpt := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: !(*skipUnpopulated),
	}

	if *ktStream {
		consumeKTStream(marshalOpt, *protoRootPath, *protoType)
	} else {
		consumeStdin(marshalOpt, *protoRootPath, *protoType)
	}
}

func printUsageDie() {
	flag.Usage()
	os.Exit(-1)
}

func listProtoFiles(root string, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func createProtoRegistry(protoRootPath string, protoFiles []string) (*protoregistry.Files, error) {
	tmpFile := "tmp.pb"

	args := []string{"--include_imports", "--descriptor_set_out=" + tmpFile, "-I" + protoRootPath}
	args = append(args, protoFiles...)

	cmd := exec.Command("protoc", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	marshalledDescriptorSet, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}

	descriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(marshalledDescriptorSet, &descriptorSet)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func findMessageDescriptor(registry *protoregistry.Files, protoFullName protoreflect.FullName) (protoreflect.MessageDescriptor, error) {
	var protoMessage protoreflect.MessageDescriptor

	registry.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		messages := desc.Messages()

		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			if message.FullName() == protoFullName {
				protoMessage = message
				return false
			}
		}

		return true
	})

	if protoMessage != nil {
		return protoMessage, nil
	} else {
		return nil, errors.New("proto not found")
	}
}

func decodeBase64(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal("error:", err)
	}
	return data
}

func consumeKTStream(marshalOpt protojson.MarshalOptions, protoRootPath string, protoName string) {
	protoMsg := protoMessage(marshalOpt, protoRootPath, protoName)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		message := s.Text()

		var result map[string]interface{}
		var data []byte
		var jsonBytes []byte
		var tmp map[string]interface{}

		err := json.Unmarshal([]byte(message), &result)
		if err != nil {
			log.Fatal("error:", err)
		}

		data = decodeBase64(result["value"].(string))

		err = proto.Unmarshal(data, protoMsg)
		if err != nil {
			log.Fatal(err)
		}

		jsonBytes, err = marshalOpt.Marshal(protoMsg)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: avoid duplicate conversion
		err = json.Unmarshal([]byte(jsonBytes), &tmp)
		if err != nil {
			log.Fatal("error:", err)
		}

		result["value"] = tmp

		jsonString, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(jsonString))
	}
}

func consumeStdin(marshalOpt protojson.MarshalOptions, protoRootPath string, protoName string) {
	data, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	marshalOpt.Indent = "  "
	protoMsg := protoMessage(marshalOpt, protoRootPath, protoName)

	err = proto.Unmarshal(data, protoMsg)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := marshalOpt.Marshal(protoMsg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}

func protoMessage(marshalOpt protojson.MarshalOptions, protoRootPath string, protoName string) *dynamicpb.Message {
	protoFiles, err := listProtoFiles(protoRootPath, "*.proto")
	if err != nil {
		log.Fatal(err)
	}

	registry, err := createProtoRegistry(protoRootPath, protoFiles)
	if err != nil {
		log.Fatal(err)
	}

	messageDescriptor, err := findMessageDescriptor(registry, protoreflect.FullName(protoName))
	if err != nil {
		log.Fatal(err)
	}

	return dynamicpb.NewMessage(messageDescriptor)
}
