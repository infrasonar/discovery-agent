package main

import (
	"encoding/xml"
	"io"
	"os"
)

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun"`
	Start   int      `xml:"start,attr"`
	Version string   `xml:"version,attr"`
	Hosts   []Host   `xml:"host"`
}

type Host struct {
	XMLName   xml.Name  `xml:"host"`
	StartTime int       `xml:"starttime,attr"`
	EndTime   int       `xml:"endtime,attr"`
	Status    Status    `xml:"status"`
	Address   Address   `xml:"address"`
	Hostnames Hostnames `xml:"hostnames"`
	Ports     Ports     `xml:"ports"`
}

type Status struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"` // for example: "up"
}

type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`
	AddrType string   `xml:"addrtype,attr"`
}

type Hostnames struct {
	XMLName  xml.Name   `xml:"hostnames"`
	Hostname []Hostname `xml:"hostname"`
}

type Hostname struct {
	XMLName xml.Name `xml:"hostname"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"` // for example: "PTR"
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Port    []Port   `xml:"port"`
}

type Port struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"` // for example: "tcp"
	Id       int      `xml:"portid,attr"`   // for example: 443
	Service  Service  `xml:"service"`
	State    State    `xml:"state"`
}

type State struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr"` // for example: "open"
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`   // for example: "ftp"
	OsType  string   `xml:"ostype,attr"` // for example: "Linux" (or does not exist)
}

func parseXml(fn string) (*NmapRun, error) {
	// Open our xmlFile
	xmlFile, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var scan NmapRun
	err = xml.Unmarshal(byteValue, &scan)
	return &scan, err
}
