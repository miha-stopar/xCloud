package messages

type Command struct {
  Name string
  Uuid string
  ListWorkers ListWorkers ",omitempty"
  MyWorker MyWorker ",omitempty"
  Reserve ReserveWorker ",omitempty"
  Execute Exec ",omitempty"
}

type ListWorkers struct {
}

type MyWorker struct {
}

type ReserveWorker struct {
    WorkerId string
}

type Exec struct {
    WorkerId string
    Cmd string
    OpType string
}
