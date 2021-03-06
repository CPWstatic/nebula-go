/* Copyright (c) 2019 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License,
 * attached with Common Clause Condition 1.0, found in the LICENSES directory.
 */

package ngdb

import (
	"fmt"
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-go/graph"
	"time"
)

type GraphOptions struct {
	Timeout time.Duration
}

type GraphOption func(*GraphOptions)

var defaultGraphOptions = GraphOptions{
	Timeout: 30 * time.Second,
}

type GraphClient struct {
	graph  graph.GraphServiceClient
	option GraphOptions
}

func WithTimeout(duration time.Duration) GraphOption {
	return func(options *GraphOptions) {
		options.Timeout = duration
	}
}

func NewClient(address string, opts ...GraphOption) (client *GraphClient, err error) {
	options := defaultGraphOptions
	for _, opt := range opts {
		opt(&options)
	}

	timeoutOption := thrift.SocketTimeout(options.Timeout)
	addressOption := thrift.SocketAddr(address)
	transport, err := thrift.NewSocket(timeoutOption, addressOption)

	if err != nil {
		return nil, err
	}

	protocol := thrift.NewBinaryProtocolFactoryDefault()
	graph := &GraphClient{
		graph: *graph.NewGraphServiceClientFactory(transport, protocol),
	}
	return graph, nil
}

func (client *GraphClient) Open() error {
	err := client.graph.Transport.Open()
	if err != nil {
		return err
	}
	return nil
}

func (client *GraphClient) Close() {
	client.graph.Close()
}

func (client *GraphClient) Authenticate(username string, password string) (response *graph.AuthResponse, err error) {
	if response, err := client.graph.Authenticate(username, password); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (client *GraphClient) Signout(sessionID int64) {
	if err := client.graph.Signout(sessionID); err != nil {
		fmt.Printf("Signout Failed : %v", err)
	}
}

func (client *GraphClient) Execute(sessionID int64, stmt string) (response *graph.ExecutionResponse, err error) {
	if response, err := client.graph.Execute(sessionID, stmt); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
