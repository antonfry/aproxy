package aproxy

import (
	"aproxy/internal/roundrobin"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Listener interface {
	ReadFrom(b []byte) (int, net.Addr, error)
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	SetLevel(log.Level)
}

type Sender interface {
}

// server is the main server struct
type server struct {
	logger      Logger
	BackendPool *roundrobin.Pool
	ListenConn  Listener
	config      *Config
	shutdown    chan struct{}
}

// New create a new server struct
func New(config *Config, pool *roundrobin.Pool) *server {
	return &server{
		logger:      log.New(),
		BackendPool: pool,
		config:      config,
		shutdown:    make(chan struct{}),
	}
}

// Start starts the server
func (s *server) Start() error {
	var err error
	if err = s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Starting HealthCheck. ", s.BackendPool.TargetGroup.HealthCheck)
	go s.healthCheckWorker()
	s.logger.Info("Binding socket: ", s.config.Host, s.config.Port)
	s.ListenConn, err = s.udpServer()
	if err != nil {
		return err
	}
	for _, b := range s.BackendPool.TargetGroup.Backends {
		s.logger.Info("Setup backend connection: ", b.Host)
		if err := b.GetConn(); err != nil {
			return err
		}
	}
	for i := 0; i < s.config.Workers; i++ {
		s.logger.Info("Starting worker: ", i)
		go s.worker()
	}
	s.logger.Info("Aproxy started listen with config: ", s.config)
	s.logger.Info("Target group: ", s.BackendPool.TargetGroup.Backends)
	s.logger.Info("HealthCheck: ", s.BackendPool.TargetGroup.HealthCheck)
	go s.gracefullShutdown()
	<-s.shutdown
	return s.Stop()
}

// Stop stops server
func (s *server) Stop() error {
	for _, h := range s.BackendPool.TargetGroup.Backends {
		if h != nil {
			if err := h.Close(); err != nil {
				s.logger.Error("Unable to close connection: ", err.Error())
			}
		}
	}
	s.logger.Info("Aproxy stopped")
	return nil
}

// configureLogger configures logger
func (s *server) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

// gracefullShutdown listnes os signals and send stop signal to server
func (s *server) gracefullShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		func() {
			sig := <-sigs
			s.logger.Info("Got os signal: ", sig.String())
			s.shutdown <- struct{}{}
		}()
	}
}

// udpServer
func (s *server) udpServer() (*net.UDPConn, error) {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	host, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		s.logger.Error("Unable to resolve: ", addr, err.Error())
		return nil, err
	}
	s.logger.Info("Starting UDP server on: ", host.String())
	conn, err := net.ListenUDP("udp", host)
	if err != nil {
		s.logger.Error("Unable to start UDP server: ", addr, err.Error())
		return nil, err
	}
	return conn, err
}

// healthCheckWorker make healchecks
func (s *server) healthCheckWorker() {
	healchCheckTicker := time.NewTicker(s.BackendPool.TargetGroup.HealthCheck.Interval)
L:
	for {
		select {
		case <-healchCheckTicker.C:
			s.logger.Debug("HealthCheck started")
			if err := s.BackendPool.TargetGroup.Check(); err != nil {
				s.logger.Error("Healcheck failed: ", err.Error())
			}
		case <-s.shutdown:
			healchCheckTicker.Stop()
			break L
		}
	}
}

// worker read and proxy packets to the backends servers
func (s *server) worker() {
	b := make([]byte, s.config.Buffer)
L:
	for {
		select {
		case <-s.shutdown:
			break L
		default:
			n, _, err := s.ListenConn.ReadFrom(b)
			if err != nil {
				s.logger.Error("Could not read a packet: ", err.Error())
				continue
			}
			targetConn, err := s.BackendPool.Next()
			if err != nil {
				s.logger.Error("No healthy target connections: ", err.Error())
				continue
			}
			s.logger.Debug("Sending packet to: ", targetConn.Conn.RemoteAddr().String())
			_, err = targetConn.Conn.Write(b[0:n])
			if err != nil {
				s.logger.Error("Fail to write into: ", targetConn.Conn.RemoteAddr().String(), err.Error())
				targetConn.SetDead()
				continue
			}
		}
	}
}
