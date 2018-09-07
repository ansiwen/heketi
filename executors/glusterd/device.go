//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterd

import (
	"encoding/json"
	"fmt"

	"github.com/heketi/heketi/executors"
	"github.com/pkg/errors"
)

// Read:
// https://access.redhat.com/documentation/en-US/Red_Hat_Storage/3.1/html/Administration_Guide/Brick_Configuration.html
//
//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

// Read:
// https://access.redhat.com/documentation/en-US/Red_Hat_Storage/3.1/html/Administration_Guide/Brick_Configuration.html
//

// deviceInfo represents structure in which devices are to be store in Peer Metadata
type deviceInfo struct {
	Name          string `json:"name"`
	State         string `json:"state"`
	VgName        string `json:"vg-name"`
	AvailableSize uint64 `json:"available-size"`
	ExtentSize    uint64 `json:"extent-size"`
	Used          bool   `json:"used"`
}

func (g *executor) DeviceSetup(host, device, vgid string, destroy bool) (*executors.DeviceInfo, error) {
	logger.Debug("BEGIN")
	logger.Debug("host: %s - port: %s", host, g.config.ClientPort)
	g.createClient(host)
	peerid := ""
	peerlist, err := g.client.Peers()
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	hostPort := host + ":" + g.config.ClientPort
	for _, peer := range peerlist {
		for _, addr := range peer.ClientAddresses {
			if addr == hostPort {
				peerid = peer.ID.String()
				break
			}
		}
	}
	//TODO implement device delete in failed scenario?
	if peerid == "" {
		msg := fmt.Sprintf("host %s is none of the peers", hostPort)
		logger.LogError(msg)
		return nil, errors.New(msg)
	}
	dev, err := g.client.DeviceAdd(peerid, device)
	if err != nil {
		logger.LogError("DeviceAdd(%s, %s) failed: %v", peerid, device, err)
		return nil, err
	}
	logger.Debug("DeviceAdd(%s, %s) successful", peerid, device)
	var devices []deviceInfo
	var info = new(executors.DeviceInfo)
	if _, exist := dev.Metadata["_devices"]; exist {

		err = json.Unmarshal([]byte(dev.Metadata["_devices"]), &devices)
		if err != nil {
			logger.Err(err)
			return nil, err
		}
		for _, d := range devices {
			if device == d.Name {
				info.ExtentSize = d.ExtentSize
				info.TotalSize = d.AvailableSize
				info.FreeSize = d.AvailableSize
			}
		}
	}
	if info == nil {
		msg := "failed to fetch device details"
		logger.LogError(msg)
		return nil, errors.New(msg)
	}
	return info, nil
}

func (g *executor) GetDeviceInfo(host, device, vgid string) (d *executors.DeviceInfo, e error) {
	logger.Debug("BEGIN")
	//TODO need to replace this function by listing device
	d = &executors.DeviceInfo{}

	g.createClient(host)

	peerlist, err := g.client.Peers()
	if err != nil {
		logger.Err(err)
		return nil, err
	}
	var devices []deviceInfo
	var info = new(executors.DeviceInfo)
	for _, peer := range peerlist {
		for _, addr := range peer.PeerAddresses {
			if addr == host+g.config.ClientPort {
				if _, exist := peer.Metadata["_devices"]; exist {

					err = json.Unmarshal([]byte(peer.Metadata["_devices"]), &devices)
					if err != nil {
						logger.Err(err)
						return nil, err
					}
					for _, d := range devices {
						if device == d.Name {
							info.ExtentSize = d.ExtentSize
							info.TotalSize = d.AvailableSize
							break
						}
					}
				}
			}

		}
	}

	if info == nil {
		logger.LogError("failed to fetch device details")
		return nil, errors.New("failed to fetch device details")
	}
	return info, nil
}

func (g *executor) DeviceTeardown(host, device, vgid string) error {
	logger.Debug("BEGIN")
	//TODO need to implement this api
	return executors.NotSupportedError
}
