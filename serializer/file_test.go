package serializer

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"grpc_learning/pb"
	"grpc_learning/sample"
	"testing"
)

func TestWriteProtobufToBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"

	message := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(message, binaryFile)
	require.NoError(t, err)
	message2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, message2)
	require.NoError(t, err)
	require.True(t, proto.Equal(message, message2))
}

func TestWriteProtobufToJsonFile(t *testing.T) {
	t.Parallel()
	jsonFile := "../tmp/laptop.json"
	message := sample.NewLaptop()
	err := WriteProtobufToJsonFile(message, jsonFile)
	require.NoError(t, err)
}
