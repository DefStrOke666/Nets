package ru.borodun.asyncapi;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class Weather {
    public Info main;
    public String name;
    public List<WeatherInfo> weather;
    public Wind wind;

    public static class Info {
        public String feels_like;
        public String humidity;
        public String pressure;
        public String temp;
    }

    public static class WeatherInfo{
        public String description;
    }

    public static class Wind{
        public Integer deg;
        public Float speed;
    }
}
