FROM gocd/gocd-agent-centos-7:v21.2.0

USER root

RUN yum update -y

RUN yum -y install maven gradle wget

RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && mkdir /golang && tar -xzf go1.16.4.linux-amd64.tar.gz --directory /golang && mv /golang/go /usr/local

COPY client.go .

USER go