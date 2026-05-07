// Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package testutil

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/kata-containers/kata-containers/src/runtime/virtcontainers/pkg/agent/protocols"
	pb "github.com/kata-containers/kata-containers/src/runtime/virtcontainers/pkg/agent/protocols/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockConn implements net.Conn for testing
type MockConn struct {
	net.Conn
	closed bool
	mu     sync.Mutex
}

func NewMockConn() *MockConn {
	return &MockConn{}
}

func (m *MockConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *MockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}

func (m *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9090}
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// MockAgentServiceClient implements pb.AgentServiceService for testing
// Only includes methods actually used in redirector_test.go
type MockAgentServiceClient struct {
	pb.AgentServiceService
	CreateContainerErr error
	StartContainerErr  error
	RemoveContainerErr error
	ExecProcessErr     error
	SignalProcessErr   error
	WaitProcessErr     error
}

func (m *MockAgentServiceClient) CreateContainer(ctx context.Context, req *pb.CreateContainerRequest) (*emptypb.Empty, error) {
	if m.CreateContainerErr != nil {
		return nil, m.CreateContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) StartContainer(ctx context.Context, req *pb.StartContainerRequest) (*emptypb.Empty, error) {
	if m.StartContainerErr != nil {
		return nil, m.StartContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) RemoveContainer(ctx context.Context, req *pb.RemoveContainerRequest) (*emptypb.Empty, error) {
	if m.RemoveContainerErr != nil {
		return nil, m.RemoveContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) ExecProcess(ctx context.Context, req *pb.ExecProcessRequest) (*emptypb.Empty, error) {
	if m.ExecProcessErr != nil {
		return nil, m.ExecProcessErr
	}
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) SignalProcess(ctx context.Context, req *pb.SignalProcessRequest) (*emptypb.Empty, error) {
	if m.SignalProcessErr != nil {
		return nil, m.SignalProcessErr
	}
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) WaitProcess(ctx context.Context, req *pb.WaitProcessRequest) (*pb.WaitProcessResponse, error) {
	if m.WaitProcessErr != nil {
		return nil, m.WaitProcessErr
	}
	return &pb.WaitProcessResponse{Status: 0}, nil
}

func (m *MockAgentServiceClient) WriteStdin(ctx context.Context, req *pb.WriteStreamRequest) (*pb.WriteStreamResponse, error) {
	return &pb.WriteStreamResponse{Len: uint32(len(req.Data))}, nil
}

func (m *MockAgentServiceClient) ReadStdout(ctx context.Context, req *pb.ReadStreamRequest) (*pb.ReadStreamResponse, error) {
	return &pb.ReadStreamResponse{Data: []byte("stdout data")}, nil
}

func (m *MockAgentServiceClient) CloseStdin(ctx context.Context, req *pb.CloseStdinRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) UpdateInterface(ctx context.Context, req *pb.UpdateInterfaceRequest) (*protocols.Interface, error) {
	return &protocols.Interface{}, nil
}

func (m *MockAgentServiceClient) UpdateRoutes(ctx context.Context, req *pb.UpdateRoutesRequest) (*pb.Routes, error) {
	return &pb.Routes{}, nil
}

func (m *MockAgentServiceClient) CreateSandbox(ctx context.Context, req *pb.CreateSandboxRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *MockAgentServiceClient) DestroySandbox(ctx context.Context, req *pb.DestroySandboxRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// MockHealthServiceClient implements pb.HealthService for testing
type MockHealthServiceClient struct {
	pb.HealthService
	CheckErr   error
	VersionErr error
}

func (m *MockHealthServiceClient) Check(ctx context.Context, req *pb.CheckRequest) (*pb.HealthCheckResponse, error) {
	if m.CheckErr != nil {
		return nil, m.CheckErr
	}
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (m *MockHealthServiceClient) Version(ctx context.Context, req *pb.CheckRequest) (*pb.VersionCheckResponse, error) {
	if m.VersionErr != nil {
		return nil, m.VersionErr
	}
	return &pb.VersionCheckResponse{
		GrpcVersion:  "1.0.0",
		AgentVersion: "2.0.0",
	}, nil
}
