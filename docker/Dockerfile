# VERSION 0.1

FROM ubuntu

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get install -y python-software-properties
RUN add-apt-repository -y ppa:chris-lea/zeromq
RUN apt-get update
RUN apt-get upgrade

RUN apt-get install -y libzmq-dev
RUN apt-get install -y golang
RUN apt-get install -y pkg-config
RUN apt-get install -y git-core
RUN apt-get install -y gcc
RUN apt-get install -y vim
RUN apt-get install -y wget
RUN apt-get install -y apache2 libconfig-general-perl librrds-perl libregexp-common-perl libhtml-parser-perl collectd-core
RUN go get github.com/alecthomas/gozmq
#RUN apt-get install -y bzr # currently not needed for worker
#RUN go get labix.org/v2/mgo/bson # currently not needed for worker

RUN wget https://raw.github.com/miha-stopar/godocker/master/worker.go
RUN go build worker.go 

RUN apt-get install -y openssh-server supervisor
ADD sshd.conf /etc/supervisor/conf.d/sshd.conf
ADD collectd-supervisor.conf /etc/supervisor/conf.d/collectd-supervisor.conf
ADD collectd.conf /etc/collectd/collectd.conf
RUN mkdir -p /var/run/sshd
RUN echo root:krop | chpasswd
EXPOSE 22

CMD ./worker -ip=192.237.169.230 -desc="this is a worker with libraries: ..."
# the following command works when running docker as a daemon, but it doesn't when running docker interactively:
CMD /etc/init.d/collectd restart
CMD /usr/bin/supervisord -n 

