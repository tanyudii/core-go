package dwarfs

import (
	"context"
	"github.com/vmihailenco/taskq/v3"
)

func CreateMessage(ctx context.Context, taskName string, args ...interface{}) *taskq.Message {
	msg := taskq.NewMessage(ctx, args...)
	msg.TaskName = taskName
	return msg
}
