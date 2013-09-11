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
RUN apt-get install -y uuid-runtime
RUN go get github.com/alecthomas/gozmq
RUN apt-get install -y bzr
RUN go get labix.org/v2/mgo/bson
RUN mkdir -p /srv/gocode/src/xCloud
RUN git clone https://github.com/miha-stopar/xCloud /srv/gocode/src/xCloud

RUN apt-get install -y openssh-server supervisor
ADD sshd.conf /etc/supervisor/conf.d/sshd.conf
RUN mkdir -p /var/run/sshd
RUN echo root:krop | chpasswd
EXPOSE 22
# somehow this ENV command doesn't work
# ENV GOPATH /srv/gocode

# the following command works when running docker as a daemon, but it doesn't when running docker interactively:
CMD /usr/bin/supervisord -n 

