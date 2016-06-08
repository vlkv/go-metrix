package metrix

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
	"xgit.tradingview.com/go-lib/go-util"
	"github.com/stretchr/testify/assert"
	"strings"
	"fmt"
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

	MetrixInstance.setCalcValue("metric3", func(input CalcFuncInput) int64 {
		return input.Values["metric1"] + input.Values["metric2"]
	})

	time.Sleep(time.Second)

	data, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)

	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")
	fmt.Println(lines)

	assert.True(t, util.FindIndex(len(lines), func(i int) bool {
		return lines[i] == "metric1 = 15"
	}) >= 0)

	assert.True(t, util.FindIndex(len(lines), func(i int) bool {
		return lines[i] == "metric2 = 200"
	}) >= 0)

	assert.True(t, util.FindIndex(len(lines), func(i int) bool {
		return lines[i] == "metric3 = 215"
	}) >= 0)
}
