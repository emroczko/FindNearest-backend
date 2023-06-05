package spdb.routingclient.services;

import com.graphhopper.GHRequest;
import com.graphhopper.GraphHopper;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import spdb.routingclient.config.GraphHopperConfiguration;
import spdb.routingclient.model.*;

import java.util.ArrayList;
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

    public List<LocationRouteDetails> calculateTimes(RouteTimesRequest routeRequest) {

        List<LocationRouteDetails> checkedLocations = new ArrayList<>();
        for (var location : routeRequest.getLocationsDetails()) {
            var req = new GHRequest(
                    routeRequest.getSourceCoordinates().getLatitude(),
                    routeRequest.getSourceCoordinates().getLongitude(),
                    location.getCoordinates().getLatitude(),
                    location.getCoordinates().getLongitude()
            )
                    .setProfile("car")
                    .setLocale(Locale.US);
            var rsp = hopper.route(req);

            if (rsp.hasErrors())
                throw new RuntimeException(rsp.getErrors().toString());

            var path = rsp.getBest();

            var distance = path.getDistance();
            var timeInMs = path.getTime();

            checkedLocations.add(LocationRouteDetails.builder()
                    .locationsDetails(location)
                    .time(timeInMs)
                    .distance(distance)
                    .build());
        }

        return checkedLocations;
    }
}
