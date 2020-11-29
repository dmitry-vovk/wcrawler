package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	task := newTask(crawlJob{
		Link: string(rune(0x7f)),
	})
	assert.Error(t, task.Process(nil).Error)
}
