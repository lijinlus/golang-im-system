package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct{
	Ip string
    Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	Message chan string
}

func NewServer(ip string, port int) *Server{
	server:=&Server{
		Ip:ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

func(this *Server)ListenMessager(){
	for{
		msg:=<-this.Message
		this.mapLock.Lock()
		for _,cli:=range this.OnlineMap{
			cli.C<-msg
		}
       this.mapLock.Unlock()
	}
}
func(this *Server)BroadCast(user *User,msg string){
	sendMsg :="["+user.Addr+"]"+user.Name+":"+msg
	this.Message<-sendMsg
}
func(this *Server)Handler(conn net.Conn){
	user :=NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
 	this.BroadCast(user,"login ok")
 fmt.Println("conn connection successful")
}


func (this *Server)Start(){
  listener,err:=net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
if err!=nil {
	fmt.Println("net.listen err:",err)
	return
}
defer listener.Close()

go this.ListenMessager()
for {
	conn,err:=listener.Accept()
	if err!=nil{
		fmt.Println("lisitener accept err",err)
		continue
	}
	go this.Handler(conn)
}

}