package signal_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/hudymi/mockice/pkg/signal"

	. "github.com/onsi/gomega"
)

func TestContext(t *testing.T) {
	// Given
	g := NewGomegaWithT(t)
	ctx := signal.Context()

	// When
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	// Then
	g.Eventually(ctx.Done(), time.Second).Should(BeClosed())
}
