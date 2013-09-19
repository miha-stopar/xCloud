Basic communication between one client and one worker
=====

* Start *client* and reserve some available *worker*.
* Copy *bworker.py* on the reserved *worker*. You might use something like:

::

	start wget https://github.com/miha-stopar/xCloud/blob/master/examples/basic-communication/bworker.py

* Execute the following command on the *client* machine (use the reserved worker id): 

::

	./bclient -ip=192.168.1.10 -workerId=0

Simple *client* program will list the workers and execute *python bworker.py* on the *worker* machine.

