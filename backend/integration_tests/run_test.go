package integration_tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestIntegration(t *testing.T) {
	suite.Run(t, new(SongSuite))
}
