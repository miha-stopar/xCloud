package messages

type Command struct {
  Name string
  ListWorkers ListWorkers ",omitempty"
  MyWorker MyWorker ",omitempty"
  Reserve ReserveWorker ",omitempty"
  Execute Exec ",omitempty"
}

type ListWorkers struct {
    Uuid string
}

type MyWorker struct {
    Uuid string
}

type ReserveWorker struct {
    WorkerId string
    Uuid string
}

type Exec struct {
    WorkerId string
    Cmd string
    Uuid string
}
