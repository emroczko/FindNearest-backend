package spdb.routingclient.config;

import com.graphhopper.GraphHopper;
import com.graphhopper.config.CHProfile;
import com.graphhopper.config.Profile;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GraphHopperConfiguration {

    private ConfigProperties configProperties;

    @Bean
    public GraphHopper create() {
        GraphHopper hopper = new GraphHopper();
        hopper.setOSMFile(configProperties.getOsmFilePath());
        hopper.setGraphHopperLocation(configProperties.getGraphHopperCachePath());
        hopper.setProfiles(new Profile("car").setVehicle("car").setWeighting("fastest").setTurnCosts(false));
        hopper.getCHPreparationHandler().setCHProfiles(new CHProfile("car"));
        hopper.importOrLoad();
        return hopper;
    }
}