package spdb.routingclient.services;

import com.graphhopper.GHRequest;
import com.graphhopper.GraphHopper;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;
import spdb.routingclient.config.GraphHopperConfiguration;
import spdb.routingclient.model.Route;
import spdb.routingclient.model.RouteRequest;

import java.util.Locale;

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

        return Route.builder().pointList(pointList).distance(distance).time(timeInMs).build();
    }
}
