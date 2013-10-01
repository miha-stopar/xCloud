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
import "xCloud/common"

var workers map[string]string = make(map[string]string) // workerId : description
var statusWorkers map[string]string = make(map[string]string) // workerId : status
var clientWorkers map[string]string = make(map[string] string) // client_uuid : worker_id

func waitRegistrations(){
  rcontext, _ := zmq.NewContext()
  rsocket, _ := rcontext.NewSocket(zmq.REP)
  rsocket.Bind(fmt.Sprintf("%s:16652", address))
  defer rcontext.Close()
  defer rsocket.Close()
  for {
    workerDesc, _ := rsocket.Recv(0)
    count := len(workers)
    workerId := strconv.Itoa(count)
    println(workerId)
    workers[string(workerId)] = string(workerDesc)
    statusWorkers[workerId] = "running"
    println("Got worker: ", string(workerDesc))
    rsocket.Send([]byte(workerId), 0)
  }
}

func serve() {
  db, err := sql.Open("mysql", "xcu:xcp@/xcdb")
  if err != nil {
    //fmt.Println("db error: %s", err)
  }
  db.Exec("CREATE TABLE IF NOT EXISTS audit (uuid VARCHAR(36), command VARCHAR(40));")
  db.Exec("CREATE TABLE IF NOT EXISTS uuids (uuid VARCHAR(36));")
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.REP)
  socket.Bind(fmt.Sprintf("%s:16653", address))
  defer context.Close()
  defer socket.Close()
  
  for {
    msg, _ := socket.Recv(0)
    var output messages.Command
    err := bson.Unmarshal(msg, &output)
    if err != nil {
      //fmt.Println("Received error: %s", err)
    }
    uuid := strings.TrimSpace(output.Uuid)
    sCmd := fmt.Sprintf("select * from uuids where uuid='%s'", uuid)
    rows, _ := db.Query(sCmd)
    if rows.Next(){
    } else {
      socket.Send([]byte("please check your uuid"), 0)
      continue
    }
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
 	sCmd := fmt.Sprintf("INSERT INTO audit VALUES ('%s', 'listWorkers')", uuid)
	db.Exec(sCmd)
	data, _ := bson.Marshal(workersRepr)
    	socket.Send(data, 0)
    case "myWorker":
	var reply string
	if val, ok := clientWorkers[uuid]; ok { 
	  reply = fmt.Sprintf("%s | %s | %s", val, statusWorkers[val], workers[val]) 
 	} else {
	  reply = "please reserve worker first"
	}
 	sCmd := fmt.Sprintf("INSERT INTO audit VALUES ('%s', 'myWorker')", uuid)
	db.Exec(sCmd)
    	socket.Send([]byte(reply), 0)
    case "reserveWorker":
	workerId := strings.TrimSpace(output.Reserve.WorkerId)
	if status, ok := statusWorkers[workerId]; !ok{ 
	  socket.Send([]byte("this worker does not exist"), 0)
	  continue
	} else {
	  if status != "running"{
	    socket.Send([]byte("this worker is not running at the moment"), 0)
	    continue
	  }
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
 	sCmd := fmt.Sprintf("INSERT INTO audit VALUES ('%s', 'reserveWorker %s')", uuid, workerId)
	db.Exec(sCmd)
	socket.Send([]byte("ok"), 0)
    case "execute":
	workerId := strings.TrimSpace(output.Execute.WorkerId)
	if len(workerId) == 0{
	  socket.Send([]byte("worker not specified"), 0)
	  continue
	}
	if _, ok := statusWorkers[workerId]; !ok{ 
	  socket.Send([]byte("this worker does not exist"), 0)
	  continue
	}
	if val, ok := clientWorkers[uuid]; !ok{ 
	  socket.Send([]byte("please reserve worker first"), 0)
	  continue
	} else {
	  if val != workerId{
	    socket.Send([]byte("this worker is not available for this client"), 0)
	    continue
	  }
	}
	cmd := strings.TrimSpace(output.Execute.Cmd)
 	reply, err := delegate(workerId, output.Execute.OpType, cmd)
	//fmt.Println(reply)
 	sCmd := fmt.Sprintf("INSERT INTO audit VALUES ('%s', 'execute %s %s')", uuid, workerId, output.Execute.Cmd)
	db.Exec(sCmd)
        if err != nil {
    	  statusWorkers[workerId] = "disconnected"
	  socket.Send([]byte("no answer"), 0)
        } else {
	  socket.Send([]byte(reply), 0)
	}
    default:
        //fmt.Println("command not known")
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

func delegate(topic string, operation string, cmd string) (string, error) {
  msg := fmt.Sprintf("%s %s %s", topic, operation, cmd)
  //fmt.Println(msg)
  psocket.Send([]byte(msg), 0)
  reply, err := wsocket.Recv(0)
  //fmt.Println(reply)
  wsocket.Send([]byte("dummy"), 0)
  //fmt.Println(err)
  return string(reply), err
}

func checkWorkers(){
  for {
    time.Sleep(1000 * time.Millisecond)
    for ind, _ := range statusWorkers{
      msg := fmt.Sprintf("%s %s", ind, "checkWorker")
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
  psocket.Bind(fmt.Sprintf("%s:16654", address))

  wcontext, _ := zmq.NewContext() // connected to workers
  wsocket, _ = wcontext.NewSocket(zmq.REP)
  //wsocket.SetRcvTimeout(1000 * time.Millisecond)
  defer wcontext.Close()
  defer wsocket.Close()
  wsocket.Bind(fmt.Sprintf("%s:16650", address))

  ccontext, _ := zmq.NewContext() // connected to workers
  csocket, _ = ccontext.NewSocket(zmq.REP)
  csocket.SetRcvTimeout(1000 * time.Millisecond)
  defer ccontext.Close()
  defer csocket.Close()
  csocket.Bind(fmt.Sprintf("%s:16651", address))

  go waitRegistrations()
  go serve()
  go checkWorkers()

  var inp string
  fmt.Scanln(&inp)
}


