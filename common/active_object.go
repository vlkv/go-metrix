package common

import (
	"runtime/debug"
	"time"
	log "github.com/Sirupsen/logrus"
)

type ActiveObject struct {
	chStopWork       chan interface{}
	cmdCh            chan func()
	messageProcessor func() bool
}

func (this *ActiveObject) Create(messageProcessor func() bool) {
	this.Create2(messageProcessor, 0)
}

func (this *ActiveObject) Create2(messageProcessor func() bool, cmdPoolSize int) {
	this.chStopWork = make(chan interface{})
	this.cmdCh = make(chan func(), cmdPoolSize)
	this.messageProcessor = messageProcessor
	go this.run()
}

func (this *ActiveObject) ExecuteAsync(f func()) {
	this.cmdCh <- f
}

func (this *ActiveObject) ExecuteSync(f func()) {
	waitCh := make(chan interface{})
	this.ExecuteAsync(func() {
		defer func() {
			err := recover()
			if err != nil {
				stack := string(debug.Stack()[:])
				log.Errorf("ActiveObject.ExecuteSync(%v) panic:\n%v\n%v\n", f, err, stack)
			}
			waitCh <- err
		}()
		f()
	})
	err := <-waitCh
	if err != nil {
		panic(err)
	}
}

func (this *ActiveObject) defaultMsgProcess() {
	for {
		select {
		case <-this.chStopWork:
			return
		case cmd := <-this.cmdCh:
			cmd()
		}
	}
}

func (this *ActiveObject) withHandlerMsgProcess() {
	for {
		select {
		case <-this.chStopWork:
			return
		case cmd := <-this.cmdCh:
			cmd()
		default:
			if !this.messageProcessor() {
				time.Sleep(1)
			}
		}
	}
}

func (this *ActiveObject) run() {
	defer func() {
		err := recover()
		if nil != err {
			log.Fatalf("ActiveObject.run panicing:%v\n%v\n", err, string(debug.Stack()[:]))
		}
	}()

	if this.messageProcessor == nil {
		this.defaultMsgProcess()
	} else {
		this.withHandlerMsgProcess()
	}
}

func (this *ActiveObject) Destroy() {
	if nil != this.chStopWork {
		this.chStopWork <- nil
		close(this.chStopWork)
	}
	this.chStopWork = nil
	close(this.cmdCh)
}
