// Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package agentproto

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"

	"github.com/confidential-containers/cloud-api-adaptor/src/cloud-api-adaptor/pkg/util/agentproto/testutil"
	pb "github.com/kata-containers/kata-containers/src/runtime/virtcontainers/pkg/agent/protocols/grpc"
)

// setupTestRedirector creates a test redirector with the given mock clients
func setupTestRedirector(t *testing.T, mockAgent *testutil.MockAgentServiceClient, mockHealth *testutil.MockHealthServiceClient) *redirector {
	t.Helper()

	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockAgent,
			HealthService:       mockHealth,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return testutil.NewMockConn(), nil
		},
	}
	r.once.Do(func() {}) // Pre-initialize to avoid connection logic
	return r
}

// TestNewRedirector tests the NewRedirector constructor
func TestNewRedirector(t *testing.T) {
	tests := []struct {
		name   string
		dialer func(context.Context) (net.Conn, error)
	}{
		{
			name: "valid dialer",
			dialer: func(ctx context.Context) (net.Conn, error) {
				return testutil.NewMockConn(), nil
			},
		},
		{
			name:   "nil dialer",
			dialer: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedirector(tt.dialer)
			if r == nil {
				t.Fatal("NewRedirector returned nil")
			}

			redirectorImpl, ok := r.(*redirector)
			if !ok {
				t.Fatal("NewRedirector did not return *redirector type")
			}

			if tt.dialer != nil && redirectorImpl.dialer == nil {
				t.Error("dialer was not set correctly")
			}
		})
	}
}

