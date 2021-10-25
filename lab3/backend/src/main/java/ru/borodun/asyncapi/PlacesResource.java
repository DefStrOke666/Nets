package ru.borodun.asyncapi;

import io.smallrye.common.annotation.Blocking;
import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import org.jboss.resteasy.reactive.RestQuery;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;
import java.util.concurrent.CompletionStage;

@Path("/v1")
public class PlacesResource {

    @RestClient
    PlacesService placesService;
    String key = System.getenv("OPENTRIPMAP_KEY");

    @GET
    @Path("/places")
    @Blocking
    public Places name(@RestQuery String radius, @RestQuery Float lng, @RestQuery Float lat) {
        return placesService.getByName(radius, lng, lat, key);
    }

    @GET
    @Path("/places-async")
    public CompletionStage<Places> nameAsync(@RestQuery String radius, @RestQuery Float lng, @RestQuery Float lat) {
        return placesService.getByNameAsync(radius, lng, lat, key);
    }

    @GET
    @Path("/places-uni")
    public Uni<Places> nameUni(@RestQuery String radius, @RestQuery Float lng, @RestQuery Float lat) {
        return placesService.getByNameAsUni(radius, lng, lat, key);
    }

}
