FROM golang:1.10-stretch

RUN go get -d -v github.com/gobuffalo/packr/...; \
    cd /go/src/github.com/gobuffalo/packr/; \
    \
    packr_version="1.10.4"; \
    \
    git checkout tags/v"$packr_version"; \
    go install ./...; \
    \
    rm -rf /go/src/*; \
    \
    chmod o+w /home

# since we don't run as root we need a home dir writable
# otherwise go will fail with permission denied when creating $HOME/.cache/ directory
ENV HOME /home

WORKDIR /go/src/github.com/rootkiwi/screen_share_remote_go
