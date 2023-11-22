/*
 * MIT-X11 Open Source License
 *
 * Copyright (c) 2022, Advanced Micro Devices, Inc.
 * All rights reserved.
 *
 * Developed by:
 *
 *                 AMD Research and AMD Software Development
 *
 *                 Advanced Micro Devices, Inc.
 *
 *                 www.amd.com
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sellcopies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 *  - The above copyright notice and this permission notice shall be included in
 *    all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Except as contained in this notice, the name of the Advanced Micro Devices,
 * Inc. shall not be used in advertising or otherwise to promote the sale, use
 * or other dealings in this Software without prior written authorization from
 * the Advanced Micro Devices, Inc.
 *
 */

package main

import (
	"io"
	"log"
	"strconv"
	"strings"
	"os/exec"
	"net/http"
	"encoding/json"
	"src/collect" // This has the implementation of the Scan() function
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gGPUProductNames[24] string
var gNodeName string
var gPod string 
/* rocm-smi output sample
{"card0": 
    {
	"Card series": "Instinct MI210", 
	"Card model": "0x0c34", 
	"Card vendor": "Advanced Micro Devices, Inc. [AMD/ATI]", 
	"Card SKU": "D67301"
    }
}*/

type Card struct {
	Cardseries string `json:"cardseries"`
	Cardmodel  string `json:"cardmodel"`
	Cardvendor string `json:"cardvendor"`
	CardSKU    string `json:"cardsku"`
}

type Info map[string]Card

func GetGpuProductNames() {
	//rocm-smi output in json format
        rocmsmiJson, err := exec.Command("rocm-smi", "--showproductname", "--json").Output()
	if err == nil {
	    //format for json unmarshalling
	    productnamesJson := strings.ReplaceAll(strings.ToLower(string(rocmsmiJson)), " ", "")
	    dec := json.NewDecoder(strings.NewReader(productnamesJson))
	    for {
		var cards Info
		if err := dec.Decode(&cards); err == io.EOF {
		    break
		} else if err != nil {
		    log.Fatal(err)
		}
		//iterate over each card
		for k, _ := range cards {
		    deviceID, err := strconv.Atoi(strings.TrimPrefix(string(k), "card"))
		    if err == nil {
			gGPUProductNames[deviceID] = cards[k].Cardseries
		    }
		}
	    }
	}
}

func GetNodeName() {
	nodename, err := exec.Command("uname", "-n").Output()
	if err == nil {
	    gNodeName = string(nodename)
	    gPod = string(nodename)
	} else {
	    log.Fatal(err)
	}
}

func main() {
	// Get Node name
	GetNodeName()
	// Get all GPU product names
	GetGpuProductNames()

	// Called on each collector.Collect.
	handle := func() (collect.AMDParams) {
		return collect.Scan()
	}

	// Make Prometheus client aware of our collector.
	c := NewCollector(handle)
	prometheus.MustRegister(c)

	// Set up HTTP handler for metrics.
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Start listening for HTTP connections.
	const addr = ":2021"
	log.Printf("starting collector exporter on %q", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("cannot start collector exporter: %s", err)
	}
}
