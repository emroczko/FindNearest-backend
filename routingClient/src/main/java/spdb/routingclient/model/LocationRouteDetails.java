package spdb.routingclient.model;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class LocationRouteDetails {
    private LocationsDetails locationsDetails;
    private Double distance;
    private Long time;
}
