package evohome

import (
	"os"
	"testing"
	"time"
)

func TestNewEvohome(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if len(client.installations) == 0 {
		t.Fatal("No installations found on account")
 	}
}


func TestEvohome_TemperatureControlSystem(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	system := client.TemperatureControlSystem()
	if system == nil {
		t.Fatal("Control system is nil")
	}
}

func TestEvohome_TemperatureControlSystemByIndex(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	system := client.TemperatureControlSystemByIndex(1)
	if system == nil {
		t.Fatal("Control system is nil")
	}
}

func TestEvoHome_getZoneStatusIncludingTemperature(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	locationID := client.Installations()[0].Location.Id
	zones := client.getZoneStatusIncludingTemperature(locationID)

	if zones == nil {
		t.Fatalf("No zones in location with id %s", locationID)
	}
}

func TestEvohome_UpdateTemperatures(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	zoneTempsBefore := map[string]map[string]float32{}
	for _, install := range client.Installations() {
		zoneTempsBefore[install.Location.Name] = getZoneTemperatureMap(install.Gateways[0].TemperatureControlSystems[0].Zones)
	}

	client.UpdateTemperatures()

	for _, install := range client.Installations() {
		zoneTempsAfter := getZoneTemperatureMap(install.Gateways[0].TemperatureControlSystems[0].Zones)

		system := install.Gateways[0].TemperatureControlSystems[0]
		for _, v := range system.Zones {
			tempBefore := zoneTempsBefore[install.Location.Name][v.Name]
			tempAfter := zoneTempsAfter[v.Name]

			if tempBefore == tempAfter {
				t.Errorf("Temperature in zone %s was not updated. Temp before: %v after: %v", v.Name, tempBefore, tempAfter)
			}
		}
	}
}

func TestEvohome_UpdateSchedules(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	client.UpdateSchedules()

	for _, install := range client.Installations() {
		for _, zone := range install.Gateways[0].TemperatureControlSystems[0].Zones {
			if zone.Schedules == nil {
				t.Errorf("%s does not have a defined schedule", zone.Name)
			}
		}
	}
}

func TestEvohome_Update(t *testing.T) {
	client, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	zoneTempsBefore := map[string]map[string]float32{}
	for _, install := range client.Installations() {
		zoneTempsBefore[install.Location.Name] = getZoneTemperatureMap(install.Gateways[0].TemperatureControlSystems[0].Zones)
	}

	client.Update()

	for _, install := range client.Installations() {
		zoneTempsAfter := getZoneTemperatureMap(install.Gateways[0].TemperatureControlSystems[0].Zones)

		system := install.Gateways[0].TemperatureControlSystems[0]
		for _, zone := range system.Zones {
			tempBefore := zoneTempsBefore[install.Location.Name][zone.Name]
			tempAfter := zoneTempsAfter[zone.Name]

			if tempBefore == tempAfter {
				t.Errorf("Temperature in zone %s was not updated. Temp before: %v after: %v", zone.Name, tempBefore, tempAfter)
			}

			if zone.Schedules == nil {
				t.Errorf("%s does not have a defined schedule", zone.Name)
			}
		}
	}
}

func TestRequireAuthentication(t *testing.T) {
	if !requireAuthentication(time.Now()) {
		t.Error("should have required authentication")
	}

	if requireAuthentication(time.Now().Add(time.Duration(100) * time.Second)) {
		t.Error("should not have required authentication")
	}
}

func getClient() (*Evohome, error) {
	username, password := getCredentials()
	return NewEvohome(username,password)
}

func getZoneTemperatureMap(zones []Zone) map[string]float32 {
	temps := map[string]float32{}

	for _, v := range zones {
		if v.TemperatureStatus == nil {
			temps[v.Name] = 0
		} else {
			temps[v.Name] = v.TemperatureStatus.Temperature
		}
	}

	return temps
}

func getCredentials()(string, string) {
	return os.Getenv("EVOHOME_USERNAME"), os.Getenv("EVOHOME_PASSWORD")
}