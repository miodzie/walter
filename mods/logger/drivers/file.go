// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package drivers

import (
	"walter"
)

type FileLogger struct {
}

func (l FileLogger) Log(message walter.Message) error {
	//TODO implement me
	panic("implement me")
	return nil
}
