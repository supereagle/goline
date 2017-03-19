FROM golang:1.6.2

ENV PROJECT_DIR /go/src/github.com/supereagle/goline
EXPOSE 8080

COPY . $PROJECT_DIR
COPY cmd/config-example.json $PROJECT_DIR/cmd/config.json
WORKDIR $PROJECT_DIR/cmd

RUN go build -v -a -o /go/bin/goline && \
    chmod +x /go/bin/goline

ENTRYPOINT ["goline"]
