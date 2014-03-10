package main

import (

	"github.com/rpraveenverma/cluster"
	"fmt"

	"time"
)

func main() {

	var s [20]cluster.Servernode
	j:=100
	for i:=0;i<20;i++{
		s[i]=cluster.New(j,"configurationfile.json")
		j++
		time.Sleep(time.Millisecond)
	}
	/*
	s1:=cluster.New(100,"configurationfile.json")
	s2:=cluster.New(101,"configurationfile.json")
	s3:=cluster.New(102,"configurationfile.json")
	s4:=clu"ster.New(105,"configurationfile.json")
	s6:=cluster.New(106,"configurationfile.json")
	s7:=cluster.New(107,"configurationfile.json")
	s8:=cluster.New(108,"configurationfile.json")
	s9:=cluster.New(109,"configurationfile.json")
	s10:=cluster.New(110,"configurationfile.json")
	s11:=cluster.New(111,"configurationfile.json")
	s12:=cluster.New(112,"configurationfile.json")
	s13:=cluster.New(113,"configurationfile.json")
	s14:=cluster.New(114,"configurationfile.json")
	s15:=cluster.New(115,"configurationfile.json")
	s16:=cluster.New(116,"configurationfile.json")
	s17:=cluster.New(117,"configurationfile.json")
	s18:=cluster.New(118,"configurationfile.json")
*/
	s[0].Outbox() <- &cluster.Envelope{Pid:cluster.BROADCAST,Msg: "hello there"}
/*
	for i:=1;i<20;i++{
		println(i)
		fmt.Println((<-s[i].Inbox()).Msg)
	}
/**/
	fmt.Println((<-s[1].Inbox()).Msg)
	fmt.Println((<-s[2].Inbox()).Msg)
	fmt.Println((<-s[3].Inbox()).Msg)
	fmt.Println((<-s[4].Inbox()).Msg)
	fmt.Println((<-s[5].Inbox()).Msg)

}
