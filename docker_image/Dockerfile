FROM debian:stretch

MAINTAINER Leandro Santiago <leandrosansilva@gmail.com>

ENV GOROOT="/go" GOPATH="/gopath" PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/gopath/bin:/go/bin"

RUN apt-get update && \
    apt-get install -y curl libqt5core5a libqt5gui5 libqt5network5 libqt5qml5 libqt5quick5 libqt5widgets5 \
    qtbase5-dev qtdeclarative5-dev qtbase5-private-dev \
    libgit2-23 libgit2-dev gcc g++ git \
    pkg-config libqt5qml5 qml-module-qtquick-controls && \
    curl -L https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar zxvf - -C /  && \
    mkdir -p /gopath/src && \
    go get -v github.com/tools/godep && \
    go get -v github.com/smartystreets/goconvey/convey

RUN useradd -m -d /home/jacodoma jacodoma

RUN chown -R jacodoma /gopath 
    
ADD docker_image/init /init

VOLUME ["/gopath/src/jacodoma"]

USER jacodoma

ENTRYPOINT ["/init"]
