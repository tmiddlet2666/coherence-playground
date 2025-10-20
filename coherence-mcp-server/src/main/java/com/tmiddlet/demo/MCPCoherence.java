package com.tmiddlet.demo;

import java.util.Collections;
import java.util.Enumeration;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.CacheService;
import com.tangosol.net.Cluster;

import org.springframework.ai.tool.annotation.Tool;
import org.springframework.stereotype.Service;


@Service
public class MCPCoherence {
    @Tool(name = "getClusterInformation",
            description = "Returns the cluster information from the coherence cluster")
    public String getClusterInformation() {
        return String.format("The cluster information is %s", CacheFactory.getCluster());
    }

    @Tool(name = "getMemberList",
            description = "Returns the list of Coherence members in the cluster")
    public String getMemberList() {
        return String.format("The list of coherence members in the cluster is %s.", CacheFactory.getCluster().getMemberSet());
    }

    @Tool(name = "getServices",
            description = "Returns the list of services running in the cluster.")
    public String getServices() {
        Cluster cluster = CacheFactory.getCluster();
        StringBuilder sb = new StringBuilder("List of services");

        Collections.list(cluster.getServiceNames()).forEach(s ->
                sb.append(SEP).append("=").append(s).append(",serviceInfo").append(cluster.getServiceInfo(s)).append(SEP));

        return String.format(sb.toString());
    }

    @Tool(name = "listCaches",
            description = "Returns the list of caches running in the cluster.")
    public String listCaches() {
        final Cluster cluster = CacheFactory.getCluster();
        StringBuilder sb = new StringBuilder("List of caches").append(SEP);

        Collections.list(cluster.getServiceNames()).stream().map(cluster::getService).forEach(s -> {
            String serviceName = s.getInfo().getServiceName();

            if (s instanceof CacheService cacheService) {
                // retrieve the list of cache names for the service if it is a CacheService
                Collections.list((Enumeration<?>) cacheService.getCacheNames()).stream().map(String::valueOf).forEach(
                        c -> sb.append("service name=").append(serviceName).append(",cache name=").append(c)
                                .append(",size=").append(CacheFactory.getCache(c).size()).append(SEP));
            }
        });

        return String.format(sb.toString());
    }

    @Tool(name = "getCacheSize",
            description = "Returns the size of the specified coherence cache. Specify cache name.")
    public String getCacheSize(String cacheName) {
        if (cacheName == null || cacheName.trim().isEmpty()) {
            return "Error: you must provide a cache name.";
        }

        return String.format("The size of cache %s is %d.", cacheName, CacheFactory.getCache(cacheName).size());
    }

    @Tool(name = "getCacheEntries",
            description = "Returns the list of cache entries for the specified coherence cache. Specify cache name.")
    public String getCacheEntries(String cacheName) {
        if (cacheName == null || cacheName.trim().isEmpty()) {
            return "Error: you must provide a cache name.";
        }

        return String.format("The list of cache entries for %s is %s.", cacheName, CacheFactory.getCache(cacheName).entrySet());
    }

    private static final String SEP = System.lineSeparator();
}
