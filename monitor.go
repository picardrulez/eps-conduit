package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func monitorRoot(w http.ResponseWriter, r *http.Request) {
	var config = ReadConfig()
	backend := config.Backends
	strings.Replace(backend, " ", "", -1)
	backends := strings.Split(backend, ",")
	t := template.New("monitor Template")
	content := `
	<h1> EPS-CONDUIT MONITOR</h1>
	<h2>Backend Servers</h2>
	<table><tr><td>HOST</td><td>PORT</td><td>AVAILABILITY</td></tr>
	`
	for _, k := range backends {
		host := strings.Split(k, ":")
		ip := host[0]
		port := host[1]
		content = content + "<tr><td>" + ip + "</td><td>" + port + "</td>"
		if checkAvailability(ip, port, false) {
			content = content + "<td><font color='green'>Available</font></td></tr>"
		} else {
			content = content + "<td><font color='red'>Offline</font></td></tr>"
		}
	}
	content = content + "</table>"
	content = content +
		`
		<h2>Connections</h2>
		<table><tr><td>Last Minute</td><td>Last 5 Minutes</td><td>Last 15 Minutes</td><td>Last Hour</td></tr>
	`
	lastmin := selectNumLastMinute("routes")
	lastfive := selectNumLastFiveMinutes("routes")
	lastfifteen := selectNumLastFifteenMinutes("routes")
	lasthour := selectNumLastHour("routes")
	content = content + "<tr><td>" + strconv.Itoa(lastmin) + "</td><td>" + strconv.Itoa(lastfive) + "</td><td>" + strconv.Itoa(lastfifteen) + "</td><td>" + strconv.Itoa(lasthour) + "</td></tr></table>"
	content = content + "</br>"
	content = content + `
	<h2>Errors</h2>
	<table><tr><td>Last Minute</td><td>Last 5 Minutes</td><td>Last 15 Minutes</td><td>Last Hour</td></tr>
	`
	errmin := selectNumLastMinute("errors")
	errfive := selectNumLastFiveMinutes("errors")
	errfifteen := selectNumLastFifteenMinutes("errors")
	errhour := selectNumLastHour("errors")
	content = content + "<tr><td>" + strconv.Itoa(errmin) + "</td><td>" + strconv.Itoa(errfive) + "</td><td>" + strconv.Itoa(errfifteen) + "</td><td>" + strconv.Itoa(errhour) + "</td></tr></table>"

	createMonitorPage := Page{Content: content}
	pageTemplate := pageBuilder()
	t, _ = t.Parse(pageTemplate)
	t.Execute(w, createMonitorPage)
}

func checkAvailability(host string, port string, updatedb bool) bool {
	_, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)

	if err != nil {
		if updatedb {
			insertError(host, err)
		}
		log.Printf("conn error:  %s", err)
		return false
	} else {
		return true
	}
}
