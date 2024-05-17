# Setting up Coherence on a Raspberry Pi-4B

## What You Will Build

This example shows how to set up Coherence CE 22.06.8 on a Raspberry Pi-4B using the Coherence CLI and 
setting it to auto-start on boot.

I was helping a colleague of mine setup Coherence on a Pi for storing state for Python programs 
and thought this may be a good alternative from a docker image.

> Note 1. The Python client requires Coherence gRPC proxy dependencies to be able to start gGRPC proxy for the client to connect. 
         
> Note 2. This example just starts a single Coherence cache server to 
> minimize memory overhead, but you could start multiple to allow for cluster HA. See the end of the readme for more details.

## What You Need

* A Raspberry Pi setup and connected to the Internet, that's all!

## Install required Linux software
 
1. Login as `pi` and run the following to install Java 17 and Maven.

   ```bash
   sudo apt install openjdk-17-jdk maven
   ```
   
   ```bash
   mvn  --version
   ````       
   
   Output:
   ```bash
   Apache Maven 3.8.7
   Maven home: /usr/share/maven
   Java version: 17.0.10, vendor: Debian, runtime: /usr/lib/jvm/java-17-openjdk-arm64
   Default locale: en_GB, platform encoding: UTF-8
   OS name: "linux", version: "6.6.20+rpt-rpi-v8", arch: "aarch64", family: "unix"
   ```
   
2. Install the Coherence CLI

   ```bash
   curl -sL https://raw.githubusercontent.com/oracle/coherence-cli/main/scripts/install.sh | bash
   ```
      
   Output:
   ```bash
   Installing Coherence CLI 1.6.1 for Linux/aarch64 into /usr/local/bin ...
   Using 'sudo' to mv cohctl binary to /usr/local/bin

   To uninstall the Coherence CLI execute the following:
     sudo rm /usr/local/bin/cohctl
   ```
  
   ```bash
   cohctl version | bash
   ``` 
   
   ```bash
   pi@pi-4b:~ $ cohctl version
   Coherence Command Line Interface
   CLI Version:  1.6.1
   Date:         2024-04-17T03:46:43Z
   Commit:       7e869a55e06476de92940b0f8b2887100d7921b0
   OS:           linux
   Architecture: arm64
   Go Version:   go1.20.14
   ```
    
3. Create a Profile to set the gRPC port

   ```bash
   cohctl set profile grpc -v "-Dcoherence.grpc.server.port=1408" -y
   cohctl get profiles
   ``` 
   
   Output:
   ```bash
   PROFILE  VALUE                            
   grpc     -Dcoherence.grpc.server.port=1408
   ```   
   
4. Create a minimal cache config to reduce what is started
   
   Create the file `/home/pi/pi-cache-config.xml` with the following contents:

   ```xml
   <?xml version="1.0"?>
   <!--
    Minimal cache config.
   -->
   <cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                 xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
                 xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd">
   
     <caching-scheme-mapping>
       <cache-mapping>
         <cache-name>*</cache-name>
         <scheme-name>server</scheme-name>
       </cache-mapping>
     </caching-scheme-mapping>
   
     <caching-schemes>
       <distributed-scheme>
         <scheme-name>server</scheme-name>
         <service-name>PartitionedCache</service-name>
         <backing-map-scheme>
           <local-scheme>
             <unit-calculator>BINARY</unit-calculator>
           </local-scheme>
         </backing-map-scheme>
         <autostart>true</autostart>
       </distributed-scheme>
     </caching-schemes>
   </cache-config>
   ```

5. Create the Cluster

   For this simple cluster we are creating and starting only 1 member. You can adjust the memory if you require more.

   ```bash
   cohctl create cluster local -v 22.06.8 -P grpc -r 1 -M 256m -a coherence-grpc-proxy --cache-config /home/pi/pi-cache-config.xml
   ```      
   
   > Note: The required Maven dependencies will be downloaded and the cluster started. This may take a short while.
   
   Output:
   ```bash
   Checking 4 Maven dependencies...
   - com.oracle.coherence.ce:coherence:22.06.8
   - com.oracle.coherence.ce:coherence-grpc-proxy:22.06.8
   - com.oracle.coherence.ce:coherence-json:22.06.8
   - org.jline:jline:3.25.0
   Starting 1 cluster members for cluster local
   Starting cluster member storage-0...
   Current context is now local
   Cluster added and started
   ```
   
   After a 5-10 seconds, check the status of the cluster using:

   ```bash
   cohctl get members
   ```
   
   Output: 
   ```bash
   Using cluster connection 'local' from current context.
   
   Total cluster members: 1
   Storage enabled count: 1
   Departure count:       0
   
   Cluster Heap - Total: 256 MB Used: 29 MB Available: 227 MB (88.7%)
   Storage Heap - Total: 256 MB Used: 29 MB Available: 227 MB (88.7%)
   
   NODE ID  ADDRESS     PORT   PROCESS  MEMBER     ROLE             STORAGE  MAX HEAP  USED HEAP  AVAIL HEAP
         1  /127.0.0.1  38975     2377  storage-0  CoherenceServer  true       256 MB      29 MB      227 MB
   ```
   
6. Setting the cluster to start on boot

   As the `pi` user, not sudo, issue the following:

   ```bash
   crontab -e
   ```                                                                   
   
   Add the following to the end of file and save:
                                                 
   ```bash
   @reboot /usr/local/bin/cohctl start cluster local -r 1 -M 256m -P grpc
   ```    
   
   Restart you Pi by using the following:
   ```bash
   sudo shutdown -r now
   ```
   
   Once the Pi restarts, login again the use following command to verify the cluster is running:
   ```bash
   cohctl get members
   ```
   
   Output: 
   ```bash
   Using cluster connection 'local' from current context.
   
   Total cluster members: 1
   Storage enabled count: 1
   Departure count:       0
   
   Cluster Heap - Total: 256 MB Used: 29 MB Available: 227 MB (88.7%)
   Storage Heap - Total: 256 MB Used: 29 MB Available: 227 MB (88.7%)
   
   NODE ID  ADDRESS     PORT   PROCESS  MEMBER     ROLE             STORAGE  MAX HEAP  USED HEAP  AVAIL HEAP
         1  /127.0.0.1  38975     2377  storage-0  CoherenceServer  true       256 MB      29 MB      227 MB
   ```  
   
7. Log file locations

   Issue the following if you wish to see the log files.

   ```bash
   ls -l ~/.cohctl/logs/local
   ```  
   
   Output: 
   ```bash
   total 16
   -rw-r--r-- 1 pi pi 13060 Apr 18 16:26 storage-0.log
   ```

## Stopping, restarting and scaling the cluster
    
   You can issue the following to stop the cluster:

   ```bash
   cohctl stop cluster local -y
   ```

   You can issue the following to restart the cluster:
     
   ```bash
   cohctl start cluster local -r 1 -M 256m -P grpc
   ```

   You can issue the following to scale the cluster to 2 nodes. (We do not use the `grpc` profile as we only have one member exposing this port)
     
   ```bash
   cohctl scale cluster local -r 2 -M 256m
   ```

   After a shot time you can issue the following command to view the members:

   ```bash
   cohctl get members
   ```
   
   Output:
   ```bash
   Using cluster connection 'local' from current context.
   
   Total cluster members: 2
   Storage enabled count: 2
   Departure count:       0
   
   Cluster Heap - Total: 512 MB Used: 88 MB Available: 424 MB (82.8%)
   Storage Heap - Total: 512 MB Used: 88 MB Available: 424 MB (82.8%)
   
   NODE ID  ADDRESS     PORT   PROCESS  MEMBER     ROLE             STORAGE  MAX HEAP  USED HEAP  AVAIL HEAP
         1  /127.0.0.1  39533      739  storage-0  CoherenceServer  true       256 MB      72 MB      184 MB
         2  /127.0.0.1  32835     1221  storage-1  CoherenceServer  true       256 MB      16 MB      240 MB
   ```

