package ru.borodun.asyncapi;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;
import java.util.concurrent.CompletionStage;

@Path("/xid")
@RegisterRestClient(configKey = "places-api")
public interface DescriptionService {

    @GET
    @Path("/{xid}")
    Description getByName(String xid, @QueryParam("apikey") String key);

    @GET
    @Path("/{xid}")
    CompletionStage<Description> getByNameAsync(String xid, @QueryParam("apikey") String key);

    @GET
    @Path("/{xid}")
    Uni<Description> getByNameAsUni(String xid, @QueryParam("apikey") String key);
}
