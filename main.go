package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const version = "0.1.0"
const usage = "Usage: proto2json <proto-file>"

func readDynamically(in []byte) {
	registry, err := createProtoRegistry(".", "person.proto")
	if err != nil {
		panic(err)
	}

	desc, err := registry.FindFileByPath("person.proto")
	if err != nil {
		panic(err)
	}
	fd := desc.Messages()
	addressBook := fd.ByName("Person")

	msg := dynamicpb.NewMessage(addressBook)
	err = proto.Unmarshal(in, msg)
	jsonBytes, err := protojson.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
	if err != nil {
		panic(err)
	}
}

func createProtoRegistry(srcDir string, filename string) (*protoregistry.Files, error) {
	// Create descriptors using the protoc binary.
	// Imported dependencies are included so that the descriptors are self-contained.
	tmpFile := filename + "-tmp.pb"
	cmd := exec.Command("protoc",
		"--include_imports",
		"--descriptor_set_out="+tmpFile,
		"-I"+srcDir,
		path.Join(srcDir, filename))

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

func main() {

	// elliot := &Person{
	// 	Name: "Elliot",
	// 	Age:  24,
	// }

	// data, err := proto.Marshal(elliot)
	// if err != nil {
	// 	log.Fatal("marshaling error: ", err)
	// }

	// // printing out our raw protobuf object
	// fmt.Println(data)

	// // let's go the other way and unmarshal
	// // our byte array into an object we can modify
	// // and use
	// newElliot := &Person{}
	// err = proto.Unmarshal(data, newElliot)
	// if err != nil {
	// 	log.Fatal("unmarshaling error: ", err)
	// }

	// // print out our `newElliot` object
	// // for good measure
	// fmt.Println(newElliot.GetAge())
	// fmt.Println(newElliot.GetName())

	// json := protojson.Format(elliot)
	// fmt.Println(json)

	// // write the whole body at once
	// err = ioutil.WriteFile("data.bin", data, 0644)
	// if err != nil {
	// 	log.Fatal("unmarshaling error: ", err)
	// }

	// read the whole file at once
	data, err := ioutil.ReadFile("data.bin")
	if err != nil {
		panic(err)
	}

	readDynamically(data)
}
