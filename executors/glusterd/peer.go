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
	"errors"

	"github.com/gluster/glusterd2/pkg/api"
	"github.com/lpabon/godbc"
)

//TODO need to hanlde glusterd2 port also
//currently its hardcoded

// :TODO: Rename this function to NodeInit or something
func (g *executor) PeerProbe(host, newnode string) error {
	logger.Debug("BEGIN")

	godbc.Require(host != "")
	godbc.Require(newnode != "")
	g.createClient(host)
	logger.Info("Probing: %v -> %v", host, newnode)
	// create the commands
	peerAddReq := api.PeerAddReq{
		Addresses: []string{newnode + ":" + g.config.ClientPort},
	}
	_, err := g.client.PeerAdd(peerAddReq)
	if err != nil {
		return err
	}

	//TODO need to check this

	// Determine if there is a snapshot limit configuration setting
	// if s.RemoteExecutor.SnapShotLimit() > 0 {
	// 	logger.Info("Setting snapshot limit")
	// 	commands = []string{
	// 		fmt.Sprintf("gluster --mode=script snapshot config snap-max-hard-limit %v",
	// 			s.RemoteExecutor.SnapShotLimit()),
	// 	}
	// 	_, err := s.RemoteExecutor.RemoteCommandExecute(host, commands, 10)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (g *executor) PeerDetach(host, detachnode string) error {
	logger.Debug("BEGIN")
	godbc.Require(host != "")
	godbc.Require(detachnode != "")
	g.createClient(host)

	var peerid string
	// create the commands
	logger.Info("Detaching node %v", detachnode)
	//list nodes in cluster
	peerlist, err := g.client.Peers()
	if err != nil {
		logger.Err(err)
		return err
	}
	for _, peer := range peerlist {
		for _, addr := range peer.PeerAddresses {
			if addr == detachnode+g.config.ClientPort {
				peerid = peer.ID.String()
			}
		}
	}
	if peerid == "" {
		logger.LogError("not able to find peer info for %s", detachnode)
		return errors.New("unable to detatch peer")
	}
	err = g.client.PeerRemove(peerid)
	if err != nil {
		logger.Err(err)
	}

	return nil
}

func (g *executor) GlusterdCheck(host string) error {
	logger.Debug("BEGIN")
	godbc.Require(host != "")

	logger.Info("Check Glusterd service status in node %v", host)
	g.createClient(host)
	err := g.client.Ping()
	if err != nil {
		logger.Err(err)
		return err
	}

	return nil
}
