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
  rsocket.Connect(fmt.Sprintf("%s:16652", address))
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
  socket.Connect(fmt.Sprintf("%s:16654", address))

  wcontext, _ := zmq.NewContext()
  wsocket, _ := wcontext.NewSocket(zmq.REQ)
  wsocket.Connect(fmt.Sprintf("%s:16650", address))
  defer wcontext.Close()
  defer wsocket.Close()

  ccontext, _ := zmq.NewContext()
  csocket, _ := ccontext.NewSocket(zmq.REQ)
  csocket.Connect(fmt.Sprintf("%s:16651", address))
  defer ccontext.Close()
  defer csocket.Close()

  for {
    datapt, _ := socket.Recv(0)
    st := strings.Replace(string(datapt), "\n", "", -1)
    temps := strings.Split(st, " ")
    if len(temps) < 2 { //todo: check when this happens
      continue 
    }
    if temps[1] == "checkWorker"{
      csocket.Send([]byte("dummy"), 0)
      _, _ = csocket.Recv(0)
    } else {
      cmd := temps[1:]
      opType := cmd[0]
      cmd = cmd[1:]
      fmt.Println(opType)
      command := strings.Join(cmd, " ")
      fmt.Println(command)
      var response []byte
      var err error
      var ecmd *exec.Cmd
      if len(cmd) == 1 {
        ecmd = exec.Command(cmd[0])
      } else if len(cmd) == 2 { 
        ecmd = exec.Command(cmd[0], cmd[1])
      } else if len(cmd) == 3 { 
        ecmd = exec.Command(cmd[0], cmd[1], cmd[2])
      } else if len(cmd) == 4 { 
        ecmd = exec.Command(cmd[0], cmd[1], cmd[2], cmd[3])
      } else if len(cmd) == 5 { 
        ecmd = exec.Command(cmd[0], cmd[1], cmd[2], cmd[3], cmd[4])
      }

      if opType == "output" {
        response, err = ecmd.Output()
      } else if opType == "start" {
        error := ecmd.Start()
	if error != nil {
	  response = []byte("error when starting a command")	
	} else {
	  error = ecmd.Wait()
	  if error != nil {
	    response = []byte("error when executing a command")	
	  } else {
	    response = []byte("command execution finished")	
	  }
	}
      }
      
      if err != nil {
        fmt.Println(err)
      }
      wsocket.Send([]byte(response), 0)
      _, _ = wsocket.Recv(0)
    }
  }
}


