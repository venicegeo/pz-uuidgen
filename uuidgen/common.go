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

import "time"

type IClient interface {
	// high-level interfaces
	GetUuid() (string, error)

	// low-level interfaces
	PostUuids(count int) (*[]string, error)
	GetStats() (*UuidGenAdminStats, error)
}

type UuidGenAdminStats struct {
	NumUUIDs    int       `json:"num_uuids"`
	NumRequests int       `json:"num_requests"`
	CreatedOn   time.Time `json:"createdOn"`
}
