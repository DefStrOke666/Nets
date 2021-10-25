package ru.borodun.asyncapi;

import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.QueryParam;
import java.util.concurrent.CompletionStage;

@RegisterRestClient(configKey = "weather-api")
public interface WeatherService {

    @GET
    Weather getByName(@QueryParam("lang") String lang, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("appid") String key);

    @GET
    CompletionStage<Weather> getByNameAsync(@QueryParam("lang") String lang, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("appid") String key);

    @GET
    Uni<Weather> getByNameAsUni(@QueryParam("lang") String lang, @QueryParam("lon") Float lon, @QueryParam("lat") Float lat, @QueryParam("appid") String key);
}
