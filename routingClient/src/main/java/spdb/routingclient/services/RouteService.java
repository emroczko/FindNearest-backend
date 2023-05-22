package spdb.routingclient.services;

import com.graphhopper.GHRequest;
import com.graphhopper.GraphHopper;
import com.graphhopper.util.Parameters;
import com.graphhopper.util.details.PathDetail;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import spdb.routingclient.config.GraphHopperConfiguration;
import spdb.routingclient.model.Coordinates;
import spdb.routingclient.model.Route;
import spdb.routingclient.model.RouteRequest;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Locale;

@Slf4j
@Service
public class RouteService {
    private final GraphHopper hopper;

    public RouteService(GraphHopperConfiguration graphHopperConfiguration) {
        this.hopper = graphHopperConfiguration.create();
    }

    public Route calculateRoute(RouteRequest routeRequest) {

        var req = new GHRequest(
                routeRequest.getSourceLatitude(),
                routeRequest.getSourceLongitude(),
                routeRequest.getDestinationLatitude(),
                routeRequest.getDestinationLongitude()
        )
                .setProfile("car")
                .setLocale(Locale.US);
        var rsp = hopper.route(req);

        if (rsp.hasErrors())
            throw new RuntimeException(rsp.getErrors().toString());

        var path = rsp.getBest();

        var pointList = path.getPoints();
        var distance = path.getDistance();
        var timeInMs = path.getTime();

        List<Coordinates> coordinates = new ArrayList<>();
        pointList.forEach(point -> coordinates.add(Coordinates.builder()
                .latitude(point.lat)
                .longitude(point.lon)
                .build()));

        return Route.builder().coordinates(coordinates).distance(distance).time(timeInMs).build();
    }
}
