package main

import "fmt"
import "flag"
import "strings"
import "labix.org/v2/mgo/bson"
import zmq "github.com/alecthomas/gozmq"
import "xCloud/common"

var uuid string
var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var workerId *string = flag.String("workerId", "1", "worker id to be reserved and used")
var address string 

func main() {
  flag.Parse();
  uuid = "c1f8cec0-9b38-41a9-8aee-6e31f962ba32"
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  address = fmt.Sprintf("tcp://%s", *ip)
  add := fmt.Sprintf("%s:16653", address)
  socket.Connect(add)

  fmt.Println(*workerId)
  l1 := messages.ReserveWorker{*workerId}
  c1 := messages.Command{"reserveWorker", uuid, messages.ListWorkers{}, messages.MyWorker{}, l1, messages.Exec{}}
  data, _ := bson.Marshal(c1)
  socket.Send(data, 0)
  r, _ := socket.Recv(0)
  fmt.Println(string(r))

  operation := "start"
  cmd := "wget https://raw.github.com/miha-stopar/xCloud/master/examples/multiple-basic/worker1.py"
  l2 := messages.Exec{*workerId, cmd, operation}
  c2 := messages.Command{"execute", uuid, messages.ListWorkers{}, messages.MyWorker{},  messages.ReserveWorker{}, l2}
  data, _ = bson.Marshal(c2)
  socket.Send(data, 0)
  reply, _ := socket.Recv(0)
  fmt.Println(string(reply) + "\n")

  operation = "output"
  cmd = "python worker1.py"
  l2 = messages.Exec{*workerId, cmd, operation}
  c2 = messages.Command{"execute", uuid, messages.ListWorkers{}, messages.MyWorker{},  messages.ReserveWorker{}, l2}
  data, _ = bson.Marshal(c2)
  socket.Send(data, 0)
  reply, _ = socket.Recv(0)
  rep := strings.TrimSpace(string(reply))
  fmt.Println(rep + "\n")
  if string(rep) != "1" {
    fmt.Println("WRONG!")
  } else {
    fmt.Println("RIGHT!")
  }
}




