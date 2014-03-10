package cluster

import (
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"io/ioutil"
	"strconv"
)

type server struct {
	Serverid int `json:"serverid"`
	Portno   int `json:"portno"`
}

type serverlist struct {
	RunningServer []server
}

const (
	BROADCAST int = -1
)

type Envelope struct {
	// On the sender side, Pid identifies the receiving peer. If instead, Pid is
	// set to cluster.BROADCAST, the message is sent to all peers. On the receiver side, the
	// Id is always set to the original sender. If the Id is not found, the message is silently dropped
	Pid int
	// An id that globally and uniquely identifies the message, meant for duplicate detection at
	// higher levels. It is opaque to this package.

	// the actual message.
	Msg interface{}
}

type Server interface {
	// Id of this server
	Pid() int

	// array of other servers' ids in the same cluster
	Peers() map[int]int

	// the channel to use to send messages to other peers
	// Note that there are no guarantees of message delivery, and messages
	// are silently dropped
	Outbox() chan *Envelope

	// the channel to receive messages from other peers.
	Inbox() chan *Envelope
}

type Servernode struct {
	sid    int
	pid    int
	Ownsocks *zmq.Socket
	peers  map[int]int
	outbox chan *Envelope
	inbox  chan *Envelope
	peersocks map[int]*zmq.Socket
}

func (s Servernode) Pid() int {
	return s.pid
}

func (s Servernode) Peers() map[int]int {
	return s.peers
}

func (s Servernode) Outbox() chan *Envelope {
	return s.outbox
}

func (s Servernode) Inbox() chan *Envelope {
	return s.inbox
}

var mapofservers map[int]int

func (s Servernode)SendMessage(){
	var tempinstance Envelope

	for {
			tempinstance = *<-s.outbox
			msg,err:=json.Marshal(tempinstance)
			if err!=nil{
				panic(err)
				}
			if tempinstance.Pid==BROADCAST{
			for _,sock:= range s.peersocks{
				sock.SendBytes(msg,zmq.DONTWAIT)
				}

			}else{

				tempsock:=s.peersocks[mapofservers[tempinstance.Pid]]

				tempsock.SendBytes(msg,zmq.DONTWAIT)
			}

	}
}

func (s Servernode)RecieveMessage(responder *zmq.Socket){
	for {
	msgbytes,err:=responder.RecvBytes(0)

	var data Envelope
			if err!=nil{
				panic(err)
			}
	err1:=json.Unmarshal(msgbytes,&data)
		if err1 !=nil {
			panic("Unmarshling the data")
		}
	s.inbox <- &data
	}
}
func New(serverid int,filename string) Servernode {

	var data serverlist
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		panic("Error parsing json File")
	}

	mapofservers = make(map[int]int)
	for i := 0; i <len(data.RunningServer); i++ {

		mapofservers[data.RunningServer[i].Serverid] = data.RunningServer[i].Portno
	}

	var s Servernode
	s.peers=make(map[int]int)
	s.peersocks=make(map[int]*zmq.Socket)
	_, exist := mapofservers[serverid]
	if exist != true {
		panic("Error in creation of new server instance")
	}
	s.sid = serverid
	s.pid = mapofservers[serverid]
	s.outbox = make(chan *Envelope)
	s.inbox = make(chan *Envelope)
	var responder *zmq.Socket
	var requester *zmq.Socket

	for id, port:= range mapofservers{
		if id ==serverid {
			responder, err= zmq.NewSocket(zmq.PULL)
			s.Ownsocks = responder
			responder.Bind("tcp://*:"+strconv.Itoa(port))

			if err != nil {
				panic("Problem Binding Server")
			}
		}else{
		requester, err = zmq.NewSocket(zmq.PUSH);
				if err != nil {
						panic("Problem Binding Server")
				}
		requester.Connect("tcp://localhost:"+strconv.Itoa(port))
		s.peers[id]=port
		s.peersocks[port]=requester
	}
	}
	go s.SendMessage()
	go s.RecieveMessage(responder)
	return  s
	}
