package agentproto

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/kata-containers/kata-containers/src/runtime/virtcontainers/pkg/agent/protocols"
	pb "github.com/kata-containers/kata-containers/src/runtime/virtcontainers/pkg/agent/protocols/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// mockConn implements net.Conn for testing
type mockConn struct {
	net.Conn
	closed bool
	mu     sync.Mutex
}

func (m *mockConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9090}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// mockAgentServiceClient implements pb.AgentServiceService for testing
type mockAgentServiceClient struct {
	pb.AgentServiceService
	createContainerErr error
	startContainerErr  error
	removeContainerErr error
	execProcessErr     error
	signalProcessErr   error
	waitProcessErr     error
	updateContainerErr error
}

func (m *mockAgentServiceClient) CreateContainer(ctx context.Context, req *pb.CreateContainerRequest) (*emptypb.Empty, error) {
	if m.createContainerErr != nil {
		return nil, m.createContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) StartContainer(ctx context.Context, req *pb.StartContainerRequest) (*emptypb.Empty, error) {
	if m.startContainerErr != nil {
		return nil, m.startContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) RemoveContainer(ctx context.Context, req *pb.RemoveContainerRequest) (*emptypb.Empty, error) {
	if m.removeContainerErr != nil {
		return nil, m.removeContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) ExecProcess(ctx context.Context, req *pb.ExecProcessRequest) (*emptypb.Empty, error) {
	if m.execProcessErr != nil {
		return nil, m.execProcessErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) SignalProcess(ctx context.Context, req *pb.SignalProcessRequest) (*emptypb.Empty, error) {
	if m.signalProcessErr != nil {
		return nil, m.signalProcessErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) WaitProcess(ctx context.Context, req *pb.WaitProcessRequest) (*pb.WaitProcessResponse, error) {
	if m.waitProcessErr != nil {
		return nil, m.waitProcessErr
	}
	return &pb.WaitProcessResponse{Status: 0}, nil
}

func (m *mockAgentServiceClient) UpdateContainer(ctx context.Context, req *pb.UpdateContainerRequest) (*emptypb.Empty, error) {
	if m.updateContainerErr != nil {
		return nil, m.updateContainerErr
	}
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) UpdateEphemeralMounts(ctx context.Context, req *pb.UpdateEphemeralMountsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) StatsContainer(ctx context.Context, req *pb.StatsContainerRequest) (*pb.StatsContainerResponse, error) {
	return &pb.StatsContainerResponse{}, nil
}

func (m *mockAgentServiceClient) PauseContainer(ctx context.Context, req *pb.PauseContainerRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) ResumeContainer(ctx context.Context, req *pb.ResumeContainerRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) RemoveStaleVirtiofsShareMounts(ctx context.Context, req *pb.RemoveStaleVirtiofsShareMountsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) WriteStdin(ctx context.Context, req *pb.WriteStreamRequest) (*pb.WriteStreamResponse, error) {
	return &pb.WriteStreamResponse{Len: uint32(len(req.Data))}, nil
}

func (m *mockAgentServiceClient) ReadStdout(ctx context.Context, req *pb.ReadStreamRequest) (*pb.ReadStreamResponse, error) {
	return &pb.ReadStreamResponse{Data: []byte("stdout data")}, nil
}

func (m *mockAgentServiceClient) ReadStderr(ctx context.Context, req *pb.ReadStreamRequest) (*pb.ReadStreamResponse, error) {
	return &pb.ReadStreamResponse{Data: []byte("stderr data")}, nil
}

func (m *mockAgentServiceClient) CloseStdin(ctx context.Context, req *pb.CloseStdinRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) TtyWinResize(ctx context.Context, req *pb.TtyWinResizeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) UpdateInterface(ctx context.Context, req *pb.UpdateInterfaceRequest) (*protocols.Interface, error) {
	return &protocols.Interface{}, nil
}

func (m *mockAgentServiceClient) UpdateRoutes(ctx context.Context, req *pb.UpdateRoutesRequest) (*pb.Routes, error) {
	return &pb.Routes{}, nil
}

func (m *mockAgentServiceClient) ListInterfaces(ctx context.Context, req *pb.ListInterfacesRequest) (*pb.Interfaces, error) {
	return &pb.Interfaces{}, nil
}

func (m *mockAgentServiceClient) ListRoutes(ctx context.Context, req *pb.ListRoutesRequest) (*pb.Routes, error) {
	return &pb.Routes{}, nil
}

func (m *mockAgentServiceClient) AddARPNeighbors(ctx context.Context, req *pb.AddARPNeighborsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) GetIPTables(ctx context.Context, req *pb.GetIPTablesRequest) (*pb.GetIPTablesResponse, error) {
	return &pb.GetIPTablesResponse{}, nil
}

func (m *mockAgentServiceClient) SetIPTables(ctx context.Context, req *pb.SetIPTablesRequest) (*pb.SetIPTablesResponse, error) {
	return &pb.SetIPTablesResponse{}, nil
}

func (m *mockAgentServiceClient) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.Metrics, error) {
	return &pb.Metrics{}, nil
}

func (m *mockAgentServiceClient) MemAgentMemcgSet(ctx context.Context, req *pb.MemAgentMemcgConfig) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) MemAgentCompactSet(ctx context.Context, req *pb.MemAgentCompactConfig) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) CreateSandbox(ctx context.Context, req *pb.CreateSandboxRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) DestroySandbox(ctx context.Context, req *pb.DestroySandboxRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) OnlineCPUMem(ctx context.Context, req *pb.OnlineCPUMemRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) ReseedRandomDev(ctx context.Context, req *pb.ReseedRandomDevRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) GetGuestDetails(ctx context.Context, req *pb.GuestDetailsRequest) (*pb.GuestDetailsResponse, error) {
	return &pb.GuestDetailsResponse{}, nil
}

func (m *mockAgentServiceClient) MemHotplugByProbe(ctx context.Context, req *pb.MemHotplugByProbeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) SetGuestDateTime(ctx context.Context, req *pb.SetGuestDateTimeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) CopyFile(ctx context.Context, req *pb.CopyFileRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) GetOOMEvent(ctx context.Context, req *pb.GetOOMEventRequest) (*pb.OOMEvent, error) {
	return &pb.OOMEvent{}, nil
}

func (m *mockAgentServiceClient) AddSwap(ctx context.Context, req *pb.AddSwapRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) AddSwapPath(ctx context.Context, req *pb.AddSwapPathRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) GetVolumeStats(ctx context.Context, req *pb.VolumeStatsRequest) (*pb.VolumeStatsResponse, error) {
	return &pb.VolumeStatsResponse{}, nil
}

func (m *mockAgentServiceClient) ResizeVolume(ctx context.Context, req *pb.ResizeVolumeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) SetPolicy(ctx context.Context, req *pb.SetPolicyRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *mockAgentServiceClient) GetDiagnosticData(ctx context.Context, req *pb.GetDiagnosticDataRequest) (*pb.GetDiagnosticDataResponse, error) {
	return &pb.GetDiagnosticDataResponse{}, nil
}

// mockHealthClient implements pb.HealthService for testing
type mockHealthClient struct {
	pb.HealthService
	checkErr   error
	versionErr error
}

func (m *mockHealthClient) Check(ctx context.Context, req *pb.CheckRequest) (*pb.HealthCheckResponse, error) {
	if m.checkErr != nil {
		return nil, m.checkErr
	}
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (m *mockHealthClient) Version(ctx context.Context, req *pb.CheckRequest) (*pb.VersionCheckResponse, error) {
	if m.versionErr != nil {
		return nil, m.versionErr
	}
	return &pb.VersionCheckResponse{
		GrpcVersion:  "1.0.0",
		AgentVersion: "2.0.0",
	}, nil
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
				return &mockConn{}, nil
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

			if redirectorImpl.agentClient != nil {
				t.Error("agentClient should be nil initially")
			}

			if redirectorImpl.ttrpcClient != nil {
				t.Error("ttrpcClient should be nil initially")
			}
		})
	}
}

