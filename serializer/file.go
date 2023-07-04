package serializer

import (
	"google.golang.org/protobuf/proto"
	"io/ioutil"
)

func WriteProtobufToJsonFile(message proto.Message, fileName string) error {
	data, err := ProtobufToJson(message)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteProtobufToBinaryFile(message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadProtobufFromBinaryFile(fileName string, message proto.Message) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(file, message)
	if err != nil {
		return err
	}
	return nil
}
