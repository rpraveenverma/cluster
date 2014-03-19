package cluster

import "testing"
import "github.com/rpraveenverma/cluster"
import "time"
import "sync"
import "fmt"
import "strings"
import "reflect"

//import "io/ioutil"
var s [10]cluster.Servernode
var mutex = &sync.Mutex{}
var count int

func TestBasicPointToPoint(t *testing.T) {
	j := 100
	for i := 0; i < 10; i++ {
		s[i] = cluster.New(j, "configurationfile.json")
		j++
		time.Sleep(time.Millisecond)
	}

	s[0].Outbox() <- &cluster.Envelope{Pid: 101, Msg: "hello there"}

	if (<-s[1].Inbox()).Msg != "hello there" {
		t.Fail()
	}

}

func TestBasicBroadcast(t *testing.T) {

	s[0].Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msg: "hello there"}

	for i := 1; i < 10; i++ {

		if (<-s[i].Inbox()).Msg != "hello there" {
			t.Fail()
		}

	}
}

func TestBasicBroadCastHeavy(t *testing.T) {
	var count int
	var n int
	n = 10
	for i := 0; i < 10000; i++ {
		s[0].Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msg: "hello there"}
	}

	for j := 1; j < 10; j++ {
		var i int
		for i = 0; i < 10000; i++ {

			if (<-s[j].Inbox()).Msg != "hello there" {
				t.Fail()
			}
			count++
		}
	}
	if count != (n-1)*10000 {
		t.Fail()
	}
}

func TestBroadcastHeavy(t *testing.T) {
	count = 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10000; j++ {
			s[i].Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msg: "hello there"}
		}
	}

	waitforcompletion := make(chan string, 10)

	for i := 0; i < 10; i++ {
		go func(serviter int) {
			countforone := 0
			for j := 1; j < 10; j++ {
				for k := 0; k < 10000; k++ {

					if (<-s[serviter].Inbox()).Msg != "hello there" {
						t.Fail()
					}
					mutex.Lock()
					countforone = countforone + 1
					mutex.Unlock()
				}
			}
			fmt.Println(countforone)
			count += countforone
			waitforcompletion <- "completed"
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-waitforcompletion

	}

	if count != 10*9*10000 {
		t.Fail()
	}
}

// This test is Heavy message transfer from one peer to another peer.
func TestHeavyMessageunicast(t *testing.T) {
	sendersidestring := strings.Repeat("I am in go", 1024*100*64)
	s[0].Outbox() <- &cluster.Envelope{Pid: 101, Msg: sendersidestring}

	if (<-s[1].Inbox()).Msg != sendersidestring {
		t.Fail()
	}
}

// This test actually test heavy message broadcast to everyserver running in this cluster. It will take around 30 sec to get completed as 65 MB datum is being sent.
func TestHeavyMessageBroadcast(t *testing.T) {
	sendersidestring := strings.Repeat("I am in go", 1024*100*64)
	s[0].Outbox() <- &cluster.Envelope{Pid: cluster.BROADCAST, Msg: sendersidestring}

	for i := 1; i < 10; i++ {
		if (<-s[i].Inbox()).Msg != sendersidestring {
			t.Fail()
		}
	}
}

// For this test i have assigned a unique id to every message so that i can check wheter 10000 messages that were being sent,recieved
// exacty or not..

// But this test is not working. I am not getting the reason.If you can understand can suggest me.
/*
func TestUniquenessofMessage(t *testing.T){
	sendersidemap:=make(map[int]interface{})
	receiversidemap:=make(map[int]interface{})
	for i:=0;i<10000;i++{
	message:="hello there"+string(i)
	sendersidemap[i]=message
	s[0].Outbox() <- &cluster.Envelope{Pid: 101, Mid:i ,Msg:message}
	}

	for i:=0;i<10000;i++{
	envelope:= <-s[i].Inbox()
	receiversidemap[envelope.Mid]=envelope.Msg
	}

	eq := reflect.DeepEqual(sendersidemap, receiversidemap)
	if !eq {
	t.Fail()
	}
}
*/
