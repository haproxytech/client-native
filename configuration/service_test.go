// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

const (
	baseSlots      = 20
	slotsIncrement = 10
)

func TestService(t *testing.T) {
	serviceSuit := ServiceTestSuit{
		client: client,
	}
	suite.Run(t, &serviceSuit)
}

type ServiceTestSuit struct {
	suite.Suite
	client        *Client
	transactionID string
}

func (s *ServiceTestSuit) testCreate() {
	serviceName := misc.RandomString(10)
	service, err := s.client.NewService(serviceName, ScalingParams{})
	s.NotNil(service)
	s.Nil(err)
	s.client.DeleteService(serviceName)
}

func (s *ServiceTestSuit) SetupSuite() {
	ver, err := s.client.GetVersion("")
	s.Nil(err)
	transaction, err := s.client.StartTransaction(ver)
	s.Nil(err)
	s.transactionID = transaction.ID
}

func (s *ServiceTestSuit) TearDownSuite() {
	//nolint
	s.client.DeleteTransaction(s.transactionID)
}

func (s *ServiceTestSuit) TestService() {
	s.Run("create service", s.testCreate)
	suite.Run(s.T(), &ServiceInitiationSuite{
		client:        s.client,
		transactionID: s.transactionID,
	})
	suite.Run(s.T(), &ServiceUpdateSuit{
		client:        s.client,
		transactionID: s.transactionID,
	})
}

type ServiceInitiationSuite struct {
	suite.Suite
	client        *Client
	service       *Service
	serviceName   string
	transactionID string
}

func (s *ServiceInitiationSuite) SetupSuite() {
	s.serviceName = misc.RandomString(10)
}

func (s *ServiceInitiationSuite) BeforeTest(suiteName, testName string) {
	service, err := s.client.NewService(s.serviceName, ScalingParams{
		BaseSlots:       baseSlots,
		SlotsGrowthType: ServiceGrowthTypeLinear,
		SlotsIncrement:  slotsIncrement,
	})
	s.NotNil(service)
	s.Nil(err)
	s.service = service
}

func (s *ServiceInitiationSuite) AfterTest(suiteName, testName string) {
	// run Init to set transactionID if test did not run it
	_, err := s.service.Init(s.transactionID)
	s.Nil(err)
	err = s.service.Delete()
	s.Nil(err)
	s.service = nil
}

func (s *ServiceInitiationSuite) TestInitAndDelete() {
	r, err := s.service.Init(s.transactionID)
	s.True(r)
	s.Nil(err)
	err = s.service.Delete()
	s.Nil(err)
}

func (s *ServiceInitiationSuite) TestFailCreatingIfExists() {
	_, err := s.client.NewService(s.serviceName, ScalingParams{})
	s.NotNil(err)
}

func (s *ServiceInitiationSuite) TestRecreatingAfterDeletion() {
	client.DeleteService(s.serviceName)
	service, err := s.client.NewService(s.serviceName, ScalingParams{})
	s.NotNil(service)
	s.Nil(err)
}

func (s *ServiceInitiationSuite) TestCreateNewBackendAndServers() {
	r, err := s.service.Init(s.transactionID)
	s.True(r)
	s.Nil(err)
	_, backned, err := s.client.GetBackend(s.serviceName, s.transactionID)
	s.Nil(err)
	s.NotNil(backned)
	servers, err := s.service.GetServers()
	s.Nil(err)
	s.NotNil(servers)
	s.Equal(baseSlots, len(servers))
}

func (s *ServiceInitiationSuite) TestLoadExistingBackend() {
	serverPort := int64(81)
	servers := models.Servers{
		{Name: "s1", Address: "127.1.1.1", Port: &serverPort},
		{Name: "s2", Address: "127.1.1.2", Port: &serverPort},
		{Name: "s3", Address: "127.1.1.3", Port: &serverPort},
		{Name: "s4", Address: "127.1.1.4", Port: &serverPort},
	}
	s.Nil(s.createExistingService(servers))

	r, err := s.service.Init(s.transactionID)
	// Only existing data for was loaded
	// No modifications on the config have been done as the server count matches base slots value so reload should be false
	s.False(r)
	s.Nil(err)
	cServers, err := s.service.GetServers()
	s.Nil(err)
	s.Equal(baseSlots, len(s.service.nodes))

	for i, server := range cServers {
		if i < len(servers) {
			s.Equal(server.Name, servers[i].Name)
			s.Equal(server.Address, servers[i].Address)
			s.Equal(*server.Port, *servers[i].Port)
		} else {
			s.Equal(server.Address, "127.0.0.1")
			s.Equal(*server.Port, int64(80))
		}
	}
}

