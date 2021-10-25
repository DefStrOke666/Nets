package ru.borodun.asyncapi;

import io.smallrye.common.annotation.Blocking;
import io.smallrye.mutiny.Uni;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import org.jboss.resteasy.reactive.RestQuery;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import java.util.concurrent.CompletionStage;

@Path("/v1")
public class DescriptionResource {

    @RestClient
    DescriptionService descriptionService;
    String key = System.getenv("OPENTRIPMAP_KEY");

    @GET
    @Path("/description/{xid}")
    @Blocking
    public Description name(String xid) {
        return descriptionService.getByName(xid, key);
    }

    @GET
    @Path("/description-async/{xid}")
    public CompletionStage<Description> nameAsync(String xid) {
        return descriptionService.getByNameAsync(xid, key);
    }

    @GET
    @Path("/description-uni/{xid}")
    public Uni<Description> nameUni(String xid) {
        return descriptionService.getByNameAsUni(xid, key);
    }

}
