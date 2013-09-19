package main

import "fmt"
import "flag"
import "bufio"
import "os"
import "strings"
import "os/exec"
import "labix.org/v2/mgo/bson"
import zmq "github.com/alecthomas/gozmq"
import "xCloud/common"

func enterCmd(socket *zmq.Socket){
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter command: ")
  command, _ := reader.ReadString('\n')
  parts := strings.Split(string(command), " ")
  if strings.Contains(parts[0], "listWorkers") {
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
  } else if strings.Contains(parts[0], "myWorker") {
    l := messages.MyWorker{uuid}
    c := messages.Command{"myWorker", messages.ListWorkers{}, l, messages.ReserveWorker{}, messages.Exec{}}
    data, _ := bson.Marshal(c)
    socket.Send(data, 0)
    reply, _ := socket.Recv(0)
    fmt.Println("My worker:")
    fmt.Printf(string(reply) + "\n\n")
  } else if strings.Contains(parts[0], "reserveWorker") {
    if len(parts) < 2 {
      fmt.Println("not enough arguments\n")
    } else {
      workerId = parts[1]
      l := messages.ReserveWorker{string(workerId), string(uuid)}
      c := messages.Command{"reserveWorker", messages.ListWorkers{}, messages.MyWorker{}, l, messages.Exec{}}
      data, _ := bson.Marshal(c)
      socket.Send(data, 0)
      reply, _ := socket.Recv(0)
      fmt.Println(string(reply) + "\n")
  }
  } else if (strings.Contains(parts[0], "start") || strings.Contains(parts[0], "output")){
    if len(parts) < 2 {
      fmt.Println("not enough arguments\n")
    } else {
      operation := parts[0]
      cmd := strings.Join(parts[1:], " ")
      fmt.Println(workerId)
      l := messages.Exec{workerId, cmd, operation, uuid}
      c := messages.Command{"execute", messages.ListWorkers{}, messages.MyWorker{},  messages.ReserveWorker{}, l}
      data, _ := bson.Marshal(c)
      socket.Send(data, 0)
      reply, _ := socket.Recv(0)
      fmt.Println(string(reply) + "\n")
    }
  } else {
    fmt.Println("command not found\n")
  }
  enterCmd(socket)
}

var uuid string
var workerId string
var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
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
  enterCmd(socket)
}
