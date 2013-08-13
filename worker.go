package main

import "fmt"
import "flag"
import "strings"
import "os/exec"
import zmq "github.com/alecthomas/gozmq"

var ip *string = flag.String("ip", "127.0.0.1", "server IP address")
var desc *string = flag.String("desc", "this is a worker with only basic libraries", "worker description")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)
  rcontext, _ := zmq.NewContext()
  rsocket, _ := rcontext.NewSocket(zmq.REQ)
  rsocket.Connect(fmt.Sprintf("%s:5002", address))
  defer rcontext.Close()
  defer rsocket.Close()
  rsocket.Send([]byte(*desc), 0)
  w_id, _ := rsocket.Recv(0)
  worker_id := string(w_id)
  fmt.Println(worker_id)

  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.SUB)
  defer context.Close()
  defer socket.Close()

  socket.SetSubscribe(string(worker_id))
  socket.Connect(fmt.Sprintf("%s:5556", address))

  wcontext, _ := zmq.NewContext()
  wsocket, _ := wcontext.NewSocket(zmq.REQ)
  wsocket.Connect(fmt.Sprintf("%s:6000", address))
  defer wcontext.Close()
  defer wsocket.Close()

  ccontext, _ := zmq.NewContext()
  csocket, _ := ccontext.NewSocket(zmq.REQ)
  csocket.Connect(fmt.Sprintf("%s:6001", address))
  defer ccontext.Close()
  defer csocket.Close()

  for {
    datapt, _ := socket.Recv(0)
    //fmt.Println(datapt)
    temps := strings.Split(string(datapt), " ")
    //fmt.Println(temps)
    cmd := temps[1]
    cmd = strings.Replace(cmd, "\n", "", -1)
    //fmt.Println(cmd)
    if cmd == "checkWorker"{
      csocket.Send([]byte("dummy"), 0)
      _, _ = csocket.Recv(0)
    } else {
      response, err := exec.Command(cmd).Output()
      if err != nil {
        fmt.Println(err)
      }
      wsocket.Send([]byte(response), 0)
      _, _ = wsocket.Recv(0)
    }
  }
}


