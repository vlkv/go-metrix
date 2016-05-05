package metrix

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrix1(t *testing.T) {
	const fileName = "./metrix.txt"
	defer os.Remove(fileName)
	
	MetrixInstance = CreateMetrix(fileName, 200*time.Millisecond)
	assert.NotNil(t, MetrixInstance)
	defer MetrixInstance.Destroy()

	AddMetrixValue("metric1", 10)
	SetMetrixValue("metric2", 200)
	AddMetrixValue("metric1", 5)

	time.Sleep(time.Second)

	data, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)

	assert.True(t, "metric1 = 15\nmetric2 = 200\n" == string(data[:]) || "metric2 = 200\nmetric1 = 15\n" == string(data[:]) )
}
