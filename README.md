# Demo chat APP (Go + websocket)



## How to run
#### Clone repo 
```shell
 git clone https://github.com/solutionstack/jobsity-demo.git
 
```
#### Then
```shell
cd ./jobsity-demo && go mod tidy 
```

#### RUN
```shell
 go run main.go
```



#### ENV
```shell
please pass the following as --build-arg and they should hold the relevant rabbitmq/amqp connection parameters

ARG_AMQ_HOST=
ARG_AMQ_USER=
ARG_AMQ_PASS=
```

This starts the app http server on port `8000` and websocket server on port `6001`

Visit `http://localhost:8000` in your browser which should take you to the register and sign in page

On login, you'd be redirected to the chat interface

To get stock quotes, type a slash  command in the format `/stock=stockcode`

An unlimited number of users can be logged in on different browsers or tabs


### Docker build
```shell

e.g

docker build -t jobsity-demo .
docker run -p 8000:8000 -p 6001:6001 jobsity-demo

Then open localhost:8000 in your browser
```

### Libraries used
```shell
 github.com/gobwas/ws
 github.com/go-chi/chi
 github.com/rs/zerolog
 github.com/pkg/errors
 
```
