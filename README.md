# pex-solution

Problem:
```
Please design and implement a web based API that steps through the Fibonacci sequence.

The API must expose 3 endpoints that can be called via HTTP requests:

current - returns the current number in the sequence
next - returns the next number in the sequence
previous - returns the previous number in the sequence
Example:

current -> 0
next -> 1
next -> 1
next -> 2
previous -> 1
Requirements:

The API must be able to handle high throughput (~1k requests per second).
The API should also be able to recover and restart if it unexpectedly crashes.
Assume that the API will be running on a small machine with 1 CPU and 512MB of RAM.
You may use any programming language/framework of your choice.
```

Solution:

Implemented fibonacci series using Golang and go-cache(non-persistent). Exposed 3 endpoints 
1) /current -----> to get current number in series
2) /next    ------> to get next number in series
3) /previous ------> to get previous number in series

Handled panics by recovering from panics. Used middleware to handle panics and keep application up and running during panics/crash.Also logging what caused panics.

Basically at the first , I am inserting 0 in the cache as current value. If someone wants to get previous value when current value is 0 then it will respond with status code 200 and  "no previous value found since current is 0(first element) in fibonacci series".

Use-case:

First 0 is  inserted in cache as  current value
current--->0
previous---> "no previous value found since current is 0(first element) in fibonacci series"
next----->1
next--->1
next--->2
current--->2
next------>3
previous---->2
current----->2
next---->3
current--->3
previous--->2
previous--->1
previous--->1
previous---->0
current---->0
previous---->"no previous value found since current is 0(first element) in fibonacci series"


Scenarios that were tested:

a) The API must be able to handle high throughput (~1k requests per second).

-- Used Vegeta(Vegeta is a versatile HTTP load testing tool built out of a need to drill HTTP services with a constant request rate. It can be used both as a command line utility and a library.)

Steps to test:(Tested in Macbook)
 
* First run your application. 
  a) go build  
  b) go run main.go

Open another terminal and follow below steps:

* brew install vegeta ---- for installing vegeta using brew package manager for Mac.

* go get -u github.com/tsenart/vegeta

* Since application is up and runnning we can use this command to load test:
```
echo "GET http://localhost:8443/current" | vegeta attack -duration=120s -rate=1000 | tee results.bin | vegeta report
```
This will send 1000 requests/sec to the current endpoint for 120 secs and at the end will produce result as below (extracted for local test) where 120000 requests were sent during that period(1000/sec) and as we can see all transactions passed with 200 status code which concludes application can handle load of 1000 tps

```
Requests      [total, rate, throughput]         120000, 1000.01, 999.99
Duration      [total, attack, wait]             2m0s, 2m0s, 1.784ms
Latencies     [min, mean, 50, 90, 95, 99, max]  1.22ms, 1.848ms, 1.701ms, 2.097ms, 2.442ms, 4.357ms, 79.445ms
Bytes In      [total, mean]                     2280000, 19.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:120000  
Error Set:
```

b) The API should also be able to recover and restart if it unexpectedly crashes.

* For this I am using a HandleServerCrash middleware. If any api panics of crashed due to run time error we are catching that runtime exceptions using recover. So any function or api fails it will come out and before reaching main func it will go to middleware where we are recovering from such panic and keeping server/application up and running for further incoming requests. Besides we are printing logs for the stack that it throws when it panics.

c) Assume that the API will be running on a small machine with 1 CPU and 512MB of RAM.
 
 Pre-Requisites
 -- install docker desktop

* To test this I ran K8's cluster locally using Docker Desktop where we can run kubernetes.We can try running local clusters in  minikube or docker. I used docker for this test. So Dockerfile and deployment.yaml and service.yaml are created to test this scenario.

* As it can be seen , I have allocated resources(cpu and memory) of 1 and "256Mi" respectively to run application in limited resources enviroment. 

Steps to run program in k8's using docker are as follows assuming you are on root directory of project folder:
a) docker build --tag trial-pex .  (it creates docker image with tag trial-pex which can be seen in docker dashboard)
b) kubectl apply -f deployment.yaml
c) kubectl apply -f service.yaml
d) kubectl port-forward deployment/pex-solution 8443 







