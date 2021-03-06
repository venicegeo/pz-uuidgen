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

package systest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/venicegeo/pz-gocommon/gocommon"
	"github.com/venicegeo/pz-uuidgen/uuidgen"
)

func sleep() {
	time.Sleep(1 * time.Second)
}

type UuidgenTester struct {
	suite.Suite
	client    *uuidgen.Client
	url       string
	apiKey    string
	apiServer string
}

func (suite *UuidgenTester) setupFixture() {
	t := suite.T()
	assert := assert.New(t)

	var err error

	suite.apiServer, err = piazza.GetApiServer()
	assert.NoError(err)

	suite.url, err = piazza.GetPiazzaServiceUrl(piazza.PzUuidgen)
	assert.NoError(err)

	suite.apiKey, err = piazza.GetApiKey(suite.apiServer)
	assert.NoError(err)

	client, err := uuidgen.NewClient(suite.url, suite.apiKey)
	assert.NoError(err)
	suite.client = client
}

func (suite *UuidgenTester) teardownFixture() {
}

func TestRunSuite(t *testing.T) {
	s := &UuidgenTester{}
	suite.Run(t, s)
}

func isValid(uuid string) bool {
	return piazza.ValidUuid(uuid)
}

func (suite *UuidgenTester) Test00Version() {
	t := suite.T()
	assert := assert.New(t)

	suite.setupFixture()
	defer suite.teardownFixture()

	client := suite.client

	version, err := client.GetVersion()
	assert.NoError(err)
	assert.EqualValues("1.0.0", version.Version)
}

func (suite *UuidgenTester) Test01Get() {
	t := suite.T()
	assert := assert.New(t)

	suite.setupFixture()
	defer suite.teardownFixture()

	client := suite.client

	uuid, err := client.GetUUID()
	assert.NoError(err)

	assert.True(isValid(uuid))
}

func (suite *UuidgenTester) Test02Post() {
	t := suite.T()
	assert := assert.New(t)

	suite.setupFixture()
	defer suite.teardownFixture()

	client := suite.client

	uuids, err := client.PostUuids(17)
	assert.NoError(err)
	assert.Len(*uuids, 17)

	for i := 0; i < 17; i++ {
		a := (*uuids)[i]
		assert.True(isValid(a))
		for j := i + 1; j < 17; j++ {
			b := (*uuids)[j]
			assert.NotEqual(a, b)
		}
	}
}

func (suite *UuidgenTester) Test03Admin() {
	t := suite.T()
	assert := assert.New(t)

	suite.setupFixture()
	defer suite.teardownFixture()

	client := suite.client

	stats, err := client.GetStats()
	assert.NoError(err, "GetFromAdminStats")

	assert.NotZero(stats.NumUUIDs)
	assert.NotZero(stats.NumRequests)
}
