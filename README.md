# ratelimiter

I've written this ratelimiter to advance my career and work with docker and kubernetes much more.

It's rate limiter gateway which can limit incoming request and it's configurable.

## How To Use?

```bash
TARGET_SERVER="http://localhost:6050" go run main.go
```

## Docker and Kubernetes

* Docker:
```bash
docker build -t ratelimiter:v1 .
```

* Kubernetes: after you configured the ratelimiter-deployment.yaml as you desire, just run below command
```bash
kubectl apply -f .
```

## Environments Variables

* TARGET_SERVER: requests should be proxied to which server.
* MAX_RATE_PER_IP: request per IP per second. default: `100`.
* TOTAL_MAX_RATE: request per second regardless of IP. default: `10000`
* PORT: server listen port. default: `9020`

## Test

I've written integeration test to ensure all things working correctly. you can run test with below command:
```bash
go test ratelimiter/tests
```
### Stress Test

I used `ab` for stress test on my gateway. first time that I tried, my app crashed and throw this error immediatly, _fatal error: concurrent map writes_. after searching about the error, I figured out that when my gateway receives two or more requests at same time, it tries to read and write on same memory address at same time and it cause the error. then I use `Mutex` that's a tool for doing stuff like that in golang, and then I fixed. next time, I tested my gateway for 1 day and it survived. 