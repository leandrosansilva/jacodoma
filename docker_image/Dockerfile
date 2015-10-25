FROM debian:stretch

MAINTAINER Leandro Santiago <leandrosansilva@gmail.com>

ENV GOROOT "/go"
ENV GOPATH "/gopath"
ENV PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/gopath/bin:/go/bin"

RUN apt-get update && \
    apt-get install -y curl libqt5core5a libqt5gui5 libqt5network5 libqt5qml5 libqt5quick5 libqt5widgets5 libgit2-23 git && \
    curl -L https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar zxvf - -C /  && \
    mkdir /gopath && \
    go get -v github.com/tools/godep
    
ADD docker_image/init /init

VOLUMES ["/source"]

ENTRYPOINT ["/init"]
