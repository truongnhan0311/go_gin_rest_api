export PATH=$(go env GOPATH)/bin:$PATH
wrk -c 500 -d 1m -t 5 http://127.0.0.1:8080/hello
