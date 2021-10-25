package ru.borodun.asyncapi;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Geocode {
    public String message;
    public List<Hit> hits;

    public static class Hit {
        public String name;
        public String country;
        public String state;
        public Point point;

        public static class Point {
            public Float lat;
            public Float lng;
        }
    }
}