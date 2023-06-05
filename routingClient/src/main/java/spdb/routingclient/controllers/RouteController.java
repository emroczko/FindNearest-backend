package spdb.routingclient.controllers;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import spdb.routingclient.model.LocationRouteDetails;
import spdb.routingclient.model.Route;
import spdb.routingclient.model.RouteRequest;
import spdb.routingclient.model.RouteTimesRequest;
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
    public ResponseEntity<List<LocationRouteDetails>> getRouteData(@RequestBody RouteTimesRequest request) {
        log.info("Route times request");
        return ResponseEntity.ok(routeService.calculateTimes(request));
    }

}
