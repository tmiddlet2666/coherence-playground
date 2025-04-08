# Simple Federation example in Coherence Operator

This example shows running Coherence Federation between two clusters in the same Kubernetes clusters.

> Note: This is an example only and normally Federation is used to federate between separate geographic locations.

## Pre-Requisites 

You must have access to a Kubernetes cluster already setup.

### Obtain your container registry Auth token

In this example we are using the Coherence Grid Edition container from Oracle's container registry.

* `container-registry.oracle.com/middleware/coherence:14.1.2.0.0`

To be able to pull the above image you need to do the following:

1. Sign in to the Container Registry at `https://container-registry.oracle.com/`. (If you don't have an account, with Oracle you will need to create one)
2. Once singed in, search for `coherence` and select the link for `Oracle Coherence`
3. If you have not already, you will need to accept the 'Oracle Standard Terms and Conditions' on the right before you are able to pull the image
4. Once you have accepted the terms and condition click on the drop-down next to your name on the top right and select `Auth Token`.
5. Click on `Generate Token` and save this in a secure place for use further down.

NOTE: See [The OCI Documentation](https://docs.oracle.com/en-us/iaas/Content/Registry/Tasks/registrygettingauthtoken.htm) for more information on creating your Auth token.
 
### Install the Coherence Operator

```bash
kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml
```

### Create the Namespace

```bash
kubectl create namespace federation-demo
```

### Create Image Pull Secret

Replace username with your username and password with your auth-token.

```bash
kubectl create secret docker-registry ocr-pull-secret \
--docker-server=container-registry.oracle.com \
--docker-username="<username>" --docker-password="<password>" \
--docker-email="<email>" -n federation-demo
```

### Create a config map with the Config

Run the following command to create the config map named `storage-config` to store the Coherence configuration files.

```bash
kubectl delete secret storage-config -n federation-demo || true
kubectl create secret generic storage-config -n federation-demo \
--from-file=config/federated-cache-config.xml \
--from-file=config/federated-override.xml
```

## Deploy the Example

### Deploy the primary-cluster

```bash
kubectl apply -n federation-demo -f yaml/primary-cluster.yaml
```

### Deploy the secondary-cluster

```bash
kubectl apply -n federation-demo -f yaml/secondary-cluster.yaml
```

## Exercise the Example

1. View federation status on primary-cluster

   ```bash
   kubectl exec -it -n federation-demo primary-cluster-0 -- /coherence-operator/utils/cohctl get federation all -W -o wide
   ```

2. Add data to the primary cluster

   ```bash
   kubectl exec -it -n federation-demo primary-cluster-0 -- /coherence-operator/utils/runner console
   ```

   ```bash
   cache test
   bulkput 10000 100 0 10
   size
   ```
   
3. Validate the primary cluster caches

   You should see the messages send being updated.

   Get the cache size from the primary cluster.

   ```bash
   kubectl exec -it -n federation-demo primary-cluster-0 -- /coherence-operator/utils/cohctl get caches -o wide -I
   ```

    You should see cache count is 10,000.   

   ```bash
   SERVICE            CACHE              COUNT   SIZE  AVG SIZE    PUTS  GETS  REMOVES  EVICTIONS  HITS   MISSES  HIT PROB
   FederatedCache     test              10,000  11 MB     1,160  10,000     0        0          0     0        0     0.00%
   ```
 
4. Validate the secondary cluster caches

   Get the cache size from the secondary cluster.

   ```bash
   kubectl exec -it -n federation-demo secondary-cluster-0 -- /coherence-operator/utils/cohctl get caches -o wide -I
   ```

    You should see cache count is 10,000 in both clusters.  

   ```bash
   SERVICE            CACHE              COUNT   SIZE  AVG SIZE    PUTS  GETS  REMOVES  EVICTIONS  HITS   MISSES  HIT PROB
   FederatedCache     test              10,000  11 MB     1,160  10,000     0        0          0     0        0     0.00%
   ```

5. Add data to the secondary cluster

   ```bash
   kubectl exec -it -n federation-demo secondary-cluster-0 -- /coherence-operator/utils/runner console
   ```

   ```bash
   cache test
   size
   ```
   This will return: 10,000, Add another 5000 entries, starting at key 15000

   ```bash
   bulkput 5000 100 15000 10
   size
   ```

   There will be 15,000 entries.

6. Confirm the data is replicated to the primary cluster.

   ```bash
   kubectl exec -it -n federation-demo primary-cluster-0 -- /coherence-operator/utils/cohctl get caches -o wide -I
   ```
 
   This will show 15,000

   ```bash
   SERVICE            CACHE              COUNT   SIZE  AVG SIZE    PUTS  GETS  REMOVES  EVICTIONS  HITS   MISSES  HIT PROB
   FederatedCache     test              15,000  12 MB       869  15,000     0        0          0     0        0     0.00%
   ```
   
## Undeploy all

```bash
kubectl delete -n federation-demo -f yaml/primary-cluster.yaml
```

```bash
kubectl delete -n federation-demo -f yaml/secondary-cluster.yaml
```

```bash
kubectl delete -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml
```

```bash
kubectl delete namespace federation-demo
```
