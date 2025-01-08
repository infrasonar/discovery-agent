package main

import (
	"fmt"
)

func getAssetName(host *Host) interface{} {
	for _, hostname := range host.Hostnames.Hostname {
		if hostname.Type == "PTR" && hostname.Name != "" {
			return hostname.Name
		}
	}
	for _, hostname := range host.Hostnames.Hostname {
		if hostname.Name != "" {
			return hostname.Name
		}
	}
	return nil
}

func getPortsFromHost(host *Host) []int {
	ports := []int{}
	for _, port := range host.Ports.Port {
		if port.State.State == "open" {
			ports = append(ports, port.Id)
		}
	}
	return ports
}

func getServicesFromHost(host *Host) []string {
	services := []string{}
	for _, port := range host.Ports.Port {
		if port.State.State == "open" {
			if port.Service.Name == "" {
				// Usually, the service name is not empty but already set to
				// unknown by the NMAP scanner
				port.Service.Name = "unknown"
			}
			services = append(services, port.Service.Name)
		}
	}
	return services
}

func getOsType(host *Host) interface{} {
	for _, port := range host.Ports.Port {
		if port.Service.OsType != "" {
			return port.Service.OsType
		}
	}
	return nil
}

func getAssets(hosts *[]Host) []map[string]any {
	assets := []map[string]any{}
	var asset map[string]any

	for _, host := range *hosts {
		asset = map[string]any{
			"name":          host.Address.Addr,
			"Name":          getAssetName(&host),
			"IpAddress":     host.Address.Addr,
			"IpAddressType": host.Address.AddrType,
			"Ports":         getPortsFromHost(&host),
			"OsType":        getOsType(&host),
			"Services":      getServicesFromHost(&host),
		}
		assets = append(assets, asset)
	}
	return assets
}

func getAllPorts(hosts *[]Host) []map[string]any {
	type MPort struct {
		Service string
		Count   int
	}
	portmap := map[int]*MPort{}

	for _, host := range *hosts {
		for _, port := range host.Ports.Port {
			m, ok := portmap[port.Id]
			if ok {
				m.Count += 1
			} else {
				portmap[port.Id] = &MPort{
					Service: port.Service.Name,
					Count:   1,
				}
			}
		}
	}

	ports := []map[string]any{}
	var port map[string]any
	for id, m := range portmap {
		port = map[string]any{
			"name":    fmt.Sprintf("%d", id),
			"Port":    id,
			"Service": m.Service,
			"Count":   m.Count,
		}
		ports = append(ports, port)
	}
	return ports
}

func getState(scan *NmapRun) map[string][]map[string]any {
	state := map[string][]map[string]any{}
	state["assets"] = getAssets(&scan.Hosts)
	state["ports"] = getAllPorts(&scan.Hosts)
	state["agent"] = []map[string]any{{
		"name":        "discovery",
		"Version":     Version,
		"NmapVersion": scan.Version,
	}}
	return state
}
