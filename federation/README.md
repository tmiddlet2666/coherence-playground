# Start Active-Active Coherence Federated clusters using the CLI

## What You Will Build

This example shows how to set up an Active-Active Coherence Federated clusters using Coherence Grid Edition
and started with the [Coherence CLI](https://github.com/oracle/coherence-cli).

The example comprises two Coherence clusters with active-active Federation between 
ClusterA and ClusterB on the same machine.

> Note: Federation is supported only with Coherence Grid Edition, and will not work with CE. 

## What You Need

1. JDK 17+
2. Maven 3.8.x
3. The Coherence CLI installed and on the PATH
4. Clone of this repository using `git clone https://github.com/tmiddlet2666/coherence-playground.git`
5. An installation of Coherence Grid Edition 14.1.1.0.X or 14.1.1.2206.x
           
> Note: Make sure java and mvn are in your PATH.

Open a terminal and change to the directory `federation`.

## Setup

### Download Oracle Coherence 14.1.1.0.0 or later.

You can download Coherence from http://www.oracle.com/technetwork/middleware/coherence/downloads/index.html.

You should also apply the latest patch available for your Coherence version.

### Set the Environment Variables

Ensure that the following environment variables are set in your configuration:

* `JAVA_HOME` -- This variable must point to the location of the JDK version supported by the Oracle Coherence version that you use. Ensure that the path is set accordingly:</br>

   For Linux/UNIX OS:

   ```bash
   export PATH=$JAVA_HOME/bin:$PATH
   ```

   For Windows OS:

   ```bash
   set PATH=%JAVA_HOME%\bin;%PATH%
   ```

* `COHERENCE_HOME` -- This variable must point to the `\coherence` directory of your Coherence installation. This is required for the Maven `install-file` commands.

* `MAVEN_HOME` -- If `mvn` command is not set in your PATH variable, then set `MAVEN_HOME` directed to the `bin` folder of Maven installation and then add `MAVEN_HOME\bin` to your PATH variable list.

### Install Coherence JARs into your Maven repository

Install Coherence and Coherence HTTP Netty installed into your local maven repository.

For Linux/UNIX/Mac OS:

```bash
mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence.jar -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence/14.1.1/coherence.14.1.1.pom
mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence-json.jar -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence-json/14.1.1/coherence-json.14.1.1.pom
```

For Windows OS:

```bash
mvn install:install-file -Dfile=%COHERENCE_HOME%\lib\coherence.jar -DpomFile=%COHERENCE_HOME%\plugins\maven\com\oracle\coherence\coherence\14.1.1\coherence.14.1.1.pom
mvn install:install-file -Dfile=%COHERENCE_HOME%\lib\coherence-json.jar -DpomFile=%COHERENCE_HOME%\plugins\maven\com\oracle\coherence\coherence-json\14.1.1\coherence-json.14.1.1.pom
```

## <a name="run"></a> Running the Example

1. Create a profile called `federation` which points to the absolute path of your cache config and override:

   E.g. if cloned dir is: /Users/timmiddleton/coherence-playground

   ```bash
   cohctl set profile federation -v "-Dcoherence.cacheconfig=/Users/timmiddleton/coherence-playground/federation/config-2-servers/federated-cache-config.xml -Dcoherence.override=/Users/timmiddleton/coherence-playground/federation/config-2-servers/federated-override.xml"
   ```

2. Create and start ClusterA using the following:

   ```bash
   $ cohctl create cluster ClusterA -C -v 14.1.1-2206-3 -P federation -H 30000 -p 7574 -M 2g -t 9612

   Cluster name:         ClusterA
   Cluster version:      14.1.1-2206-3
   Cluster port:         7574
   Management port:      30000
   Replica count:        3
   Initial memory:       512m
   Persistence mode:     on-demand
   Group ID:             com.oracle.coherence
   Additional artifacts: 
   Startup Profile:      federation
   Log destination root: 
   Dependency tool:      mvn
   Are you sure you want to create the cluster with the above details? (y/n) y

   Checking 3 Maven dependencies...
   - com.oracle.coherence:coherence:14.1.1-2206-3
   - com.oracle.coherence:coherence-json:14.1.1-2206-3
   - org.jline:jline:3.20.0
   Starting 3 cluster members for cluster ClusterA
   Starting cluster member storage-0...
   Starting cluster member storage-1...
   Starting cluster member storage-2...
   Current context is now ClusterA
   ```    

3. Create and start ClusterB using the following:

   ```bash
   $ cohctl create cluster ClusterB -C -v 14.1.1-2206-3 -P federation -H 30001 -p 7575 -M 2g -t 9615

   Cluster name:         ClusterB
   Cluster version:      14.1.1-2206-3
   Cluster port:         7575
   Management port:      30001
   Replica count:        3
   Initial memory:       512m
   Persistence mode:     on-demand
   Group ID:             com.oracle.coherence
   Additional artifacts: 
   Startup Profile:      federation
   Log destination root: 
   Dependency tool:      mvn
   Are you sure you want to create the cluster with the above details? (y/n) y

   Checking 3 Maven dependencies...
   - com.oracle.coherence:coherence:14.1.1-2206-3
   - com.oracle.coherence:coherence-json:14.1.1-2206-3
   - org.jline:jline:3.20.0
   Starting 3 cluster members for cluster ClusterB
   Starting cluster member storage-0...
   Starting cluster member storage-1...
   Starting cluster member storage-2...
   Current context is now ClusterB
   ```    

4. Add data to ClusterA

    ```bash
    cohctl set context ClusterA
   
    cohctl start console -P federation   
   
    Map (?): cache test
           
    Map (test): bulkput 100000 100 0 100      

    ```

5. Query data in ClusterB
   
    Open a new terminal and connect to ClusterB. (Make sure you have JDK17 and CLI in the PATH)

    ```bash
    cohctl set context ClusterB
   
    cohctl start console -P federation   
   
    Map (?): cache test
           
    Map (test): size
    100000
    ```
    
    You can see the data has been replicated:

## Restarting the clusters
       
```bash
cohctl start cluster ClusterA -P federation -t 9612 -M 2g

cohctl start cluster ClusterB -P federation -t 9615 -M 2g 
```


## Shutting everything down

1. Stop ClusterA using `cohctl stop cluster ClusterA -y`
1. Stop ClusterB using `cohctl stop cluster ClusterB -y`
