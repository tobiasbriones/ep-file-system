// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

package server

import (
	"fmt"
	"os"
)

const (
	fsRootPath = "fs"
	port       = 8080
	network    = "tcp"
)

func main() {

}

func getFilePath(relPath string) string {
	return fmt.Sprintf("%v%v%v", fsRootPath, os.PathSeparator, relPath)
}

func getServerAddress() string {
	return fmt.Sprintf("localhost:%v", port)
}
