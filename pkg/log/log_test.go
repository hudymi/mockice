package log_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/michal-hudy/mockice/pkg/log"
)

func TestSetup(t *testing.T) {
	t.Run("Verbose", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)

		// When
		log.Setup(true)

		// Then
		g.Expect(logrus.GetLevel()).To(Equal(logrus.InfoLevel))
	})

	t.Run("Warn", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)

		// When
		log.Setup(false)

		// Then
		g.Expect(logrus.GetLevel()).To(Equal(logrus.WarnLevel))
	})
}
