package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobufToJson(message proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		Indent:          "  ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	return marshaler.Marshal(message)
}

func JsonToProtoBuf(data []byte) (proto.Message, error) {
	var message proto.Message
	err := protojson.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
