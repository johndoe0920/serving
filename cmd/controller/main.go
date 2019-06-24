/*
Copyright 2018 The Knative Authors

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
	"log"
	"os"
	"github.com/knative/pkg/configmap"
	"github.com/knative/serving/pkg/logging"
	"github.com/knative/pkg/metrics"
	"go.uber.org/zap"
	"github.com/knative/pkg/logging/logkey"
	// The set of controllers this controller process runs.
	"github.com/knative/serving/pkg/reconciler/configuration"
	"github.com/knative/serving/pkg/reconciler/labeler"
	"github.com/knative/serving/pkg/reconciler/revision"
	"github.com/knative/serving/pkg/reconciler/route"
	"github.com/knative/serving/pkg/reconciler/serverlessservice"
	"github.com/knative/serving/pkg/reconciler/service"

	// This defines the shared main for injected controllers.
	"github.com/knative/pkg/injection/sharedmain"
)

const (
    component = "controller"
)

func main() {
	cm, err := configmap.Load("/etc/config-logging")
	if err != nil {
		log.Fatal("Error loading logging configuration:", err)
	}
	logConfig, err := logging.NewConfigFromMap(cm)
	if err != nil {
		log.Fatal("Error loading logging configuration:", err)
	}
	createdLogger, _ := logging.NewLoggerFromConfig(logConfig, component)
	logger := createdLogger.With(zap.String(logkey.ControllerType, "activator"))
	defer flush(logger)

	logger.Error("Wolverine")
	sharedmain.Main("controller",
		configuration.NewController,
		labeler.NewRouteToConfigurationController,
		revision.NewController,
		route.NewController,
		serverlessservice.NewController,
		service.NewController,
	)
}

func flush(logger *zap.SugaredLogger) {
	logger.Sync()
	os.Stdout.Sync()
	os.Stderr.Sync()
	metrics.FlushExporter()
}

