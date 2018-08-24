//
// Copyright (c) 2018 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterd

import (
	"github.com/heketi/heketi/executors"
)

func (g *executor) snapshotActivate(host string, snapshot string) error {
	logger.Debug("BEGIN")
	return executors.NotSupportedError
}

func (g *executor) SnapshotCloneVolume(host string, vcr *executors.SnapshotCloneRequest) (*executors.Volume, error) {
	logger.Debug("BEGIN")
	return nil, executors.NotSupportedError
}

func (g *executor) SnapshotCloneBlockVolume(host string, vcr *executors.SnapshotCloneRequest) (*executors.BlockVolumeInfo, error) {
	logger.Debug("BEGIN")
	// TODO: cloning of block volume is not implemented yet
	return nil, executors.NotSupportedError
}

func (g *executor) SnapshotDestroy(host string, snapshot string) error {
	logger.Debug("BEGIN")
	return executors.NotSupportedError
}
