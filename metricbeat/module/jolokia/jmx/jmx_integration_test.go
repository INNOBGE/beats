// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build integration

package jmx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
)

func TestFetch(t *testing.T) {
	compose.EnsureUp(t, "jolokia")

	for _, config := range getConfigs() {
		f := mbtest.NewEventsFetcher(t, config)
		events, err := f.Fetch()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		t.Logf("%s/%s events: %+v", f.Module().Name(), f.Name(), events)
		if len(events) == 0 || len(events[0]) <= 1 {
			t.Fatal("Empty events")
		}
	}
}

func TestData(t *testing.T) {
	compose.EnsureUp(t, "jolokia")

	for _, config := range getConfigs() {
		f := mbtest.NewEventsFetcher(t, config)
		err := mbtest.WriteEvents(f, t)
		if err != nil {
			t.Fatal("write", err)
		}
	}
}

func getConfigs() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"module":     "jolokia",
			"metricsets": []string{"jmx"},
			"hosts":      []string{getEnvHost() + ":" + getEnvPort()},
			"namespace":  "testnamespace",
			"jmx.mappings": []map[string]interface{}{
				{
					"mbean": "java.lang:type=Runtime",
					"attributes": []map[string]string{
						{
							"attr":  "Uptime",
							"field": "uptime",
						},
					},
				},
				{
					"mbean": "java.lang:type=GarbageCollector,name=ConcurrentMarkSweep",
					"attributes": []map[string]string{
						{
							"attr":  "CollectionTime",
							"field": "gc.cms_collection_time",
						},
						{
							"attr":  "CollectionCount",
							"field": "gc.cms_collection_count",
						},
					},
				},
				{
					"mbean": "java.lang:type=Memory",
					"attributes": []map[string]string{
						{
							"attr":  "HeapMemoryUsage",
							"field": "memory.heap_usage",
						},
						{
							"attr":  "NonHeapMemoryUsage",
							"field": "memory.non_heap_usage",
						},
					},
				},
			},
		},
		{
			"module":     "jolokia",
			"metricsets": []string{"jmx"},
			"hosts":      []string{getEnvHost() + ":" + getEnvPort()},
			"namespace":  "testnamespace",
			"jmx.mappings": []map[string]interface{}{
				{
					"mbean": "Catalina:name=*,type=ThreadPool",
					"attributes": []map[string]string{
						{
							"attr":  "maxConnections",
							"field": "max_connections",
						},
						{
							"attr":  "port",
							"field": "port",
						},
					},
				},
				{
					"mbean": "Catalina:type=Server",
					"attributes": []map[string]string{
						{
							"attr":  "serverNumber",
							"field": "server_number_dosntconnect",
						},
					},
					"target": &TargetBlock{
						URL:      "service:jmx:rmi:///jndi/rmi://localhost:7091/jmxrmi",
						User:     "monitorRole",
						Password: "IGNORE",
					},
				},
				{
					"mbean": "Catalina:type=Server",
					"attributes": []map[string]string{
						{
							"attr":  "serverInfo",
							"field": "server_info_proxy",
						},
					},
					"target": &TargetBlock{
						URL:      "service:jmx:rmi:///jndi/rmi://localhost:7091/jmxrmi",
						User:     "monitorRole",
						Password: "QED",
					},
				},
			},
		},
	}
}

func getEnvHost() string {
	host := os.Getenv("JOLOKIA_HOST")

	if len(host) == 0 {
		host = "127.0.0.1"
	}
	return host
}

func getEnvPort() string {
	port := os.Getenv("JOLOKIA_PORT")

	if len(port) == 0 {
		port = "8778"
	}
	return port
}
