// Copyright 2019 Kaleido

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	eventmocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/events/service/mocks"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/firefly-fabconnect/internal/conf"
	eventsapi "github.com/hyperledger/firefly-fabconnect/internal/events/api"
	"github.com/hyperledger/firefly-fabconnect/internal/fabric/utils"
	"github.com/hyperledger/firefly-fabconnect/internal/kvstore"
	mockfabric "github.com/hyperledger/firefly-fabconnect/mocks/fabric/client"
	mockkvstore "github.com/hyperledger/firefly-fabconnect/mocks/kvstore"
	"github.com/stretchr/testify/mock"
)

func tempdir(t *testing.T) string {
	dir, _ := ioutil.TempDir("", "fly")
	t.Logf("tmpdir/create: %s", dir)
	return dir
}

func cleanup(t *testing.T, dir string) {
	t.Logf("tmpdir/cleanup: %s [dir]", dir)
	os.RemoveAll(dir)
}

type mockWebSocket struct {
	capturedNamespace string
	sender            chan interface{}
	broadcast         chan interface{}
	receiver          chan error
	closing           chan struct{}
}

func (m *mockWebSocket) GetChannels(namespace string) (chan<- interface{}, chan<- interface{}, <-chan error, <-chan struct{}) {
	m.capturedNamespace = namespace
	return m.sender, m.broadcast, m.receiver, m.closing
}

func (m *mockWebSocket) SendReply(message interface{}) {}

func newMockWebSocket() *mockWebSocket {
	return &mockWebSocket{
		sender:    make(chan interface{}),
		broadcast: make(chan interface{}),
		receiver:  make(chan error),
		closing:   make(chan struct{}),
	}
}

type mockSubMgr struct {
	stream        *eventStream
	subscription  *subscription
	err           error
	subscriptions []*subscription
}

func (m *mockSubMgr) getConfig() *conf.EventstreamConf {
	return &conf.EventstreamConf{}
}

func (m *mockSubMgr) streamByID(string) (*eventStream, error) {
	return m.stream, m.err
}

func (m *mockSubMgr) subscriptionByID(string) (*subscription, error) {
	return m.subscription, m.err
}

func (m *mockSubMgr) subscriptionsForStream(string) []*subscription {
	return m.subscriptions
}

func (m *mockSubMgr) loadCheckpoint(string) (map[string]uint64, error) { return nil, nil }

func (m *mockSubMgr) storeCheckpoint(string, map[string]uint64) error { return nil }

func testSubInfo(name string) *eventsapi.SubscriptionInfo {
	return &eventsapi.SubscriptionInfo{ID: "test", Stream: "streamID", Name: name}
}

func newTestStream(submgr subscriptionManager) *eventStream {
	a, _ := newEventStream(submgr, &StreamInfo{
		ID:   "123",
		Type: "WebHook",
		Webhook: &webhookActionInfo{
			URL: "http://hello.example.com/world",
		},
	}, nil)
	return a
}

func newTestSubscriptionManager() *subscriptionMGR {
	smconf := &conf.EventstreamConf{}
	rpc := mockRPCClient("")
	sm := NewSubscriptionManager(smconf, rpc, newMockWebSocket()).(*subscriptionMGR)
	sm.db = &mockkvstore.KVStore{}
	sm.config.WebhooksAllowPrivateIPs = true
	sm.config.PollingIntervalSec = 0
	return sm
}

func newTestStreamForBatching(spec *StreamInfo, db kvstore.KVStore, status ...int) (*subscriptionMGR, *eventStream, *httptest.Server, chan []*eventsapi.EventEntry) {
	mux := http.NewServeMux()
	eventStream := make(chan []*eventsapi.EventEntry)
	count := 0
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		var events []*eventsapi.EventEntry
		_ = json.NewDecoder(req.Body).Decode(&events)
		eventStream <- events
		idx := count
		if idx >= len(status) {
			idx = len(status) - 1
		}
		res.WriteHeader(status[idx])
		count++
	})
	svr := httptest.NewServer(mux)
	if spec.Type == "" {
		spec.Type = "webhook"
		spec.Webhook.URL = svr.URL
		spec.Webhook.Headers = map[string]string{"x-my-header": "my-value"}
	}
	sm := newTestSubscriptionManager()
	sm.config.WebhooksAllowPrivateIPs = true
	sm.config.PollingIntervalSec = 0
	if db != nil {
		sm.db = db
	}
	mockstore, ok := sm.db.(*mockkvstore.KVStore)
	if ok {
		mockstore.On("Get", mock.Anything).Return([]byte(""), nil)
		mockstore.On("Put", mock.Anything, mock.Anything).Return(nil)
	}

	_ = sm.addStream(spec)
	return sm, sm.streams[spec.ID], svr, eventStream
}

