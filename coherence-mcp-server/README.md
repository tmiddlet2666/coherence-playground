# Coherence MCP Server Example

This repository contains an demo of an Model Context Protocol (MCP) server to 
expose various Coherence related functionality to AI modesl.

In this example you will:

1. Compile the Coherence MCP Server example using JDK21+
2. Configure Claude Desktop with your MCP Server
3. Interact with the LLM to ask specific questions on the Coherence Cluster

Please see the Medium Article [available here](https://medium.com/@middleton.music.au/integrating-coherence-with-claude-desktop-using-an-mcp-server-9f665feec989) for full instructions and walkthrough.

## Project Structure

- `src/main/java`: Contains the main Java source code
- `src/main/resources`: Contains application resources

## Pre-requisites

* JDK - You must have Maven and JDK21+ or latest installed
* LLM - I'm using Claude Desktop for this example

## Build the Project

Use the following command to build the project:

```bash
./mvnw clean package
```

Issue the following to see the runnable jar:

```bash
ls -l target/*.jar
```

You should see output similar to the following, note this jar name for the config below.

```bash
-rw-r--r--@ 1 timbo  staff  49274136 16 Oct 10:44 target/coherence-mcp-server-0.0.1-SNAPSHOT.jar
```

## Run the Example

### Start Coherence Cluster

In this example we can start a Coherence storage member, scoped to the current machine, using following:

```bash
./mvnw exec:exec
```

Run the above in two separate terminals.  Maven will be downloaded if you dont already have it.

> Note: You can change the `mvnw` command to `mvn` if you have Maven configured.

### Start the Coherence and add data

1. Run the following command to start the console:
   ```bash
   ./mvwn exec:exec -Pconsole
   ```

2. At the `Map (?):` prompt type:
   ```bash
   cache test1
   ```
   
3. Add two entries using:
    ```bash
    put "1" "Tim"
    put "2" "John"
    list
    ```
   
4. Create a new cache called `test2` by using the following:
   ```bash
   cache test2
   ```

5. Add some random binary data using:
    ```bash
    bulkput 10000
    ```
   
Type `bye` to exit the console.

## Configure Claude Desktop

### Add the Coherence MCP Server

Open Claude Desktop and go into `Settings`, click `Developer` and `Edit Config`. 

Use the editor of you choice to add the following to the contents of the file `claude_desktop_config.json` and include the following entry,
ensuring the `/path/to/coherence-playground` is replaced with full path such as `/Users/tmiddlet/github/path/to/coherence-playground` in my case.


```json
{
  "mcpServers": {
    "coherence-mcp-server-local": {
      "command": "java",
      "args": [
        "-Xmx1g",
        "-Xms1g",
        "-Dcoherence.wka=127.0.0.1",
        "-Dcoherence.cluster=demo-cluster",
        "-Dcoherence.distributed.localstorage=false",
        "-jar",
        "/path/to/coherence-playground/coherence-mpc-server/target/coherence-mcp-server-0.0.1-SNAPSHOT.jar"
      ]
    }
  }
}

```

You need to quit and reload Claude Desktop for the changes to be picked up.

### Interact with the LLM.

Follow the Article [available here](https://medium.com/@middleton.music.au/integrating-coherence-with-claude-desktop-using-an-mcp-server-9f665feec989) for full steps taking to query the cluster.


