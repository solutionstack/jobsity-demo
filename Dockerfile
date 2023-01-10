FROM amd64/golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.* ./

RUN go build -o ./jobsity-demo

RUN chmod +x ./jobsity-demo
EXPOSE 8000

CMD [ "./jobsity-demo" ]