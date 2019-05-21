// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()
	secs := now.Unix()

	// current UNIX time in UTC formatted as RFC3339
	fmt.Println(now.Format(time.RFC3339))
	// the number of seconds since the UNIX epoch
	fmt.Println(secs)
}
