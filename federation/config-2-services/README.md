# Federation of 2 services

## Overview

We have 2 clusters, one in Perth and one in Sydney

Perth

* Cluster Port: 7574
* FederatedCache1 service - port 40000
* FederatedCache2 service - port 40001

Sydney

* Cluster Port: 7574
* FederatedCache1 service - port 50000
* FederatedCache2 service - port 50001


System properties for cache config
* `coherence.local.federation.port1` = local port for service 1
* `coherence.local.federation.port2` = local port for service 2

System properties for override:

* `coherence.primary.federation.port1` = primary cluster port 1
* `coherence.primary.federation.port2` = primary cluster port 2
* `coherence.secondary.federation.port1` = secondly cluster port 1
* `coherence.secondary.federation.port2` = secondly cluster port 2

## Create the following profiles

Profile for `perth`

```bash
cohctl set profile Perth -v "-Dcoherence.local.federation.port1=40000 -Dcoherence.local.federation.port2=40001 -Dcoherence.primary.federation.port1=40000 -Dcoherence.primary.federation.port2=40001 -Dcoherence.secondary.federation.port1=50000 -Dcoherence.secondary.federation.port2=50001"
```

```bash
cohctl set profile Sydney -v "-Dcoherence.local.federation.port1=50000 -Dcoherence.local.federation.port2=50001 -Dcoherence.primary.federation.port1=40000 -Dcoherence.primary.federation.port2=40001 -Dcoherence.secondary.federation.port1=50000 -Dcoherence.secondary.federation.port2=50001"
```

We are starting the cluster nodes with `-r 1` to only start up a single server otherwise the ports start auto-adjusting.

## Perth Cluster

```bash
cohctl create cluster Perth -C -v 14.1.2-0-0 -H 30000 -p 7574 -M 1g -P Perth -r 1 \
--cache-config /Users/timbo/Documents/CoherenceEngineering/github/coherence-playground/federation/config-2-services/federated-cache-config.xml \
--override-config /Users/timbo/Documents/CoherenceEngineering/github/coherence-playground/federation/config-2-services/federated-override.xml
```

If you want to restart:

```bash
cohctl start cluster Perth -P Perth -r 1 -M 2g
```

## Sydney Cluster

```bash
cohctl create cluster Sydney -C -v 14.1.2-0-0 -H 30001 -p 7575 -M 1g -P Sydney  -r 1 \
--cache-config /Users/timbo/Documents/CoherenceEngineering/github/coherence-playground/federation/config-2-services/federated-cache-config.xml \
--override-config /Users/timbo/Documents/CoherenceEngineering/github/coherence-playground/federation/config-2-services/federated-override.xml
```

If you want to restart:

```bash
cohctl start cluster Sydney -P Sydney -r 1 -M 2g
```


