# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build

WORKDIR /find_the_flag

COPY ./Find_the_FLAG/go.mod ./

RUN go mod download

COPY ./Find_the_FLAG/*.go ./

RUN go build -o /find_flag find_flag.go

WORKDIR /cryptopals

COPY ./Cryptopals_SET_1/go.mod ./

RUN go mod download

COPY ./Cryptopals_SET_1/*.go ./

RUN go build -o /cryptopals_set_1 main.go

WORKDIR /

COPY ./Find_the_FLAG/*.png ./

COPY ./Cryptopals_SET_1/*.txt ./

RUN mkdir output

CMD /find_flag;/cryptopals_set_1;mv decoded1.png output/

##RUN /find_flag

##RUN /cryptopals_set_1
