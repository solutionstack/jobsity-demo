FROM golang:1.19
ENV AMQ_HOST=woodpecker.rmq.cloudamqp.com

ENV AMQ_USER=kfimzwup
ENV AMQ_PASS=T_q0ycPW1gDy89h-l_TWyB1kNvmsLffo

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