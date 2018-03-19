/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

var runtimeStatusCommand = cli.Command{
	Name:      "info",
	Usage:     "Display information of the container runtime",
	ArgsUsage: "",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "json",
			Usage: "Output format, One of: json|yaml",
		},
	},
	Action: func(context *cli.Context) error {
		err := Info(context, runtimeClient)
		if err != nil {
			return fmt.Errorf("getting status of runtime failed: %v", err)
		}
		return nil
	},
	Before: getRuntimeClient,
	After:  closeConnection,
}

// Info sends a StatusRequest to the server, and parses the returned StatusResponse.
func Info(cliContext *cli.Context, client pb.RuntimeServiceClient) error {
	request := &pb.StatusRequest{Verbose: true}
	logrus.Debugf("StatusRequest: %v", request)
	r, err := client.Status(context.Background(), request)
	logrus.Debugf("StatusResponse: %v", r)
	if err != nil {
		return err
	}

	status, err := protobufObjectToJSON(r.Status)
	if err != nil {
		return err
	}
	return outputStatusInfo(status, r.Info, cliContext.String("output"))
}
