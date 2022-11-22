# Monitoring Coherence under WebLogic Server

## What You Will Build

This document explains how to monitor Coherence under WebLogic Server, using the demo application, and view metrics via Grafana using `docker-compose`.

In this example we initially deploy the WebLogic Server examples and startup the following docker images
1. Grafana - 8.5.6
1. Prometheus - v2.36.2

The Docker images expose the following ports:

* Grafana
  * 3000 - Grafana UI
* Prometheus
  * 9090 - Prometheus UI

> Note: As of writing this, the most stable release of Grafana is 8.5.6. It is not recommended to use the 9.X versions yet as they have many bugs still.

>You can edit `docker-compose.yaml` to change the Grafana, Prometheus or Coherence CE version.

## What You Need

You must have the following:
* Docker Desktop for Mac or the equivalent Docker environment for your O/S
* docker-compose - at least version 1.20+
* Clone of this repository using git clone https://github.com/tmiddlet2666/coherence-playground.git
* WebLogic Server 14.1.1.0 Downloaded
                                 
## Setup

### Download and Install WebLogic Server

Download WebLogic Server from https://www.oracle.com/in/middleware/technologies/fusionmiddleware-downloads.html and make sure
you install with the examples.  You do not have to start the example domain at the end.

### Install the Coherence Demo Domain

1. After install, change to the directory `WLS_INSTALL/wlserver/server/bin`

2. Run `. ./setWLSEnv.sh` for Linux/Mac or `setWLSEnv.cmd` for Windows

3. Run the following command to install and startup the demo

   ```bash    
   cd $WLS_INSTALL/wlserver/samples/server/examples/src/examples/coherence/managed-coherence-servers/multi-server
   
   ant -Ddo.delete=y -Dsingle.server.password=welcome1 -Dsingle.server.host=127.0.0.1 -Dsingle.server.port=7001 -Dmulti.server.port=7003 deploy
   ```                                                                                                                                         

   Once this is completed, you will have a 6 node WebLogic Server cluster with 4 storage nodes and 2 client nodes
   the Coherence examples installed.
   
   Default username/password is `weblogic/welcome1`

5. Install the `coherence-metrics-resource.war` to the `ClientTier` and `StorageTier`.
 
   Refer to the [Coherence Documentation for more information](https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/manage/using-coherence-metrics.html#GUID-76E06390-E71B-4AC2-B1BA-DFB11FD2BDCC).
### Download the Grafana Dashboards

#### Linx/MacOS

Run the following to download the dashboards:

```bash
cd coherence-docker-metrics/grafana
./download-dashboards.sh
```
#### Windows

1. Clone the Coherence Operator repository - https://github.com/oracle/coherence-operator.git

1. Copy all the `*.json` files from `coherence-operator/dashboards/grafana` to the directory `grafana/dashboards` under cloned `coherence-docker-metrics` repository.
   
## Running the Example

from the `monitoring-docker` directory run:

```bash
docker-compose up -d
Creating network "coherence-docker-metrics_coherence" with the default driver
Creating coherence-docker-metrics_prometheus_1 ... done
Creating coherence-docker-metrics_grafana_1    ... done
Creating coherence-docker-metrics_coherence1_1 ... done
Creating coherence-docker-metrics_coherence2_1 ... done

$ docker ps
CONTAINER ID   IMAGE                               COMMAND                  CREATED          STATUS                             PORTS                                                                                                                     NAMES
dd27c452d0f6   grafana/grafana:8.5.6               "/run.sh"                13 seconds ago   Up 11 seconds                      0.0.0.0:3000->3000/tcp                                                                                                    coherence-docker-metrics_grafana_1
cff1586818ba   ghcr.io/oracle/coherence-ce:22.09   "java -cp /coherence…"   13 seconds ago   Up 11 seconds (health: starting)   1408/tcp, 6676/tcp, 20000-20001/tcp, 30000/tcp, 0.0.0.0:9613->9612/tcp                                                    coherence-docker-metrics_coherence2_1
fc9f62a6be4e   prom/prometheus:v2.36.2             "/bin/prometheus --c…"   13 seconds ago   Up 11 seconds                      0.0.0.0:9090->9090/tcp                                                                                                    coherence-docker-metrics_prometheus_1
4e04f10ac3af   ghcr.io/oracle/coherence-ce:22.09   "java -cp /coherence…"   13 seconds ago   Up 11 seconds (health: starting)   0.0.0.0:1408->1408/tcp, 0.0.0.0:9612->9612/tcp, 6676/tcp, 0.0.0.0:20000->20000/tcp, 0.0.0.0:30000->30000/tcp, 20001/tcp   coherence-docker-metrics_coherence1_1
```

### Access Grafana Dashboards

Access the Grafana dashboard using the following URL:

http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main

Default username/password is admin/admin

### Access Prometheus

If required, you can access the Prometheus UI using the following URL:

http://127.0.0.1:9000/

## Troubleshooting

If you have issues with the Grafana dashboards showing data, then check the Prometheus targets using the following URL:

http://127.0.0.1:9090/targets

If the connection is not working, then you may need to change the hostname in [prometheus/prometheus.yaml](prometheus/prometheus.yaml)
to the IP/hostname you are running on.

The current hostname is docker0, IP of 172.17.0.1, which works for most environments. You can also use `host.docker.internal` on Mac.

> Note: You may also need to change the hostname in [grafana/datasources.yaml](grafana/datasources.yaml).
      
## Shutting everything Down

1. Stop all Docker processes using `docker-compose down`
 
2. Shutdown the WebLogic Instance
         
   ```bash
   ant -Ddo.delete=y -Dsingle.server.password=welcome1 -Dsingle.server.host=127.0.0.1 -Dsingle.server.port=7001 -Dmulti.server.port=7003 shutdownMultiServer
   ```   