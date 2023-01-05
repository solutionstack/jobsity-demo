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
 go run cmd/cmd.go

```
This starts the app http server on port `8000` and websocket server on port `6001`

Visit `http://localhost:8000` in your browser which should take you to the register and sign in page

On login, you'd be redirected to the chat interface

An unlimited number of users can be logged in on different browser

### Libraries used
```shell
 github.com/gobwas/ws
 github.com/go-chi/chi
 github.com/rs/zerolog
 github.com/pkg/errors
 
```