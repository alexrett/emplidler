FROM ubuntu:18.04

RUN apt-get update && apt-get install -y locales wget tar gcc libxss-dev libgtk-3-dev libappindicator3-dev libwebkit2gtk-4.0-dev && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8

RUN wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz && tar -xvf go1.14.linux-amd64.tar.gz && mv go /usr/local && rm -f go1.14.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN go version

