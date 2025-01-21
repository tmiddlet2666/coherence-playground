package com.oracle.coherence.demo.frameworks.helidon;

import static javax.ws.rs.core.MediaType.APPLICATION_JSON;

import com.oracle.coherence.cdi.Name;

import com.tangosol.net.NamedMap;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;

import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Response;

@Path("/api/customers")
@ApplicationScoped
public class CustomerResource {

    @Inject
    @Name("tasks")
    private NamedMap<Integer, Customer> customers;

    @POST
    @Consumes(APPLICATION_JSON)
    public Response createCustomer(Customer customer) {
        customers.put(customer.getId(), customer);
        return Response.accepted(customer).build();
    }

    @GET
    @Produces(APPLICATION_JSON)
    public Response getCustomers() {
        return Response.ok(customers.values()).build();
    }

    @GET
    @Path("{id}")
    @Produces(APPLICATION_JSON)
    public Response getTask(@PathParam("id") int id) {
        Customer customer = customers.get(id);

        return customer == null ? Response.status(Response.Status.NOT_FOUND).build() : Response.ok(customer).build();
    }

    @DELETE
    @Path("{id}")
    @Produces(APPLICATION_JSON)
    public Response deleteTask(@PathParam("id") int id) {
        Customer oldCustomer = customers.remove(id);
        return oldCustomer == null ? Response.status(Response.Status.NOT_FOUND).build() : Response.ok(oldCustomer).build();
    }
}
