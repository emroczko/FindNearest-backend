package spdb.routingclient.controllers;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import spdb.routingclient.model.LocationRouteDetails;
import spdb.routingclient.model.Route;
import spdb.routingclient.model.RouteRequest;
import spdb.routingclient.model.RouteDetailsRequest;
import spdb.routingclient.services.RouteService;

import java.util.List;

@Slf4j
@RestController
@AllArgsConstructor
@RequestMapping(value = "/api/v1")
public class RouteController {

    private RouteService routeService;

    @PostMapping(path = "/getRoute")
    public ResponseEntity<Route> getRoute(@RequestBody RouteRequest request) {
        log.info("Route request: {}", request);

        return ResponseEntity.ok(routeService.calculateRoute(request));
    }

    @PostMapping(path = "/getRouteData")
    public ResponseEntity<List<LocationRouteDetails>> getRouteData(@RequestBody RouteDetailsRequest request) {
        log.info("Route details request for {} by {}", request.getSourceCoordinates(), request.getMeanOfTransport());
        return ResponseEntity.ok(routeService.calculateDetails(request));
    }

}
