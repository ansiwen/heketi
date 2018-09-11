//
// Copyright (c) 2015 The heketi Authors
//
// This file is licensed to you under your choice of the GNU Lesser
// General Public License, version 3 or any later version (LGPLv3 or
// later), or the GNU General Public License, version 2 (GPLv2), in all
// cases as published by the Free Software Foundation.
//

package glusterfs

import (
	stdErrors "errors"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrNoSpace          = stdErrors.New("No space")
	ErrFound            = stdErrors.New("Id already exists")
	ErrNotFound         = stdErrors.New("Id not found")
	ErrConflict         = stdErrors.New("The target exists, contains other items, or is in use.")
	ErrMaxBricks        = stdErrors.New("Maximum number of bricks reached.")
	ErrMinimumBrickSize = stdErrors.New("Minimum brick size limit reached.  Out of space.")
	ErrDbAccess         = stdErrors.New("Unable to access db")
	ErrAccessList       = stdErrors.New("Unable to access list")
	ErrKeyExists        = stdErrors.New("Key already exists in the database")
	ErrNoReplacement    = stdErrors.New("No Replacement was found for resource requested to be removed")
	ErrCloneBlockVol    = stdErrors.New("Cloning of block hosting volumes is not supported")

	// well known errors for cluster device source
	ErrEmptyCluster = stdErrors.New("No nodes in cluster")
	ErrNoStorage    = stdErrors.New("No online storage devices in cluster")

	// returned by code related to operations load
	ErrTooManyOperations = stdErrors.New("Server handling too many operations")
)

// IsRetry returns true if the error-generating operation should be retried.
func IsRetry(err error) bool {
	err = errors.Cause(err)
	te, ok := err.(interface {
		Retry() bool
	})
	return ok && te.Retry()
}

// Original returns a nested error if present or nil.
func Original(err error) error {
	err = errors.Cause(err)
	if ne, ok := err.(interface {
		Original() error
	}); ok {
		return ne.Original()
	}
	return nil
}

type retryError struct {
	originalError error
}

// NewRetryError wraps err in a retryError
func NewRetryError(err error) error {
	return retryError{originalError: err}
}

func (ore retryError) Error() string {
	return fmt.Sprintf("Operation Should Be Retried; Error: %v",
		ore.originalError.Error())
}

func (ore retryError) Original() error {
	return ore.originalError
}

func (ore retryError) Retry() bool {
	return true
}
