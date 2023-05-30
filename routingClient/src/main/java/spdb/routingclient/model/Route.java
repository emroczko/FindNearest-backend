package spdb.routingclient.model;

import com.graphhopper.util.PointList;
import lombok.Builder;
import lombok.Data;

import java.util.List;

@Data
@Builder
public class Route {
    private List<Coordinates> coordinates;
    private Double distance;
    private Long time;
}
