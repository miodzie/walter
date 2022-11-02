// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package drivers

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/logger"
)

type MultiLogger struct {
	loggers []logger.Logger
}

func NewMultiLogger(args ...logger.Logger) *MultiLogger {
	return &MultiLogger{loggers: args}
}

func (l MultiLogger) Log(message walter.Message) error {
	var err error
	for _, l := range l.loggers {
		err = l.Log(message)
	}

	// Just return the last error for now. w/e
	return err
}
