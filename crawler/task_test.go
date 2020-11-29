package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	task := NewTask(CrawlJob{
		Link: string(rune(0x7f)),
	})
	assert.Panics(t, func() {
		task.Process(nil)
	})
}
