package storage

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFileStorage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "mdump-")
	if err != nil {
		fmt.Println(err)
	}
	defer os.RemoveAll(tempDir)

	s := NewLocalStorage(tempDir)
	t.Logf("dir path --- %s", tempDir)
	err = s.Save(context.Background(), "test", strings.NewReader("my request"))
	assert.NoError(t, err)

	_, err = os.Stat(fmt.Sprintf("%s/test", tempDir))
	assert.NoError(t, err)
}
