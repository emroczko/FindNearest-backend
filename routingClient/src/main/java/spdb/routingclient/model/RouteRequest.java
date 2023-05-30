package spdb.routingclient.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NonNull;

@Data
@AllArgsConstructor
@Builder
public class RouteRequest {
    @NonNull
    private Double sourceLongitude;
    @NonNull
    private Double sourceLatitude;
    @NonNull
    private Double destinationLongitude;
    @NonNull
    private Double destinationLatitude;
}
