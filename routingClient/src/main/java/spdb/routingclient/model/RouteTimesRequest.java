package spdb.routingclient.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NonNull;

import java.util.List;

@Data
@AllArgsConstructor
@Builder
public class RouteTimesRequest {
    @NonNull
    private Coordinates sourceCoordinates;
    @NonNull
    private List<LocationsDetails> locationsDetails;
}
