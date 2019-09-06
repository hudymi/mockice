package service_test

import (
	"context"
	"github.com/michal-hudy/mockice/pkg/endpoint"
	"github.com/michal-hudy/mockice/pkg/service"
	"sync"
	"testing"

	. "github.com/onsi/gomega"
)

func TestService_Start(t *testing.T) {
	// Given
	g := NewGomegaWithT(t)
	svc := service.New("localhost:8080")
	ctx, cancel := context.WithCancel(context.TODO())

	// When
	var err error
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		err = svc.Start(ctx)
		wait.Done()
	}()
	cancel()
	wait.Wait()

	// Then
	g.Expect(err).To(Succeed())
}

func TestService_Register(t *testing.T) {
	for testName, testCase := range map[string]struct {
		shouldFail       bool
		endpointsConfigs []endpoint.Config
	}{
		"Default": {
			shouldFail:       false,
			endpointsConfigs: endpoint.DefaultConfig(),
		},
		"Multiple": {
			shouldFail: false,
			endpointsConfigs: []endpoint.Config{
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
			},
		},
		"Noname": {
			shouldFail:       true,
			endpointsConfigs: []endpoint.Config{{}},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			// Given
			g := NewGomegaWithT(t)
			svc := service.New("localhost:8080")

			for _, config := range testCase.endpointsConfigs {
				end := endpoint.New(config)

				// When
				err := svc.Register(end)

				// Then
				if testCase.shouldFail {
					g.Expect(err).ToNot(Succeed())
				} else {
					g.Expect(err).To(Succeed())
				}
			}
		})
	}
}
