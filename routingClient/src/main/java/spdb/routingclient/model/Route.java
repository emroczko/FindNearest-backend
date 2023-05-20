package spdb.routingclient.model;

import com.graphhopper.util.PointList;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class Route {
    private PointList pointList;
    private Double distance;
    private Long time;
}
