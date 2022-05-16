// Copyright (c) 2022 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/ep-file-system-server

// Package io Models a file system according to
// https://github.com/tobiasbriones/cp-unah-mm545-distributed-text-file-system/tree/main/model
package io

const (
	Root           = ""
	Separator      = "/"
	ValidPathRegex = "^$|\\\\w+/*\\\\.*-*"
)
