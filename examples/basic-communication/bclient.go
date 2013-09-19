package main

import "fmt"
import "flag"
import "os/exec"
import "labix.org/v2/mgo/bson"
import zmq "github.com/alecthomas/gozmq"
import "xCloud/common"

var uuid string
var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var workerId *string = flag.String("workerId", "0", "worker id to be reserved and used")
var address string 

func main() {
  flag.Parse();
  u, err := exec.Command("uuidgen").Output()
  if err != nil {
    fmt.Println(err)
  } else {
    uuid = string(u)
  }
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REQ)
  address = fmt.Sprintf("tcp://%s", *ip)
  add := fmt.Sprintf("%s:16653", address)
  socket.Connect(add)

  l := messages.ListWorkers{uuid}
  c := messages.Command{"listWorkers", l, messages.MyWorker{}, messages.ReserveWorker{}, messages.Exec{}}
  data, _ := bson.Marshal(c)
  socket.Send(data, 0)
  reply, _ := socket.Recv(0)
  var output map[string]string
  _ = bson.Unmarshal(reply, &output)
  fmt.Println("Workers:")
  for k,v := range output{
    fmt.Printf("%s: %s\n\n", string(k), v)
  }

  operation := "output"
  cmd := "python bla.py"
  l2 := messages.Exec{*workerId, cmd, operation, uuid}
  c2 := messages.Command{"execute", messages.ListWorkers{}, messages.MyWorker{},  messages.ReserveWorker{}, l2}
  data, _ = bson.Marshal(c2)
  socket.Send(data, 0)
  reply, _ = socket.Recv(0)
  fmt.Println(string(reply) + "\n")

}
