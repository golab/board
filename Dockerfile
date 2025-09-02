FROM ubuntu:latest

RUN apt-get update && apt-get -y upgrade

RUN apt-get install -y golang ca-certificates

RUN update-ca-certificates

ADD frontend /root/frontend
ADD backend /root/backend
ADD cmd /root/cmd

ADD go.mod /root/go.mod
ADD go.sum /root/go.sum

RUN sed -i -e 's/localhost/0.0.0.0/g' /root/cmd/main.go

EXPOSE 8080

ADD scripts/run.sh /root
WORKDIR /root
CMD ["bash", "/root/run.sh"]
