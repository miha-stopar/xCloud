package main

import "fmt"
import "strings"
import "flag"
import "strconv"
import "time"
import "labix.org/v2/mgo/bson"
import zmq "github.com/alecthomas/gozmq"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var workers map[string]string = make(map[string]string) // workerId : description
var statusWorkers map[string]string = make(map[string]string) // workerId : status
var clientWorkers map[string]string = make(map[string] string) // client_uuid : worker_id

func waitRegistrations(){
  rcontext, _ := zmq.NewContext()
  rsocket, _ := rcontext.NewSocket(zmq.REP)
  rsocket.Bind(fmt.Sprintf("%s:5002", address))
  defer rcontext.Close()
  defer rsocket.Close()
  for {
    workerDesc, _ := rsocket.Recv(0)
    count := len(workers)
    workerId := strconv.Itoa(count)
    println(workerId)
    workers[string(workerId)] = string(workerDesc)
    statusWorkers[string(workerId)] = "running"
    println("Got worker: ", string(workerDesc))
    rsocket.Send([]byte(workerId), 0)
  }
}

func serve() {
  db, err := sql.Open("mysql", "kropdb:kropdb@/krop")
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("%s:5000", address))
  defer context.Close()
  defer socket.Close()
  fmt.Println("serving ... ")
  
  for {
    msg, _ := socket.Recv(0)
    var output messages.Command
    err := bson.Unmarshal(msg, &output)
    fmt.Println("Received error: %s", err)

    switch string(output.Name) {
    case "listWorkers":
	workersRepr :=  make(map[string] string)
        for ind, desc := range workers{
	  ownerClient := getClientBehindWorker(ind)
	  occupied := "available"
	  if ownerClient != ""{
	    occupied = "reserved"
	  }
	  workersRepr[ind] = fmt.Sprintf("%s | %s | %s", statusWorkers[ind], occupied,  desc) 
	} 
	data, _ := bson.Marshal(workersRepr)
    	socket.Send(data, 0)
    case "myWorker":
	var reply string
 	uuid := strings.TrimSpace(output.MyWorker.Uuid)
	fmt.Println(clientWorkers)
	fmt.Println(uuid)
	if val, ok := clientWorkers[uuid]; ok { 
	  reply = fmt.Sprintf("%s | %s | %s", val, statusWorkers[val], workers[val]) 
 	} else {
	  reply = "please reserve worker first"
	}
    	socket.Send([]byte(reply), 0)
    case "reserveWorker":
	fmt.Println("reserveWorker")
	workerId := strings.TrimSpace(output.Reserve.WorkerId)
        uuid := strings.TrimSpace(output.Reserve.Uuid)
	if _, ok := statusWorkers[workerId]; !ok{ 
	  socket.Send([]byte("this worker does not exist"), 0)
	  continue
	}
	if _, ok := clientWorkers[uuid]; ok{ 
	  socket.Send([]byte("this client already has a worker"), 0)
	  continue
	} else {
	  ownerClient := getClientBehindWorker(workerId)
	  if ownerClient != ""{
	    socket.Send([]byte("this worker is already reserved"), 0)
	    continue
 	  }
	}
	clientWorkers[uuid] = strings.TrimSpace(workerId)
	socket.Send([]byte("ok"), 0)
    case "execute":
	workerId := output.Execute.WorkerId
 	uuid := strings.TrimSpace(output.Execute.Uuid)
	if _, ok := statusWorkers[workerId]; !ok{ 
	  socket.Send([]byte("this worker does not exist"), 0)
	  continue
	}
	fmt.Println(clientWorkers)
	fmt.Println(uuid)
	if val, ok := clientWorkers[uuid]; !ok{ 
	  socket.Send([]byte("please reserve worker first"), 0)
	  continue
	} else {
	  if val != workerId{
	    socket.Send([]byte("this worker is not available for this client"), 0)
	    continue
	  }
	}
 	reply, err := delegate(workerId, output.Execute.Cmd)
	fmt.Println(reply)
        if err != nil {
    	  statusWorkers[workerId] = "disconnected"
	  socket.Send([]byte("no answer"), 0)
        } else {
	  socket.Send([]byte(reply), 0)
	}
    default:
        fmt.Println("command not known")
    }
  }
}

func getClientBehindWorker(workerId string) string {
  var uuid string
  for cUuid, wId := range clientWorkers{
    if wId == workerId {
     uuid = cUuid
     break
    }
  }
  return uuid
}

func delegate(topic string, cmd string) (string, error) {
  msg := fmt.Sprintf("%s %s", topic, cmd)
  fmt.Println(msg)
  psocket.Send([]byte(msg), 0)
  reply, err := wsocket.Recv(0)
  wsocket.Send([]byte("dummy"), 0)
  fmt.Println(err)
  //fmt.Println(reply)
  return string(reply), err
}

func checkWorkers(){
  for {
    time.Sleep(4000 * time.Millisecond)
    for ind, _ := range statusWorkers{
      msg := fmt.Sprintf("%s %s", ind, "checkWorker")
      //fmt.Println(msg)
      psocket.Send([]byte(msg), 0)
      _, err := csocket.Recv(0)
      csocket.Send([]byte("dummy"), 0)
      //fmt.Println(err)
      if err != nil{
	statusWorkers[ind] = "disconnected"
      } else {
        statusWorkers[ind] = "running"
      }
    }
  } 
}

var psocket *zmq.Socket
var wsocket *zmq.Socket
var csocket *zmq.Socket
var ip *string = flag.String("ip", "127.0.0.1", "public IP address of this very computer")
var address string

func main() {
  flag.Parse();
  address = fmt.Sprintf("tcp://%s", *ip)
  pcontext, _ := zmq.NewContext()
  psocket, _ = pcontext.NewSocket(zmq.PUB)
  defer pcontext.Close()
  defer psocket.Close()
  psocket.Bind(fmt.Sprintf("%s:5556", address))

  wcontext, _ := zmq.NewContext() // connected to workers
  wsocket, _ = wcontext.NewSocket(zmq.REP)
  wsocket.SetRcvTimeout(1000 * time.Millisecond)
  defer wcontext.Close()
  defer wsocket.Close()
  wsocket.Bind(fmt.Sprintf("%s:6000", address))

  ccontext, _ := zmq.NewContext() // connected to workers
  csocket, _ = ccontext.NewSocket(zmq.REP)
  csocket.SetRcvTimeout(1000 * time.Millisecond)
  defer ccontext.Close()
  defer csocket.Close()
  csocket.Bind(fmt.Sprintf("%s:6001", address))

  go waitRegistrations()
  go serve()
  go checkWorkers()

  var inp string
  fmt.Scanln(&inp)
}


