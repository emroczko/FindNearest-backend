package spdb.routingclient.config;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Data
@Configuration
@ConfigurationProperties(prefix = "app")
public class ConfigProperties {

    private String osmFilePath;
    private String graphHopperCachePath;
}
