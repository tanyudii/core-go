package dwarfs

import (
	"github.com/vmihailenco/taskq/v3"
)

type Worker interface {
	GetTasks() []*taskq.TaskOptions
}
