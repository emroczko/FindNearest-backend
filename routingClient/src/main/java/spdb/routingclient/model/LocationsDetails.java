package spdb.routingclient.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NonNull;

import java.math.BigInteger;

@Data
@AllArgsConstructor
@Builder
public class LocationsDetails {
    @NonNull
    private BigInteger id;
    @NonNull
    private Coordinates coordinates;
}