func newTestStreamForWebSocket(spec *StreamInfo, db kvstore.KVStore, status ...int) (*subscriptionMGR, *eventStream, *mockWebSocket) {
	sm := newTestSubscriptionManager()
	sm.config.PollingIntervalSec = 0
	if db != nil {
		sm.db = db
	}
	_ = sm.addStream(spec)
	return sm, sm.streams[spec.ID], sm.wsChannels.(*mockWebSocket)
}

func testEvent(subID string) *eventData {
	entry := &eventsapi.EventEntry{
		SubID: subID,
	}
	return &eventData{
		event:         entry,
		batchComplete: func(*eventsapi.EventEntry) {},
	}
}

func mockRPCClient(fromBlock string, withReset ...bool) *mockfabric.RPCClient {
	rpc := &mockfabric.RPCClient{}
	blockEventChan := make(chan *fab.BlockEvent)
	ccEventChan := make(chan *fab.CCEvent)
	var roBlockEventChan <-chan *fab.BlockEvent = blockEventChan
	var roCCEventChan <-chan *fab.CCEvent = ccEventChan
	res := &fab.BlockchainInfoResponse{
		BCI: &common.BlockchainInfo{
			Height: 10,
		},
	}
	rawBlock := &utils.RawBlock{
		Header: &common.BlockHeader{
			Number: uint64(20),
		},
	}
	block := &utils.Block{
		Number:    uint64(20),
		Timestamp: int64(1000000),
	}
	rpc.On("SubscribeEvent", mock.Anything, mock.Anything).Return(nil, roBlockEventChan, roCCEventChan, nil)
	rpc.On("QueryChainInfo", mock.Anything, mock.Anything).Return(res, nil)
	rpc.On("QueryBlock", mock.Anything, mock.Anything, mock.Anything).Return(rawBlock, block, nil)
	rpc.On("Unregister", mock.Anything).Return()

	go func() {
		if fromBlock == "0" {
			blockEventChan <- &fab.BlockEvent{
				Block: constructBlock(1),
			}
		}
		blockEventChan <- &fab.BlockEvent{
			Block: constructBlock(11),
		}
		ccEventChan <- &fab.CCEvent{
			BlockNumber: uint64(10),
			TxID:        "3144a3ad43dcc11374832bbb71561320de81fd80d69cc8e26a9ea7d3240a5e84",
			ChaincodeID: "asset_transfer",
		}
		if len(withReset) > 0 {
			blockEventChan <- &fab.BlockEvent{
				Block: constructBlock(11),
			}
		}
	}()

	return rpc
}

func setupTestSubscription(sm *subscriptionMGR, stream *eventStream, subscriptionName, fromBlock string, withReset ...bool) *eventsapi.SubscriptionInfo {
	rpc := mockRPCClient(fromBlock, withReset...)
	sm.rpc = rpc
	spec := &eventsapi.SubscriptionInfo{
		Name:   subscriptionName,
		Stream: stream.spec.ID,
	}
	if fromBlock != "" {
		spec.FromBlock = fromBlock
	}
	_ = sm.addSubscription(spec)

	return spec
}

func constructBlock(number uint64) *common.Block {
	mockTx := eventmocks.NewTransactionWithCCEvent("testTxID", peer.TxValidationCode_VALID, "testChaincodeID", "testCCEventName", []byte("testPayload"))
	mockBlock := eventmocks.NewBlock("testChannelID", mockTx)
	mockBlock.Header.Number = number
	return mockBlock
}