package com.oracle.coherence.demo.queues;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import com.tangosol.net.Coherence;
import com.tangosol.net.NamedQueue;

import java.io.IOException;
import java.io.OutputStream;

import java.net.InetSocketAddress;
import java.util.Random;

/**
 * Run a Coherence cluster and expose a "/publish/startOrder,numberOrders" endpoint on
 * :8888 to publish the requested number of orders from a starting order.
 * To publish 1000 orders starting at order 1:
 *
 *  curl http://localhost:8888/publish/1,1000
 */
public class RunCoherence {

    private static final   Random RANDOM = new Random();
    private static         Coherence coherence;

    public static void main(String[] args) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress(8888), 0);
        server.createContext("/publish", RunCoherence::publish);
        server.setExecutor(null);
        server.start();
        System.err.println("Server started on http://127.0.0.1:8888");

        coherence = Coherence.clusterMember().start().join();
    }

    /**
     * Called to publish orders to the Queue. Path param must be "/publisher/start,orders"
     * @param t
     * @throws IOException
     */
    private static void publish(HttpExchange t) throws IOException {
        try {
            String   path   = t.getRequestURI().getPath().replaceAll("/publish/", "");
            String[] values = path.split(",");
            if (values.length != 2) {
                throw new IllegalArgumentException("you must provide a start order number and number of orders");
            }

            NamedQueue<Order> queue      = coherence.getSession().getQueue("orders-queue");
            int               startOrder = Integer.parseInt(values[0]);
            int               numOrders  = Integer.parseInt(values[1]);

            for (int i = 0; i < numOrders; i++) {
                Order order = new Order(startOrder + i, "Customer-" + i + 1, "NEW", RANDOM.nextFloat() * 1000f + 10,
                        System.currentTimeMillis(), 0);
                if (!queue.offer(order)) {
                    throw new IllegalStateException("unable to offer order to queue");
                }
            }
            String message = "Added " + numOrders + " orders\n";
            System.err.println(message);
            send(t, 200, message);
        }
        catch (Exception e) {
            send(t, 500, "Error");
            e.printStackTrace();
        }
    }

    private static void send(HttpExchange t, int status, String body) throws IOException {
        t.sendResponseHeaders(status, body.length());
        OutputStream os = t.getResponseBody();
        os.write(body.getBytes());
        os.close();
    }
}
