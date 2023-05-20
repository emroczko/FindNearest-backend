package spdb.routingclient.model;

import lombok.Data;

@Data
public class RouteRequest {
    private Double sourceLongitude;
    private Double sourceLatitude;
    private Double destinationLongitude;
    private Double destinationLatitude;
}
