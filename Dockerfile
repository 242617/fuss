FROM golang:1.8

WORKDIR /go/src/github.com/242617/torture
COPY . .
COPY torture.yaml /go

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run", "-config=torture.yaml"]