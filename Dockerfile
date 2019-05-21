FROM golang:1.11

WORKDIR /go/src/github.com/bitum-project/bitumd
COPY . .

RUN env GO111MODULE=on go install . ./cmd/...

EXPOSE 9208

CMD [ "bitumd" ]
