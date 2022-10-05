# Monitoring Coherence with Grafana and Prometheus using Docker images

## What You Will Build

This example shows how to set up Grafana and Prometheus to monitor a Coherence cluster 
started with the Coherence CLI.

The example startup Grafana 8.5.6 and Prometheus v2.36.2 images using Docker and then startup a cluster using 
`cohctl` with Coherence metrics enabled.

## What You Need

1. JDK 17+
2. Docker and docker-compose
3. Maven 3.8.x
4. The Coherence CLI installed and in the PATH
5. Clone of this repository using `git clone https://github.com/tmiddlet2666/coherence-playground.git`
      
Open a terminal and change to the directory `monitoring`

## Setup

You will need to download the latest Grafana dashboards that are available in the Coherence Operator.

### Linux/ OSX

1. Ensure you are in the `monitoring/grafana` directory 
2. Run the following to download the dashboards.

    ```bash
    $ ./download-dashboards.sh
    
    alerts-dashboard.json
    coherence-dashboard-main.json
    federation-details-dashboard.json
    http-servers-summary-dashboard.json
    member-details-dashboard.json
    proxy-server-detail-dashboard.json
    services-summary-dashboard.json
    cache-details-dashboard.json
    federation-summary-dashboard.json
    kubernetes-summary-dashboard.json
    members-summary-dashboard.json
    proxy-servers-summary-dashboard.json
    caches-summary-dashboard.json
    elastic-data-summary-dashboard.json
    grpc-summary-dashboard.json
    machines-summary-dashboard.json
    persistence-summary-dashboard.json
    service-details-dashboard.json
    coherence-executors-summary.json
    coherence-executor-detail.json
    ```

### Windows
 
1. Clone the Coherence Operator repository - https://github.com/oracle/coherence-operator.git

1. Copy all the `*.json` files from `coherence-operator/dashboards/grafana` to the directory `monitoring/grafana/dashboards` directory.

## Running the Example
     
1. Create and start a Coherence Cluster with metrics enabled using the following:

   ```bash
   $ cohctl create cluster my-cluster -t 9612

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
   
   $ cohctl set context my-cluster
   Current context is now my-cluster
   ```   

2. Set the current context to `my-cluster`

   ```bash
   $ cohctl set context my-cluster
   Current context is now my-cluster
      
3. Startup Grafana and Prometheus

    Ensure you are in the `monitoring` directory and run:

    ```bash
    $ docker-compose up -d  
   
   [+] Running 2/2
   ⠿ Container grafana-prometheus-1  Running                                                                                                                                                                               0.0s
   ⠿ Container grafana-grafana-1     Started  
   
   $ docker ps
   
   CONTAINER ID   IMAGE                     COMMAND                  CREATED         STATUS          PORTS                                       NAMES
   4bc1250ec605   prom/prometheus:v2.36.2   "/bin/prometheus --c…"   2 minutes ago   Up 2 minutes    0.0.0.0:9090->9090/tcp, :::9090->9090/tcp   grafana-prometheus-1
   844114ca12a8   grafana/grafana:8.5.6     "/run.sh"                2 minutes ago   Up 22 seconds   0.0.0.0:3000->3000/tcp, :::3000->3000/tcp   grafana-grafana-
   ```
   
4. Check the status of Prometheus targets

   Open http://localhost:9090/targets and check that at least 3 of the targets are UP. 

   It may take a minute for them to be discovered.

5. Check Grafana
   
   Open the main Grafana dashboard at http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main.
           
   The default username and password is: `admin`. You can change this or just press `Skip` on the first login.

   You should now see the main dashboard similar to the following:

   ![Coherence Demo](assets/coherence-dashboard-main.png "Coherence Dashboard Main")
   
6. Add some data using the console

   ```bash
   $ cohctl start console
   ``` 
   
   At the `Map (?):` prompt type `cache test` and press return.

   Next, type `bulkput 10000 100 0 1000` to add 10,000 objects of size 100, starting at 0 in batches of 1000.
       
7. Open the `Caches Summary Dashboard` by clicking on the `Available Dashboards` link on the top right.
 
8. Scale your cluster to 6 nodes using `cohctl scale cluster -r 6` and observe the Grafana dashboards updating.

## Shutting everything down

1. Stop all Docker processes using `docker-compose down`
2. Stop the Coherence cluster using `cohctl stop cluster my-cluster -y`
