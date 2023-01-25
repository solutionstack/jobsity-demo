FROM golang:1.19

ARG ARG_AMQ_HOST=

ARG ARG_AMQ_USER=
ARG ARG_AMQ_PASS=

ENV AMQ_HOST=$ARG_AMQ_HOST

ENV AMQ_USER=$ARG_AMQ_USER
ENV AMQ_PASS=$ARG_AMQ_PASS

EXPOSE 8000
EXPOSE 6001

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN go build -o ./jobsity-demo .

RUN chmod +x ./jobsity-demo


CMD [ "./jobsity-demo" ]
