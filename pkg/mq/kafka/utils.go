package kafka

import (
	"strconv"
	"strings"
)

// buildKafkaMsgID create a unique id for kafka message
func buildKafkaMsgID(message *Message) string {
	return strings.Join([]string{"topic:" + message.Topic, "partition:" + strconv.Itoa(message.Partition),
		"offset:" + strconv.FormatInt(message.Offset, 10)}, "-")
}
