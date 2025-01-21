# Frameworks starter example

This has three started projects for using Coherence with Helidon, Spring Boot and Micronaut.
      
## Requirements

* JDK21
* Maven 3.8+

## Building each of the projects.
              
Run the following in each of the subdirectories:

```bash
mvn clean install
```

## Starting the application

### Run Helidon (Jakarta)

```bash
cd helidon
java -jar target/helidon.jar
```

Run additional cache server:

```bash
java -Dmain.class=com.tangosol.net.Coherence -Dserver.port=-1 -Dcoherence.management.http=none -jar target/helidon.jar  
```

### Run Helidon (javax)

```bash
cd helidon-javax
java -jar target/helidon-javax.jar
```

Run additional cache server:

```bash
java -Dmain.class=com.tangosol.net.Coherence -Dserver.port=-1 -Dcoherence.management.http=none -jar target/helidon-javax.jar  
```

### Run Spring Boot

```bash
cd springboot
java -jar target/springboot-1.0-SNAPSHOT.jar
```
      
Run additional cache server:

```bash
java -Dserver.port=-1 -Dloader.main=com.tangosol.net.Coherence -Dcoherence.management.http=none -jar target/springboot-1.0-SNAPSHOT.jar
```

### Run Micronaut

```bash
cd micronaut
java -jar target/micronaut-1.0-SNAPSHOT-shaded.jar
```
    
Run additional cache server:

```bash
java -Dmicronaut.main.class=com.tangosol.net.Coherence -Dcoherence.management.http=none -Dmicronaut.server.port=-1 -jar target/micronaut-1.0-SNAPSHOT-shaded.jar
```

## Using the application

Each of the application servers start on 8080 and expose the following REST endpoints.

### List Customers

```bash
curl http://localhost:8080/api/customers
```

### Add a Customer

```bash
curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "name": "Tim", "balance": 1000}' http://localhost:8080/api/customers
```

```bash
curl -X POST -H "Content-Type: application/json" -d '{"id": 2, "name": "John", "balance": 5000}' http://localhost:8080/api/customers
```

### Get a specific customer

```bash
curl http://localhost:8080/api/customers/1
```

### Delete a specific customer

```bash
curl -X DELETE http://localhost:8080/api/customers/1
```