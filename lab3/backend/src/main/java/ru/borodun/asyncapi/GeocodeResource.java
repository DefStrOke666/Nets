package ru.borodun.asyncapi;


import io.smallrye.common.annotation.Blocking;
import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import java.util.Set;
import java.util.concurrent.CompletionStage;

@Path("/v1")
public class GeocodeResource {
    @RestClient
    GeocodeService geocodeService;
    String key = System.getenv("GEOCODE_KEY");

    @GET
    @Path("/geocode/{place}")
    @Blocking
    public Geocode name(String place) {
        return geocodeService.getByName(place, key);
    }

    @GET
    @Path("/geocode-async/{place}")
    public CompletionStage<Geocode> nameAsync(String place) {
        return geocodeService.getByNameAsync(place, key);
    }

    @GET
    @Path("/geocode-uni/{place}")
    public Uni<Geocode> nameUni(String place) {
        return geocodeService.getByNameAsUni(place, key);
    }
}
