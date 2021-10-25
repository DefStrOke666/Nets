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
public class WeatherResource {

    @RestClient
    WeatherService weatherService;
    String key = System.getenv("OPENWEATHERMAP_KEY");
    String lang = "ru";

    @GET
    @Path("/weather")
    @Blocking
    public Weather name(@RestQuery Float lng, @RestQuery Float lat) {
        return weatherService.getByName(lang, lng, lat, key);
    }

    @GET
    @Path("/weather-async")
    public CompletionStage<Weather> nameAsync(@RestQuery Float lng, @RestQuery Float lat) {
        return weatherService.getByNameAsync(lang, lng, lat, key);
    }

    @GET
    @Path("/weather-uni")
    public Uni<Weather> nameUni(@RestQuery Float lng, @RestQuery Float lat) {
        return weatherService.getByNameAsUni(lang, lng, lat, key);
    }
}
