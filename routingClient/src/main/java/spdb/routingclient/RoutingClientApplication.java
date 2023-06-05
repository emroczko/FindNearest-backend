package spdb.routingclient;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import spdb.routingclient.config.ConfigProperties;

@SpringBootApplication
@EnableConfigurationProperties(ConfigProperties.class)
public class RoutingClientApplication {

    public static void main(String[] args) {
        SpringApplication.run(RoutingClientApplication.class, args);
    }

}
