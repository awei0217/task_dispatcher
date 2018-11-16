package task_manager

import (
	"testing"
)

func TestTaskManager_StopTask(t *testing.T) {
	GetTaskManger().StopTask(1)
}
