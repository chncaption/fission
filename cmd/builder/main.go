/*
Copyright 2018 The Fission Authors.

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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/fission/fission/cmd/builder/app"
)

// Usage: builder <shared volume path>
func main() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}()

	shareVolume := os.Args[1]
	if _, err := os.Stat(shareVolume); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(shareVolume, os.ModeDir|0700)
			if err != nil {
				logger.Fatal("error creating directory: %v", zap.Error(err), zap.String("directory", shareVolume))
			}
		}
	}

	err = app.Run(logger, shareVolume)
	logger.Error("error running builder", zap.Error(err))
}
