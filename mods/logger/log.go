// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package logger

import "github.com/miodzie/walter"

// Logger is meant for logging all incoming and outgoing Messages on a per-connection
// basis.
// There is a default null, file option, and ElasticSearch drivers available.
type Logger interface {
	Log(message walter.Message) error
}

type NullLogger struct {
}

func (n NullLogger) Log(message walter.Message) error {
	return nil
}

type FileLogger struct {
	file string
}

func (f FileLogger) Log(message walter.Message) error {
	//TODO implement me
	panic("implement me")
}
