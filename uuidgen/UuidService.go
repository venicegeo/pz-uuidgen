// Copyright 2016, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uuidgen

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pborman/uuid"
	piazza "github.com/venicegeo/pz-gocommon/gocommon"
	pzlogger "github.com/venicegeo/pz-logger/logger"
)

//---------------------------------------------------------------------

type LockedAdminStats struct {
	sync.Mutex
	UuidGenAdminStats
}

type UuidService struct {
	logger pzlogger.IClient
	stats  LockedAdminStats
}

//---------------------------------------------------------------------

func (service *UuidService) Init(logger pzlogger.IClient) error {
	service.logger = logger
	service.stats.CreatedOn = time.Now()

	err := service.logger.Info("uuidgen started")
	if err != nil {
		return err
	}

	return nil
}

func (service *UuidService) GetAdminStats() *piazza.JsonResponse {
	service.stats.Lock()
	t := service.stats.UuidGenAdminStats
	service.stats.Unlock()

	resp := &piazza.JsonResponse{StatusCode: http.StatusOK, Data: t}
	err := resp.SetType()
	if err != nil {
		return &piazza.JsonResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}
	return resp
}

// request body is ignored
// we allow a count of zero, for testing
func (service *UuidService) PostUuids(params *piazza.HttpQueryParams) *piazza.JsonResponse {
	var count int
	var err error
	var key string

	// ?count=INT
	key = params.Get("count")

	if key == "" {
		count = 1
	} else {
		count, err = strconv.Atoi(key)
		if err != nil {
			s := fmt.Sprintf("query argument invalid: %s", key)
			return &piazza.JsonResponse{StatusCode: http.StatusBadRequest, Message: s}
		}
	}

	if count < 0 || count > 255 {
		s := fmt.Sprintf("query argument out of range: %d", count)
		return &piazza.JsonResponse{StatusCode: http.StatusBadRequest, Message: s}
	}

	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		uuids[i] = uuid.New()
	}

	service.stats.Lock()
	service.stats.NumUUIDs += count
	service.stats.NumRequests++
	service.stats.Unlock()

	//log.Printf("INFO: uuidgen created %d", count)

	resp := &piazza.JsonResponse{StatusCode: http.StatusCreated, Data: uuids}
	err = resp.SetType()
	if err != nil {
		log.Printf("UuidService.PostUuids: returning %#v", nil)
		return &piazza.JsonResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	log.Printf("UuidService.PostUuids: returning %#v", resp)
	return resp
}
