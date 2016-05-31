package metrix

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"
	"github.com/vlkv/go-util"
)

type Metrix interface {
	add(name string, increment int64)
	set(name string, value int64)
	setCalcValue(name string, calcFun func(prevValues map[string]int64, values map[string]int64)int64)
	Destroy()
}

func AddMetrixValue(name string, inc int64){
	if MetrixInstance!=nil{
		MetrixInstance.add(name, inc)
	}
}

func SetMetrixValue(name string, value int64){
	if MetrixInstance!=nil{
		MetrixInstance.set(name, value)
	}
}

func SetMetrixCalcValue(name string, calcFun func(map[string]int64, map[string]int64)int64){
	if MetrixInstance!=nil{
		MetrixInstance.setCalcValue(name, calcFun)
	}
}

var MetrixInstance Metrix

type metrixImpl struct {
	util.ActiveObject

	file       string
	values     map[string]int64
	prevValues map[string]int64
	calcFuns   map[string]func(map[string]int64, map[string]int64)int64
}
var _ Metrix = (*metrixImpl)(nil)

func CreateMetrix(file string, flushInterval time.Duration) Metrix {
	this := new(metrixImpl)
	this.file = file
	this.values = make(map[string]int64)
	this.calcFuns = make(map[string]func(map[string]int64,map[string]int64)int64)
	this.prevValues = make(map[string]int64)
	this.ActiveObject.Create2(nil, 10000)
	this.runTimer(flushInterval)
	return this
}

func (this *metrixImpl) Destroy(){
	this.ActiveObject.Destroy()
}

func (this *metrixImpl) runTimer(flushInterval time.Duration) {
	time.AfterFunc(flushInterval, func() {
		defer func() { recover() }()
		this.ExecuteAsync(this.flush)
		this.runTimer(flushInterval)
	})
}

func (this *metrixImpl) add(name string, increment int64) {
	this.ExecuteAsync(func() { this.addImpl(name, increment) })
}

func (this *metrixImpl) set(name string, value int64) {
	this.ExecuteAsync(func() { this.setImpl(name, value) })
}

func(this *metrixImpl) setCalcValue(name string, calcFun func(map[string]int64, map[string]int64)int64) {
	this.ExecuteAsync(func() { this.setCalcValueImpl(name, calcFun) })
}

func (this *metrixImpl) addImpl(name string, increment int64) {
	this.values[name] += increment
}

func (this *metrixImpl) setImpl(name string, value int64) {
	this.values[name] = value
}

func (this *metrixImpl) setCalcValueImpl(name string, calcFun func(map[string]int64, map[string]int64)int64) {
	this.calcFuns[name] = func (prevValues map[string]int64, values map[string]int64) (res int64) {
		return calcFun(prevValues, values)
	}
}

func (this *metrixImpl) flush() {
	if len(this.values) == 0 {
		return
	}

	buf := new(bytes.Buffer)
	for k, v := range this.values {
		fmt.Fprintf(buf, "%v = %v\n", k, v)
	}
	for k, v := range this.calcFuns {
		fmt.Fprintf(buf, "%v = %v\n", k, v(this.prevValues, this.values))
	}
	ioutil.WriteFile(this.file, buf.Bytes(), 0666)

	for k, v := range this.values {
		this.prevValues[k] = v
	}
}
