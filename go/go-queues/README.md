# Using Coherence Queues in Go

This example is the code from a blog post (TBC) that talks about how to use and
access Coherence Queues from  Go.

## Prerequisites

1. Docker setup using either Docker or Rancher Desktop
2. Go 1.20 or above 

## Run the Demo

### Start a Coherence cluster

```bash
docker run -d -p 1408:1408 -p 30000:30000 ghcr.io/oracle/coherence-ce:24.03
```

### Start one or more receivers

```bash
go run receiver/main.go 
2024/05/16 14:02:54 session: af36618b-d8fe-493c-ad39-4a442aa30b8b connected to address localhost:1408
2024/05/16 14:02:54 Waiting for orders
```

### Start one or more publisher

Provide the start order number and number of orders to publish. 
For example to start a order 1 and publish 100 orders use the following:

```bash
go run publisher/main.go 1 100
```