// TestConnect tests the Connect method with various scenarios
func TestConnect(t *testing.T) {
	tests := []struct {
		name      string
		dialer    func(context.Context) (net.Conn, error)
		cancelCtx bool
		wantErr   bool
	}{
		{
			name: "successful connection",
			dialer: func(ctx context.Context) (net.Conn, error) {
				return testutil.NewMockConn(), nil
			},
			wantErr: false,
		},
		{
			name: "connection failure",
			dialer: func(ctx context.Context) (net.Conn, error) {
				return nil, errors.New("connection failed")
			},
			wantErr: true,
		},
		{
			name: "context cancellation",
			dialer: func(ctx context.Context) (net.Conn, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
			cancelCtx: true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redirector{dialer: tt.dialer}

			ctx := context.Background()
			if tt.cancelCtx {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			err := r.Connect(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConnectConcurrency tests concurrent Connect calls
func TestConnectConcurrency(t *testing.T) {
	mockAgent := &testutil.MockAgentServiceClient{}
	r := setupTestRedirector(t, mockAgent, nil)

	ctx := context.Background()
	const numGoroutines = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.Connect(ctx); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("concurrent Connect failed: %v", err)
	}
}

// TestClose tests the Close method
func TestClose(t *testing.T) {
	mockAgent := &testutil.MockAgentServiceClient{}
	r := setupTestRedirector(t, mockAgent, nil)

	if err := r.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestContainerOperations tests critical container operations with error handling
func TestContainerOperations(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*testutil.MockAgentServiceClient)
		operation func(*redirector, context.Context) error
		wantErr   bool
	}{
		{
			name: "CreateContainer success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.CreateContainerErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.CreateContainer(ctx, &pb.CreateContainerRequest{})
				return err
			},
			wantErr: false,
		},
		{
			name: "CreateContainer failure",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.CreateContainerErr = errors.New("create failed")
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.CreateContainer(ctx, &pb.CreateContainerRequest{})
				return err
			},
			wantErr: true,
		},
		{
			name: "StartContainer success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.StartContainerErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.StartContainer(ctx, &pb.StartContainerRequest{})
				return err
			},
			wantErr: false,
		},
		{
			name: "RemoveContainer success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.RemoveContainerErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.RemoveContainer(ctx, &pb.RemoveContainerRequest{})
				return err
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAgent := &testutil.MockAgentServiceClient{}
			tt.setupMock(mockAgent)
			r := setupTestRedirector(t, mockAgent, nil)

			ctx := context.Background()
			err := tt.operation(r, ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

// TestProcessOperations tests critical process operations
func TestProcessOperations(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*testutil.MockAgentServiceClient)
		operation func(*redirector, context.Context) error
		wantErr   bool
	}{
		{
			name: "ExecProcess success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.ExecProcessErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.ExecProcess(ctx, &pb.ExecProcessRequest{})
				return err
			},
			wantErr: false,
		},
		{
			name: "SignalProcess success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.SignalProcessErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				_, err := r.SignalProcess(ctx, &pb.SignalProcessRequest{})
				return err
			},
			wantErr: false,
		},
		{
			name: "WaitProcess success",
			setupMock: func(m *testutil.MockAgentServiceClient) {
				m.WaitProcessErr = nil
			},
			operation: func(r *redirector, ctx context.Context) error {
				resp, err := r.WaitProcess(ctx, &pb.WaitProcessRequest{})
				if err == nil && resp == nil {
					return errors.New("nil response")
				}
				return err
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAgent := &testutil.MockAgentServiceClient{}
			tt.setupMock(mockAgent)
			r := setupTestRedirector(t, mockAgent, nil)

			ctx := context.Background()
			err := tt.operation(r, ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

// TestHealthService tests health check functionality
func TestHealthService(t *testing.T) {
	tests := []struct {
		name       string
		checkErr   error
		versionErr error
	}{
		{
			name:       "successful health check",
			checkErr:   nil,
			versionErr: nil,
		},
		{
			name:       "health check failure",
			checkErr:   errors.New("health check failed"),
			versionErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAgent := &testutil.MockAgentServiceClient{}
			mockHealth := &testutil.MockHealthServiceClient{
				CheckErr:   tt.checkErr,
				VersionErr: tt.versionErr,
			}
			r := setupTestRedirector(t, mockAgent, mockHealth)
			ctx := context.Background()

			_, err := r.Check(ctx, &pb.CheckRequest{})
			if (err != nil) != (tt.checkErr != nil) {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.checkErr != nil)
			}

			_, err = r.Version(ctx, &pb.CheckRequest{})
			if (err != nil) != (tt.versionErr != nil) {
				t.Errorf("Version() error = %v, wantErr %v", err, tt.versionErr != nil)
			}
		})
	}
}

// TestStreamOperations tests I/O stream operations
func TestStreamOperations(t *testing.T) {
	mockAgent := &testutil.MockAgentServiceClient{}
	r := setupTestRedirector(t, mockAgent, nil)
	ctx := context.Background()

	// Test WriteStdin
	resp, err := r.WriteStdin(ctx, &pb.WriteStreamRequest{Data: []byte("test")})
	if err != nil || resp == nil {
		t.Errorf("WriteStdin() error = %v, resp = %v", err, resp)
	}

	// Test ReadStdout
	resp2, err := r.ReadStdout(ctx, &pb.ReadStreamRequest{})
	if err != nil || resp2 == nil {
		t.Errorf("ReadStdout() error = %v, resp = %v", err, resp2)
	}

	// Test CloseStdin
	_, err = r.CloseStdin(ctx, &pb.CloseStdinRequest{})
	if err != nil {
		t.Errorf("CloseStdin() error = %v", err)
	}
}

// TestSandboxOperations tests sandbox lifecycle
func TestSandboxOperations(t *testing.T) {
	mockAgent := &testutil.MockAgentServiceClient{}
	r := setupTestRedirector(t, mockAgent, nil)
	ctx := context.Background()

	// Test CreateSandbox
	_, err := r.CreateSandbox(ctx, &pb.CreateSandboxRequest{})
	if err != nil {
		t.Errorf("CreateSandbox() error = %v", err)
	}

	// Test DestroySandbox
	_, err = r.DestroySandbox(ctx, &pb.DestroySandboxRequest{})
	if err != nil {
		t.Errorf("DestroySandbox() error = %v", err)
	}
}

// TestNetworkOperations tests network configuration
func TestNetworkOperations(t *testing.T) {
	mockAgent := &testutil.MockAgentServiceClient{}
	r := setupTestRedirector(t, mockAgent, nil)
	ctx := context.Background()

	// Test UpdateInterface
	_, err := r.UpdateInterface(ctx, &pb.UpdateInterfaceRequest{})
	if err != nil {
		t.Errorf("UpdateInterface() error = %v", err)
	}

	// Test UpdateRoutes
	_, err = r.UpdateRoutes(ctx, &pb.UpdateRoutesRequest{})
	if err != nil {
		t.Errorf("UpdateRoutes() error = %v", err)
	}
}

// TestInterfaceCompliance verifies interface implementation
func TestInterfaceCompliance(t *testing.T) {
	var _ Redirector = (*redirector)(nil)
	var _ pb.AgentServiceService = (*client)(nil)
	var _ pb.HealthService = (*client)(nil)
}
