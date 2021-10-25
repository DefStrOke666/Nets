package ru.borodun.asyncapi;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Places {
    public List<Feature> features;

    public static class Feature {
        public Properties properties;

        public static class Properties {
            public String xid;
            @NotEmpty
            public String name;
        }
    }
}
