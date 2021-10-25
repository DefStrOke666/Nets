package ru.borodun.asyncapi;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;
import java.util.concurrent.CompletionStage;

@Path("/radius")
@RegisterRestClient(configKey = "places-api")
public interface PlacesService {

    @GET
    Places getByName(@QueryParam("radius") String radius, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("apikey") String key);

    @GET
    CompletionStage<Places> getByNameAsync(@QueryParam("radius") String radius, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("apikey") String key);

    @GET
    Uni<Places> getByNameAsUni(@QueryParam("radius") String radius, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("apikey") String key);
}
