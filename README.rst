About
=====

This is a framework which exposes your Docker container (called *worker* in this context) to the outside world. 
Others (*clients*) can access its terminal and execute commands inside it.
A central *server* is needed to enable the communication between *clients* and *workers*. In case you are not running *workers* and *clients* on the same subnet, the *server* needs to be running on a public IP.

The system could be handy if you are running out of resources on your local machine and would like to delegate some tasks on *workers* (if they are available of course).

.. image:: https://bitbucket.org/miha_stopar/godocker/raw/tip/godocker.png


Run server
=====

* install Go (for this and the steps below you can check docker/Dockerfile and see which libraries are needed)
* install ZeroMQ
* install gozmq (you will need Git for this):

::

	go get github.com/alecthomas/gozmq

* install gobson:

::

	go get labix.org/v2/mgo/bson
	
* install go-sql-driver:

::

	go get github.com/go-sql-driver/mysql

* download godocker
* build server.go:

::

	go build server.go

* run *server* (example command for running server on a local subnet): 

::

	./server -ip=192.168.1.12


Run worker
=====

* install docker
* build docker container from a Dockerfile (execute the following command when in folder *godocker/docker*) - you might add some libraries to be installed during the process and you need to modify the *ip* argument at the end of the Dockerfile (it has to be the *server* ip), also you might change the *desc* argument to reflect your changes related to the installed libraries:

::

	docker build -t godocker-img .

If some problems appear when building container, the following command executed on host might help:

::

        sysctl -w net.ipv4.ip_forward=1

* run docker container:

::

	docker run -d godocker-img

*Worker* will be automatically started. You can connect to the container using SSH:

::

        ssh root@localhost -p 49164

Find out the port number using the command:

::

        docker ps

Run client
=====

* install Go
* install ZeroMQ
* install gozmq and gobson
* download godocker
* build worker.go:

::

	go build worker.go

* run *worker* - ip has to be the IP of a *server*: 

::

	./client -ip=192.168.1.12

* list available *workers* (type into *worker* console):

::

	listWorkers

* reserve worker:

::

	reserveWorker 0

* list my *worker* (one *client* can connect to one *worker*, but you can have many *clients*)

::

	myWorker

* execute something on *worker*:

::

	execute 0 ls -al	


.. image:: https://bitbucket.org/miha_stopar/godocker/raw/tip/godocker_screenshot.png

Note
=====

Use ZeroMQ version 2.2 or higher (due to SetRcvTimeout call in server.go).