func (s *ServiceInitiationSuite) createExistingService(servers models.Servers) error {
	err := client.CreateBackend(&models.Backend{
		Name: s.serviceName,
	}, s.transactionID, 0)
	if err != nil {
		return err
	}
	for _, server := range servers {
		err := s.client.CreateServer(s.serviceName, server, s.transactionID, 0)
		if err != nil {
			return err
		}
	}

	defaultPort := int64(80)
	maintServer := &models.Server{
		Address:     "127.0.0.1",
		Port:        &defaultPort,
		Maintenance: "enabled",
	}

	for i := len(servers); i < baseSlots; i++ {
		maintServer.Name = fmt.Sprintf("s%d", i+1)
		err := s.client.CreateServer(s.serviceName, maintServer, s.transactionID, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

type ServiceUpdateSuit struct {
	suite.Suite
	client        *Client
	serviceName   string
	transactionID string
	service       *Service
}

func (s *ServiceUpdateSuit) SetupSuite() {
	s.serviceName = misc.RandomString(10)
}

func (s *ServiceUpdateSuit) BeforeTest(suiteName, testName string) {
	service, err := s.client.NewService(s.serviceName, ScalingParams{
		BaseSlots:       baseSlots,
		SlotsGrowthType: ServiceGrowthTypeLinear,
		SlotsIncrement:  slotsIncrement,
	})
	s.NotNil(service)
	s.Nil(err)
	r, err := service.Init(s.transactionID)
	s.True(r)
	s.Nil(err)
	s.service = service
}

func (s *ServiceUpdateSuit) AfterTest(suiteName, testName string) {
	err := s.service.Delete()
	s.Nil(err)
	s.service = nil
}

func (s *ServiceUpdateSuit) TestFirstUpdate() {
	servers := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.2", Port: 82},
		{Address: "127.1.1.3", Port: 83},
		{Address: "127.1.1.4", Port: 84},
	}
	r, err := s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	s.validateUpdateResult(servers)
}

func (s *ServiceUpdateSuit) validateUpdateResult(expected []ServiceServer) {
	servers, err := s.service.GetServers()
	s.Nil(err)
	for i, server := range servers {
		if i < len(expected) {
			s.Equal(expected[i].Address, server.Address)
			s.Equal(expected[i].Port, int(*server.Port))
			s.NotEqual("enabled", server.Maintenance)
		} else {
			s.Equal("127.0.0.1", server.Address)
			s.Equal(int64(80), *server.Port)
			s.Equal("enabled", server.Maintenance)
		}
	}
}

func (s *ServiceUpdateSuit) TestSecondUpdateWithDeletedServer() {
	servers := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.2", Port: 82},
		{Address: "127.1.1.3", Port: 83},
		{Address: "127.1.1.4", Port: 84},
	}
	r, err := s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	servers = []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.3", Port: 83},
		{Address: "127.1.1.4", Port: 84},
	}
	r, err = s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	// When the new server list has less servers than the previous one "holes" can be left
	// in the middle of enabled servers.
	// In this case server 127.1.1.2 got removed and the last enabled server from the remaining ones
	// gets moved into its slot, in this case 127.1.1.4
	expected := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.4", Port: 84},
		{Address: "127.1.1.3", Port: 83},
	}
	s.validateUpdateResult(expected)
}

func (s *ServiceUpdateSuit) TestSecondUpdateWithNewAndRemovedServers() {
	servers := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.2", Port: 82},
		{Address: "127.1.1.3", Port: 83},
		{Address: "127.1.1.4", Port: 84},
	}
	r, err := s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	servers = []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.4", Port: 84},
		{Address: "127.1.1.5", Port: 85},
		{Address: "127.1.1.6", Port: 86},
		{Address: "127.1.1.7", Port: 87},
		{Address: "127.1.1.2", Port: 82},
	}
	r, err = s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	// Server 127.1.1.3 is the only that will be marked as deleted
	// as server 127.1.1.2 reapears at the end of the server list with the same port.
	// The first new server will be placed in place of server 127.1.1.3 which is 127.1.1.5.
	// The remaining new servers will be added to the end.
	expected := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.2", Port: 82},
		{Address: "127.1.1.5", Port: 85},
		{Address: "127.1.1.4", Port: 84},
		{Address: "127.1.1.6", Port: 86},
		{Address: "127.1.1.7", Port: 87},
	}
	s.validateUpdateResult(expected)
}

