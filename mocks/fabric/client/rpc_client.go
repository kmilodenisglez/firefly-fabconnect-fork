// Code generated by mockery v1.0.0. DO NOT EDIT.

package mockfabric

import (
	api "github.com/hyperledger/firefly-fabconnect/internal/events/api"
	client "github.com/hyperledger/firefly-fabconnect/internal/fabric/client"

	fab "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	mock "github.com/stretchr/testify/mock"

	utils "github.com/hyperledger/firefly-fabconnect/internal/fabric/utils"
)

// RPCClient is an autogenerated mock type for the RPCClient type
type RPCClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *RPCClient) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Invoke provides a mock function with given fields: channelId, signer, chaincodeName, method, args, isInit
func (_m *RPCClient) Invoke(channelId string, signer string, chaincodeName string, method string, args []string, isInit bool) (*client.TxReceipt, error) {
	ret := _m.Called(channelId, signer, chaincodeName, method, args, isInit)

	var r0 *client.TxReceipt
	if rf, ok := ret.Get(0).(func(string, string, string, string, []string, bool) *client.TxReceipt); ok {
		r0 = rf(channelId, signer, chaincodeName, method, args, isInit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.TxReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, []string, bool) error); ok {
		r1 = rf(channelId, signer, chaincodeName, method, args, isInit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: channelId, signer, chaincodeName, method, args, strongread
func (_m *RPCClient) Query(channelId string, signer string, chaincodeName string, method string, args []string, strongread bool) ([]byte, error) {
	ret := _m.Called(channelId, signer, chaincodeName, method, args, strongread)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, string, string, []string, bool) []byte); ok {
		r0 = rf(channelId, signer, chaincodeName, method, args, strongread)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, []string, bool) error); ok {
		r1 = rf(channelId, signer, chaincodeName, method, args, strongread)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryBlock provides a mock function with given fields: channelId, blockNumber, signer
func (_m *RPCClient) QueryBlock(channelId string, blockNumber uint64, signer string) (*utils.RawBlock, *utils.Block, error) {
	ret := _m.Called(channelId, blockNumber, signer)

	var r0 *utils.RawBlock
	if rf, ok := ret.Get(0).(func(string, uint64, string) *utils.RawBlock); ok {
		r0 = rf(channelId, blockNumber, signer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*utils.RawBlock)
		}
	}

	var r1 *utils.Block
	if rf, ok := ret.Get(1).(func(string, uint64, string) *utils.Block); ok {
		r1 = rf(channelId, blockNumber, signer)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.Block)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, uint64, string) error); ok {
		r2 = rf(channelId, blockNumber, signer)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryChainInfo provides a mock function with given fields: channelId, signer
func (_m *RPCClient) QueryChainInfo(channelId string, signer string) (*fab.BlockchainInfoResponse, error) {
	ret := _m.Called(channelId, signer)

	var r0 *fab.BlockchainInfoResponse
	if rf, ok := ret.Get(0).(func(string, string) *fab.BlockchainInfoResponse); ok {
		r0 = rf(channelId, signer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fab.BlockchainInfoResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(channelId, signer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryTransaction provides a mock function with given fields: channelId, signer, txId
func (_m *RPCClient) QueryTransaction(channelId string, signer string, txId string) (map[string]interface{}, error) {
	ret := _m.Called(channelId, signer, txId)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(string, string, string) map[string]interface{}); ok {
		r0 = rf(channelId, signer, txId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(channelId, signer, txId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeEvent provides a mock function with given fields: subInfo, since
func (_m *RPCClient) SubscribeEvent(subInfo *api.SubscriptionInfo, since uint64) (*client.RegistrationWrapper, <-chan *fab.BlockEvent, <-chan *fab.CCEvent, error) {
	ret := _m.Called(subInfo, since)

	var r0 *client.RegistrationWrapper
	if rf, ok := ret.Get(0).(func(*api.SubscriptionInfo, uint64) *client.RegistrationWrapper); ok {
		r0 = rf(subInfo, since)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.RegistrationWrapper)
		}
	}

	var r1 <-chan *fab.BlockEvent
	if rf, ok := ret.Get(1).(func(*api.SubscriptionInfo, uint64) <-chan *fab.BlockEvent); ok {
		r1 = rf(subInfo, since)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan *fab.BlockEvent)
		}
	}

	var r2 <-chan *fab.CCEvent
	if rf, ok := ret.Get(2).(func(*api.SubscriptionInfo, uint64) <-chan *fab.CCEvent); ok {
		r2 = rf(subInfo, since)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(<-chan *fab.CCEvent)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(*api.SubscriptionInfo, uint64) error); ok {
		r3 = rf(subInfo, since)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Unregister provides a mock function with given fields: _a0
func (_m *RPCClient) Unregister(_a0 *client.RegistrationWrapper) {
	_m.Called(_a0)
}
