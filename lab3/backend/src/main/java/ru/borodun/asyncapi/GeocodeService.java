package ru.borodun.asyncapi;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;
import org.jboss.logging.annotations.Param;

import javax.ws.rs.*;
import java.util.Set;
import java.util.concurrent.CompletionStage;

@Path("/1")
@RegisterRestClient(configKey="geocode-api")
public interface GeocodeService {

    @GET
    @Path("/geocode")
    Geocode getByName(@QueryParam("q") String place, @QueryParam("key") String key);

    @GET
    @Path("/geocode")
    CompletionStage<Geocode> getByNameAsync(@QueryParam("q") String place, @QueryParam("key") String key);

    @GET
    @Path("/geocode")
    Uni<Geocode> getByNameAsUni(@QueryParam("q") String place, @QueryParam("key") String key);
}