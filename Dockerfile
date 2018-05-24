from ubuntu:14.04

#install packages
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    golang

#add eps-conduit user account and set up go environment
RUN useradd -d /var/lib/eps-conduit -m -s /bin/bash eps-conduit
RUN mkdir /var/lib/eps-conduit/.go && chown eps-conduit: /var/lib/eps-conduit/.go
RUN echo "export GOPATH=\$HOME/.go" >> /var/lib/eps-conduit.bashrc && echo "export PATH=\$PATH:\$GOPATH/bin" >> /var/lib/eps-conduit/.bashrc
ENV GOPATH=/var/lib/eps-conduit/.go
RUN mkdir /var/lib/eps-conduit/eps-conduit
ADD main.go /var/lib/eps-conduit/eps-conduit/main.go
RUN chown -R eps-conduit: /var/lib/eps-conduit

#build binary
WORKDIR /var/lib/eps-conduit/eps-conduit
USER eps-conduit
RUN go build

CMD /var/lib/eps-conduit/eps-conduit/eps-conduit
