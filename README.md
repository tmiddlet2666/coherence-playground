# Coherence Playground

## Introduction

This repository contains my own personal Coherence related projects/demos mostly related to managing and monitoring Coherence.

> Note: These projects don't necessarily reflect the recommended 'Production' way of setting and running various Coherence
bits and pieces, but are a good starting point for playing with and trying things out with Coherence. 
> They are for educational purposes only.
> For the official way/ best practises of doing setting up Coherence, please refer to the [official Coherence commercial documentation](https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/).

## Requirements

To get the best out of these projects, you should have the following minimum requirements:

* JDK 17
* Maven 3.8.x
* Installed the latest Coherence CLI - [Install instructions](https://oracle.github.io/coherence-cli/docs/latest/#/docs/installation/01_installation)

Some examples require docker and docker-compose.

>Important: If you find any problems with these projects or see something that needs to 
> be changed or clarified, please [submit an issue](https://github.com/tmiddlet2666/coherence-playground/issues/new/choose).

## Useful Coherence Resources

* [Coherence Community](https://coherence.community/)
* [Coherence on Medium](https://medium.com/oracle-coherence)
* [Latest Coherence Commercial Documentation](https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/)

## Open Source Projects

* [Coherence CE on GitHub](https://github.com/oracle/coherence)
* [Coherence CLI](https://github.com/oracle/coherence-cli)
* [Coherence VisualVM Plugin](https://github.com/oracle/coherence-visualvm)
* [Coherence Go Client](https://github.com/oracle/coherence-go-client)
* [Coherence Demo](https://github.com/coherence-community/coherence-demo)

## Completed Projects

As I add projects I will try and organise them in some useful way, but that's not a guarantee.

### 1. Monitoring Coherence with Grafana and Prometheus in Docker Images
    
* [Starting clusters using the CLI or your own cluster](monitoring)
* [Starting cluster using Docker images](monitoring-docker)
* [Start federated clusters using `cohctl`](federation)

### 2. Capturing and searching Coherence logs in Kibana

* [Loading logs into Kibana Using Fluentd](logging)
                 
### 3. Start Federated clusters using the CLI

* [Start Active-Active Federated Coherence Clusters](federation)

## In Progress Projects
 
* [Storing HTTP sessions in Coherence with Go](go/sessions)
* [Example using Coherence Queues in Go and Hava](go/queues)
* [Setting up Coherence on a Raspberry Pi-4B](pi)

## Projects TODO/ Ideas

* [Submit an idea](https://github.com/tmiddlet2666/coherence-playground/issues/new/choose)



