#Cluster Library
`CLUSTER` is a library that provide you a facility to create server instance.You can create several servers in `Distributed Environment`.`Cluster` has a very easy client API.One need not to understand the internal details to create a instance of server and to run a network communication at scale.
#Client API
Client API is quite simpler to implement.Methods we one needs to implement are:
##Methods:
```
New()      New Method to create an instance of server
Outbox()   Outbox Provide you a way to either broadcast or send a message to peer
Inbox()    Inbox is where you can get a message
Envelope Envelope is not a Method but its a group of Peer id and message of any type

```

##Instructions for Client to use
'Cluster'
 To Instantiate cluster one needs to import "github.com/rpraveenverma/cluster" to main program and has to call Cluster.New(parameters). In Parameters first argument is server id and second argument is configuration.json file.
`Client Example Program to clearify the use`
```
package main

import (

	"github.com/rpraveenverma/cluster"
)
func main() {

	s1:=clusterlib.New(100,"configurationfile.json")

	&cluster.Envelope{Pid:cluster.BROADCAST,Msg: "hello there"}   Pass This reference to s1.Outbox() that provide you a channel to communicate
 	
	s2.Inbox() provides you channel that is used to receive the envelope. Say `envelope` is instance of Envelope.Now you can access Message using `envelope.Msg`
 	
}
```

# Install
$ shows your terminal representation.
```
To get the Cluster repository you need to type
$ go get github.com/rpraveenverma/cluster
```
#Feedback
To submit a feedback mail me on [rpraveenverma@gmail.com]
