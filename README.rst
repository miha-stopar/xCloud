About
=====

NOTE: xCloud is under development.

Using xCloud:

* in case your computer has some free resources you can start Docker container and expose it as a *worker* to the outside world
* in case you lack processing power on your local machine you might check if some *workers* are available and if they are - you might locally run *client* to start exploiting their processing power

A central *server* is needed to enable the communication between *clients* and *workers*. In case you are not running *workers* and *clients* on the same subnet, the *server* needs to be running on a public IP.


.. image:: https://raw.github.com/miha-stopar/xCloud/master/img/xcloud.png


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

* download xCloud

* install MySQL, create database *xcdb*, add user *xcu* with password *xcp* (or change DB settings in server.go):

::

	create database xcdb;
	create user 'xcu'@'localhost' identified by 'xcp';
	GRANT ALL PRIVILEGES ON *.* TO 'xcu'@'localhost';

* build server.go:

::

	go build server.go

* run *server* (example command for running server on a local subnet): 

::

	./server -ip=192.168.1.12

* at the moment there is only a primitive authentication available for a client that wants to connect to the server - a client uuid needs to be registered manually in the DB (uuids table is created when *server* is started):

::

	insert into uuids values ("b1f8cec0-9b38-41a9-8aee-6e31f962ba32");


* install `collectd <http://collectd.org/>`_

* configure collectd network plugin - /etc/collectd/collectd.conf has to contain:

::

	LoadPlugin "network"
	<Plugin "network">
  	    Server "192.168.1.12"
	</Plugin>

* install `collectd-web <https://github.com/httpdss/collectd-web>`_
* run collectd-web (graphs will be available on localhost:8888 once the *workers* will be running - see the image below showing some CPU statistics for Docker container)

.. image:: https://raw.github.com/miha-stopar/xCloud/master/img/collectd-web.png

Run worker
=====

* install docker
* download xCloud/docker directory
* modify *ip* parameter at the end of the Dockerfile (needs to be the central *server* IP) and SSH username/password (see *chpasswd*)
* you might add some additional libraries to be installed inside worker (see Dockerfile) and you might change the *worker* description accordingly (at the end of Dockerfile)
* modify *Server* parameter for network plugin in collectd.conf (needs to be the central *server* IP)
* build docker container from a Dockerfile (execute the following command when in folder *xCloud/docker*):

::

	docker build -t xcloud-img .

If some problems appear when building container, the following command executed on host might help:

::

        sysctl -w net.ipv4.ip_forward=1

* run docker container (you might want to limit the CPU and RAM of the container using *-c* and *-m* options):

::

	docker run -d xcloud-img

*Worker* will be automatically started. You can connect to the container using SSH:

::

        ssh root@localhost -p 49164

Find out the port number using the command:

::

        docker ps

Run client
=====

There are two possibilities:

Run client from within Docker container:
-------------------------------

* install docker
* download xCloud/docker-client directory
* build docker container from a Dockerfile (execute the following command when in folder *xCloud/docker-client*):

::

	docker build -t xclient .

* run docker container:

::

	docker run -d xclient

* go into Docker container and set GOPATH variable:

::

	export GOPATH=/srv/gocode

* configure uuid in the client.go inside /srv/gocode/srv/xCloud (uuid needs to be registered manually in the *server* database)
* build client.go:

::

	go build client.go

* start *client*

Run client without Docker container:
-------------------------------

* install Go
* install ZeroMQ
* install gozmq and gobson
* download xCloud
* configure uuid in the client.go (uuid needs to be registered manually in the *server* database)
* build client.go:

::

	go build client.go

* start *client*

How to start and use client
-------------------------------

* run *client* - ip has to be the IP of a *server*: 

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

	output ls -al	


.. image:: https://raw.github.com/miha-stopar/xCloud/master/img/xcloud_screenshot.png

Note
=====

Use ZeroMQ version 2.2 or higher (due to SetRcvTimeout call in server.go).



