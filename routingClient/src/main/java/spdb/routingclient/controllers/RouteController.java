package spdb.routingclient.controllers;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import spdb.routingclient.model.Route;
import spdb.routingclient.model.RouteRequest;
import spdb.routingclient.services.RouteService;

@Slf4j
@RestController
@AllArgsConstructor
public class RouteController {

    private RouteService routeService;

    @GetMapping(path = "/getRoute")
    public ResponseEntity<Route> getRoute(RouteRequest request) {
        log.info("Route request: {}", request);

        return ResponseEntity.ok(routeService.calculateRoute(request));
    }

}
