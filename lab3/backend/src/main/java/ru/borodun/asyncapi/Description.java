package ru.borodun.asyncapi;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Description {
    public String name;
    public String wikipedia;
    public String otm;
    public String image;
    public Info info;

    public static class Info {
        public String descr;
    }
}