// TestConnect tests the Connect method
func TestConnect(t *testing.T) {
	tests := []struct {
		name        string
		dialer      func(context.Context) (net.Conn, error)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful connection",
			dialer: func(ctx context.Context) (net.Conn, error) {
				return &mockConn{}, nil
			},
			expectError: false,
		},
		{
			name: "dialer returns error",
			dialer: func(ctx context.Context) (net.Conn, error) {
				return nil, errors.New("connection failed")
			},
			expectError: true,
			errorMsg:    "agent connection is not established",
		},
		{
			name: "context cancelled",
			dialer: func(ctx context.Context) (net.Conn, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
			expectError: true,
			errorMsg:    "agent connection is not established",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedirector(tt.dialer)

			ctx := context.Background()
			if tt.name == "context cancelled" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			err := r.Connect(ctx)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestConnectIdempotency tests that Connect is idempotent (only connects once)
func TestConnectIdempotency(t *testing.T) {
	callCount := 0
	dialer := func(ctx context.Context) (net.Conn, error) {
		callCount++
		return &mockConn{}, nil
	}

	r := NewRedirector(dialer)
	ctx := context.Background()

	// Call Connect multiple times
	for i := 0; i < 5; i++ {
		err := r.Connect(ctx)
		if err != nil {
			t.Fatalf("Connect failed on iteration %d: %v", i, err)
		}
	}

	// Dialer should only be called once due to sync.Once
	if callCount != 1 {
		t.Errorf("expected dialer to be called once, but was called %d times", callCount)
	}
}

// TestConnectConcurrency tests concurrent calls to Connect
func TestConnectConcurrency(t *testing.T) {
	callCount := 0
	var mu sync.Mutex
	dialer := func(ctx context.Context) (net.Conn, error) {
		mu.Lock()
		callCount++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Simulate some work
		return &mockConn{}, nil
	}

	r := NewRedirector(dialer)
	ctx := context.Background()

	// Launch multiple goroutines calling Connect concurrently
	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			err := r.Connect(ctx)
			if err != nil {
				t.Errorf("Connect failed: %v", err)
			}
		}()
	}

	wg.Wait()

	// Dialer should only be called once even with concurrent calls
	mu.Lock()
	defer mu.Unlock()
	if callCount != 1 {
		t.Errorf("expected dialer to be called once, but was called %d times", callCount)
	}
}

// TestClose tests the Close method
func TestClose(t *testing.T) {
	tests := []struct {
		name            string
		setupRedirector func() Redirector
		expectError     bool
	}{
		{
			name: "close connected redirector",
			setupRedirector: func() Redirector {
				r := NewRedirector(func(ctx context.Context) (net.Conn, error) {
					return &mockConn{}, nil
				})
				_ = r.Connect(context.Background())
				return r
			},
			expectError: false,
		},
		{
			name: "close unconnected redirector",
			setupRedirector: func() Redirector {
				return NewRedirector(func(ctx context.Context) (net.Conn, error) {
					return &mockConn{}, nil
				})
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setupRedirector()
			err := r.Close()

			if tt.expectError && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestCreateContainer tests the CreateContainer method
func TestCreateContainer(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *mockAgentServiceClient
		expectError bool
	}{
		{
			name: "successful create container",
			setupMock: func() *mockAgentServiceClient {
				return &mockAgentServiceClient{}
			},
			expectError: false,
		},
		{
			name: "create container with error",
			setupMock: func() *mockAgentServiceClient {
				return &mockAgentServiceClient{
					createContainerErr: errors.New("create failed"),
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.setupMock()
			r := &redirector{
				agentClient: &client{
					AgentServiceService: mockClient,
				},
				dialer: func(ctx context.Context) (net.Conn, error) {
					return &mockConn{}, nil
				},
			}
			// Mark as already connected
			r.once.Do(func() {})

			ctx := context.Background()
			req := &pb.CreateContainerRequest{ContainerId: "test-container"}

			_, err := r.CreateContainer(ctx, req)

			if tt.expectError && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestCreateContainerWithoutConnection tests CreateContainer when not connected
func TestCreateContainerWithoutConnection(t *testing.T) {
	r := NewRedirector(func(ctx context.Context) (net.Conn, error) {
		return nil, errors.New("connection failed")
	})

	ctx := context.Background()
	req := &pb.CreateContainerRequest{ContainerId: "test-container"}

	_, err := r.CreateContainer(ctx, req)
	if err == nil {
		t.Fatal("expected error when not connected")
	}
}

// TestStartContainer tests the StartContainer method
func TestStartContainer(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *mockAgentServiceClient
		expectError bool
	}{
		{
			name: "successful start container",
			setupMock: func() *mockAgentServiceClient {
				return &mockAgentServiceClient{}
			},
			expectError: false,
		},
		{
			name: "start container with error",
			setupMock: func() *mockAgentServiceClient {
				return &mockAgentServiceClient{
					startContainerErr: errors.New("start failed"),
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.setupMock()
			r := &redirector{
				agentClient: &client{
					AgentServiceService: mockClient,
				},
				dialer: func(ctx context.Context) (net.Conn, error) {
					return &mockConn{}, nil
				},
			}
			r.once.Do(func() {})

			ctx := context.Background()
			req := &pb.StartContainerRequest{ContainerId: "test-container"}

			_, err := r.StartContainer(ctx, req)

			if tt.expectError && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestRemoveContainer tests the RemoveContainer method
func TestRemoveContainer(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()
	req := &pb.RemoveContainerRequest{ContainerId: "test-container"}

	_, err := r.RemoveContainer(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestExecProcess tests the ExecProcess method
func TestExecProcess(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()
	req := &pb.ExecProcessRequest{
		ContainerId: "test-container",
		ExecId:      "test-exec",
	}

	_, err := r.ExecProcess(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestSignalProcess tests the SignalProcess method
func TestSignalProcess(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()
	req := &pb.SignalProcessRequest{
		ContainerId: "test-container",
		ExecId:      "test-exec",
		Signal:      9, // SIGKILL
	}

	_, err := r.SignalProcess(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestWaitProcess tests the WaitProcess method
func TestWaitProcess(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()
	req := &pb.WaitProcessRequest{
		ContainerId: "test-container",
		ExecId:      "test-exec",
	}

	res, err := r.WaitProcess(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected response but got nil")
	}
	if res.Status != 0 {
		t.Errorf("expected status 0, got %d", res.Status)
	}
}

// TestUpdateContainer tests the UpdateContainer method
func TestUpdateContainer(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()
	req := &pb.UpdateContainerRequest{ContainerId: "test-container"}

	_, err := r.UpdateContainer(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestStreamOperations tests stream-related operations
func TestStreamOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("WriteStdin", func(t *testing.T) {
		req := &pb.WriteStreamRequest{
			ContainerId: "test-container",
			ExecId:      "test-exec",
			Data:        []byte("test data"),
		}
		res, err := r.WriteStdin(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Len != uint32(len(req.Data)) {
			t.Errorf("expected len %d, got %d", len(req.Data), res.Len)
		}
	})

	t.Run("ReadStdout", func(t *testing.T) {
		req := &pb.ReadStreamRequest{
			ContainerId: "test-container",
			ExecId:      "test-exec",
			Len:         100,
		}
		res, err := r.ReadStdout(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Data) == 0 {
			t.Error("expected data but got empty")
		}
	})

	t.Run("ReadStderr", func(t *testing.T) {
		req := &pb.ReadStreamRequest{
			ContainerId: "test-container",
			ExecId:      "test-exec",
			Len:         100,
		}
		res, err := r.ReadStderr(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Data) == 0 {
			t.Error("expected data but got empty")
		}
	})

	t.Run("CloseStdin", func(t *testing.T) {
		req := &pb.CloseStdinRequest{
			ContainerId: "test-container",
			ExecId:      "test-exec",
		}
		_, err := r.CloseStdin(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestNetworkOperations tests network-related operations
func TestNetworkOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("UpdateInterface", func(t *testing.T) {
		req := &pb.UpdateInterfaceRequest{}
		_, err := r.UpdateInterface(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("UpdateRoutes", func(t *testing.T) {
		req := &pb.UpdateRoutesRequest{}
		_, err := r.UpdateRoutes(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ListInterfaces", func(t *testing.T) {
		req := &pb.ListInterfacesRequest{}
		_, err := r.ListInterfaces(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ListRoutes", func(t *testing.T) {
		req := &pb.ListRoutesRequest{}
		_, err := r.ListRoutes(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("AddARPNeighbors", func(t *testing.T) {
		req := &pb.AddARPNeighborsRequest{}
		_, err := r.AddARPNeighbors(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestHealthServiceMethods tests HealthService methods
func TestHealthServiceMethods(t *testing.T) {
	mockHealth := &mockHealthClient{}
	r := &redirector{
		agentClient: &client{
			HealthService: mockHealth,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("Check", func(t *testing.T) {
		req := &pb.CheckRequest{}
		res, err := r.Check(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != pb.HealthCheckResponse_SERVING {
			t.Errorf("expected SERVING status, got %v", res.Status)
		}
	})

	t.Run("Version", func(t *testing.T) {
		req := &pb.CheckRequest{}
		res, err := r.Version(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.GrpcVersion == "" {
			t.Error("expected grpc version but got empty")
		}
		if res.AgentVersion == "" {
			t.Error("expected agent version but got empty")
		}
	})

	t.Run("Check with error", func(t *testing.T) {
		mockHealth.checkErr = errors.New("health check failed")
		req := &pb.CheckRequest{}
		_, err := r.Check(ctx, req)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("Version with error", func(t *testing.T) {
		mockHealth.versionErr = errors.New("version check failed")
		req := &pb.CheckRequest{}
		_, err := r.Version(ctx, req)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}

// TestSandboxOperations tests sandbox-related operations
func TestSandboxOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("CreateSandbox", func(t *testing.T) {
		req := &pb.CreateSandboxRequest{}
		_, err := r.CreateSandbox(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("DestroySandbox", func(t *testing.T) {
		req := &pb.DestroySandboxRequest{}
		_, err := r.DestroySandbox(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestContainerLifecycleOperations tests container lifecycle operations
func TestContainerLifecycleOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("PauseContainer", func(t *testing.T) {
		req := &pb.PauseContainerRequest{ContainerId: "test-container"}
		_, err := r.PauseContainer(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ResumeContainer", func(t *testing.T) {
		req := &pb.ResumeContainerRequest{ContainerId: "test-container"}
		_, err := r.ResumeContainer(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("StatsContainer", func(t *testing.T) {
		req := &pb.StatsContainerRequest{ContainerId: "test-container"}
		_, err := r.StatsContainer(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestMemoryOperations tests memory-related operations
func TestMemoryOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("OnlineCPUMem", func(t *testing.T) {
		req := &pb.OnlineCPUMemRequest{}
		_, err := r.OnlineCPUMem(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("MemHotplugByProbe", func(t *testing.T) {
		req := &pb.MemHotplugByProbeRequest{}
		_, err := r.MemHotplugByProbe(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("MemAgentMemcgSet", func(t *testing.T) {
		req := &pb.MemAgentMemcgConfig{}
		_, err := r.MemAgentMemcgSet(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("MemAgentCompactSet", func(t *testing.T) {
		req := &pb.MemAgentCompactConfig{}
		_, err := r.MemAgentCompactSet(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestStorageOperations tests storage-related operations
func TestStorageOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("GetVolumeStats", func(t *testing.T) {
		req := &pb.VolumeStatsRequest{}
		_, err := r.GetVolumeStats(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ResizeVolume", func(t *testing.T) {
		req := &pb.ResizeVolumeRequest{}
		_, err := r.ResizeVolume(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("AddSwap", func(t *testing.T) {
		req := &pb.AddSwapRequest{}
		_, err := r.AddSwap(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("AddSwapPath", func(t *testing.T) {
		req := &pb.AddSwapPathRequest{}
		_, err := r.AddSwapPath(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestMiscellaneousOperations tests miscellaneous operations
func TestMiscellaneousOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("TtyWinResize", func(t *testing.T) {
		req := &pb.TtyWinResizeRequest{
			ContainerId: "test-container",
			ExecId:      "test-exec",
			Row:         24,
			Column:      80,
		}
		_, err := r.TtyWinResize(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		req := &pb.CopyFileRequest{}
		_, err := r.CopyFile(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("SetGuestDateTime", func(t *testing.T) {
		req := &pb.SetGuestDateTimeRequest{}
		_, err := r.SetGuestDateTime(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ReseedRandomDev", func(t *testing.T) {
		req := &pb.ReseedRandomDevRequest{}
		_, err := r.ReseedRandomDev(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetGuestDetails", func(t *testing.T) {
		req := &pb.GuestDetailsRequest{}
		_, err := r.GetGuestDetails(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetMetrics", func(t *testing.T) {
		req := &pb.GetMetricsRequest{}
		_, err := r.GetMetrics(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetOOMEvent", func(t *testing.T) {
		req := &pb.GetOOMEventRequest{}
		_, err := r.GetOOMEvent(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("SetPolicy", func(t *testing.T) {
		req := &pb.SetPolicyRequest{}
		_, err := r.SetPolicy(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetDiagnosticData", func(t *testing.T) {
		req := &pb.GetDiagnosticDataRequest{}
		_, err := r.GetDiagnosticData(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestIPTablesOperations tests IPTables operations
func TestIPTablesOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("GetIPTables", func(t *testing.T) {
		req := &pb.GetIPTablesRequest{}
		_, err := r.GetIPTables(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("SetIPTables", func(t *testing.T) {
		req := &pb.SetIPTablesRequest{}
		_, err := r.SetIPTables(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestMountOperations tests mount-related operations
func TestMountOperations(t *testing.T) {
	mockClient := &mockAgentServiceClient{}
	r := &redirector{
		agentClient: &client{
			AgentServiceService: mockClient,
		},
		dialer: func(ctx context.Context) (net.Conn, error) {
			return &mockConn{}, nil
		},
	}
	r.once.Do(func() {})

	ctx := context.Background()

	t.Run("UpdateEphemeralMounts", func(t *testing.T) {
		req := &pb.UpdateEphemeralMountsRequest{}
		_, err := r.UpdateEphemeralMounts(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("RemoveStaleVirtiofsShareMounts", func(t *testing.T) {
		req := &pb.RemoveStaleVirtiofsShareMountsRequest{}
		_, err := r.RemoveStaleVirtiofsShareMounts(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// TestRedirectorInterfaceCompliance tests that redirector implements the Redirector interface
func TestRedirectorInterfaceCompliance(t *testing.T) {
	var _ Redirector = (*redirector)(nil)
}

// TestClientStructCompliance tests that client struct implements required interfaces
func TestClientStructCompliance(t *testing.T) {
	var _ pb.AgentServiceService = (*client)(nil)
	var _ pb.HealthService = (*client)(nil)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