func (s *ServiceUpdateSuit) TestSecondUpdateWithNoChanges() {
	servers := []ServiceServer{
		{Address: "127.1.1.1", Port: 81},
		{Address: "127.1.1.2", Port: 82},
		{Address: "127.1.1.3", Port: 83},
		{Address: "127.1.1.4", Port: 84},
	}
	r, err := s.service.Update(servers)
	s.True(r)
	s.Nil(err)
	r, err = s.service.Update(servers)
	s.False(r)
	s.Nil(err)
}

func (s *ServiceUpdateSuit) TestUpdateScalingParams() {
	newBaseSlots := 30
	err := s.service.UpdateScalingParams(ScalingParams{
		BaseSlots:       newBaseSlots,
		SlotsGrowthType: ServiceGrowthTypeLinear,
		SlotsIncrement:  slotsIncrement,
	})
	s.Nil(err)
	servers, err := s.service.GetServers()
	s.Nil(err)
	s.Equal(newBaseSlots, len(servers))
}

func (s *ServiceUpdateSuit) TestLinearUpscaling() {
	expectedSlotsCount := baseSlots + slotsIncrement
	expected := s.generateServers(baseSlots + 2)
	r, err := s.service.Update(expected)
	s.True(r)
	s.Nil(err)
	servers, err := s.service.GetServers()
	s.Nil(err)
	s.Equal(expectedSlotsCount, len(servers))
	s.validateUpdateResult(expected)
}

func (s *ServiceUpdateSuit) generateServers(count int) []ServiceServer {
	servers := make([]ServiceServer, 0)
	for i := 0; i < count; i++ {
		servers = append(servers, ServiceServer{
			Address: fmt.Sprintf("127.1.1.%d", i),
			Port:    80 + i,
		})
	}
	return servers
}

func (s *ServiceUpdateSuit) TestExponentialUpscaling() {
	expectedSlotsCount := baseSlots * 2
	// Switch from linear to exponential scaling
	err := s.service.UpdateScalingParams(ScalingParams{
		BaseSlots:       baseSlots,
		SlotsGrowthType: ServiceGrowthTypeExponential,
		SlotsIncrement:  slotsIncrement,
	})
	s.Nil(err)
	expected := s.generateServers(baseSlots + 2)
	r, err := s.service.Update(expected)
	s.True(r)
	s.Nil(err)
	servers, err := s.service.GetServers()
	s.Nil(err)
	s.Equal(expectedSlotsCount, len(servers))
	s.validateUpdateResult(expected)
}

func (s *ServiceUpdateSuit) TestLinearDownscaling() {
	// First we need to increase the server count to be one at least one threshold above our final expected value
	upscaleExpectedSlots := baseSlots + 2*slotsIncrement
	upscaleServerCount := baseSlots + slotsIncrement + 2
	s.scaleServiceAndValidate(upscaleExpectedSlots, upscaleServerCount)
	// Now we can update with the expected servers and validate if the service downscaled
	expectedSlotsCount := baseSlots + slotsIncrement
	serverCount := baseSlots + 2
	s.scaleServiceAndValidate(expectedSlotsCount, serverCount)
}

func (s *ServiceUpdateSuit) scaleServiceAndValidate(expectedSlots, serverCount int) {
	upscaleServers := s.generateServers(serverCount)
	r, err := s.service.Update(upscaleServers)
	s.True(r)
	s.Nil(err)
	servers, err := s.service.GetServers()
	s.Nil(err)
	s.Equal(expectedSlots, len(servers))
	s.validateUpdateResult(upscaleServers)
}

func (s *ServiceUpdateSuit) TestExponentialDownscaling() {
	// Switch from linear to exponential scaling
	err := s.service.UpdateScalingParams(ScalingParams{
		BaseSlots:       baseSlots,
		SlotsGrowthType: ServiceGrowthTypeExponential,
		SlotsIncrement:  slotsIncrement,
	})
	s.Nil(err)
	// First we need to increase the server count to be one at least one threshold above our final expected value
	upscaleExpectedSlots := baseSlots * 2 * 2
	upscaleServerCount := baseSlots*2 + 2
	s.scaleServiceAndValidate(upscaleExpectedSlots, upscaleServerCount)
	// Now we can update with the expected servers and validate if the service downscaled
	expectedSlotsCount := baseSlots * 2
	serverCount := baseSlots + 2
	s.scaleServiceAndValidate(expectedSlotsCount, serverCount)
}
