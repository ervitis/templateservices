package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type (
	serverTestSuite struct {
		suite.Suite
	}

	mockLogger struct {
		mock.Mock
	}
	mockLogrus struct{}
)

func (m *mockLogrus) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (m *mockLogrus) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (m *mockLogger) L() *logrus.Entry {
	return &logrus.Entry{}
}

func (s *serverTestSuite) TestServerInstance() {
	r := http.NewServeMux()
	assert.NotNil(s.T(), CreateServer(WithPort("8080"), WithRouter(r)))
}

func (s *serverTestSuite) TestGetStringConnection() {
	r := http.NewServeMux()
	l := new(mockLogger)

	l.On("L").Return(&mockLogrus{})

	srv := CreateServer(WithRouter(r), WithLogger(l))

	assert.Equal(s.T(), "127.0.0.1:8080", srv.getStringConnection())
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}
