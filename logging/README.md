# Loading logs into Kibana Using Fluentd

## What You Will Build

This example shows how to set up Fluentd and Kibana to scrape logs from a Coherence cluster 
started with the [Coherence CLI](https://github.com/oracle/coherence-cli).

The example startup Kibana 7.13.1 and Elastic Search:7.13.1images using Docker and then startup a cluster using 
`cohctl`.

Please refer to the *official* Coherence Operator docs on logging [here](https://oracle.github.io/coherence-operator/docs/latest/#/docs/logging/010_overview).

## What You Need

1. JDK 17+
2. Docker and docker-compose
3. Maven 3.8.x
4. The Coherence CLI installed and on the PATH
5. Clone of this repository using `git clone https://github.com/tmiddlet2666/coherence-playground.git`
           
> Note: Make sure java and mvn are in your PATH.

Open a terminal and change to the directory `logging`.

## Setup fluentd image
 
### Edit Fluentd config

You must edit the file `fluentd/conf/fluentd.conf` and ensure the following entry matches the cluster
name you are creating below. E.g. in this case it is `my-cluster`.

```bash
  path /logs/my-cluster/storage-*.log
```
        
### Update the docker-compose.yaml

You must edit the file `docker-compose.yaml` and change the following line to match your
home directory. E.g. replace `timmiddleton` with your username.

```bash
      - /Users/timmiddleton/.cohctl/logs:/logs
```

You must also change the TZ below to match your own timezone so that Kibana will show data correctly.

```bash
       - TZ=Australia/Perth
```

### Build the image

Run the following to build the fluentd_logging image:

```bash
docker build -t coherence/fluentd_logging:latest fluentd
```

## Running the Example
     
1. Create and start a Coherence Cluster using the following:

   ```bash
   $ cohctl create cluster my-cluster

   Cluster name:         my-cluster
   Cluster version:      22.09
   Cluster port:         7574
   Management port:      30000
   Replica count:        3
   Initial memory:       512m
   Persistence mode:     on-demand
   Group ID:             com.oracle.coherence.ce
   Additional artifacts: 
   Startup Profile:      
   Dependency Tool:      mvn
   Are you sure you want to create the cluster with the above details? (y/n) y

   Checking 3 Maven dependencies...
    - com.oracle.coherence.ce:coherence:22.09
    - com.oracle.coherence.ce:coherence-json:22.09
    - org.jline:jline:3.20.0
   Starting 3 cluster members for cluster my-cluster
   Starting cluster member storage-0...
   Starting cluster member storage-1...
   Starting cluster member storage-2...
   Cluster added and started
   Current context is now my-cluster
   ```    

2. Startup Grafana and Prometheus

   Ensure you are in the `logging` directory, have Docker running and then issue the following:

   ```bash
   $ docker-compose up -d  
   
   [+] Running 4/4
   Network logging_default      Created                                                                                                                                                                        0.0s
   Container elasticsearch      Started                                                                                                                                                                        0.6s
   Container logging-fluentd-1  Started                                                                                                                                                                        1.5s
   Container logging-kibana-1   Started  
   
   $ docker ps
   
   CONTAINER ID   IMAGE                                                  COMMAND                  CREATED          STATUS          PORTS                                                                                                    NAMES
   11b0180ae4e5   docker.elastic.co/kibana/kibana:7.13.1                 "/bin/tini -- /usr/l…"   27 seconds ago   Up 23 seconds   0.0.0.0:5601->5601/tcp, :::5601->5601/tcp                                                                logging-kibana-1
   b686efc1b65d   coherence/fluentd_logging:latest                       "tini -- /bin/entryp…"   27 seconds ago   Up 23 seconds   5140/tcp, 0.0.0.0:24224->24224/tcp, 0.0.0.0:24224->24224/udp, :::24224->24224/tcp, :::24224->24224/udp   logging-fluentd-1
   a88fbd2c2497   docker.elastic.co/elasticsearch/elasticsearch:7.13.1   "/bin/tini -- /usr/l…"   27 seconds ago   Up 24 seconds   0.0.0.0:9200->9200/tcp, :::9200->9200/tcp, 9300/tcp                                                      elasticsearch
   ```

3. Check Kibana
   
    Open Kibana dashboard at http://127.0.0.1:5601/app/management/kibana/indexPatterns
         
    Click on `Create Index Pattern`, and enter `coherence-cluster*` in the index pattern name, and you should see a pattern 
    below called `coherence-cluster*`.  Click `Next Step` and choose `@timestamp` as your time field. 
    Click `Create Index Pattern`.

4. View logs in Kibana.

   Open Kibana at http://127.0.0.1:5601/app/discover and search for data.

   Refer to the following regarding searching for logs: https://www.elastic.co/guide/en/enterprise-search/current/logging-view-query-logs.html

5. Import dashboards

   If you would like import the out-of-the-box dashboards from the Coherence Operator, please do the following:

   1. Clone the Coherence Operator repository - https://github.com/oracle/coherence-operator.git
   2. Access: http://127.0.0.1:5601/app/management/kibana/objects and click on `Import`
   3. Select the file `kibana-dashboard-data.ndjson` from `coherence-operator/dashboards/kibana` directory and import
   4. Access the dashboards via http://127.0.0.1:5601/app/dashboards


## Shutting everything down

1. Stop all Docker processes using `docker-compose down`
2. Stop the Coherence cluster using `cohctl stop cluster my-cluster -y`
