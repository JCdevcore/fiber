package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ExponentialBackoff_Retry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expErr     error
		expBackoff *ExponentialBackoff
		f          func() error
		name       string
	}{
		{
			name:       "With default values - successful",
			expBackoff: NewExponentialBackoff(),
			f: func() error {
				return nil
			},
		},
		{
			name: "Successful function",
			expBackoff: &ExponentialBackoff{
				InitialInterval: 1 * time.Millisecond,
				MaxBackoffTime:  100 * time.Millisecond,
				Multiplier:      2.0,
				MaxRetryCount:   5,
			},
			f: func() error {
				return nil
			},
		},
		{
			name: "Unsuccessful function",
			expBackoff: &ExponentialBackoff{
				InitialInterval: 2 * time.Millisecond,
				MaxBackoffTime:  100 * time.Millisecond,
				Multiplier:      2.0,
				MaxRetryCount:   5,
			},
			f: func() error {
				return errors.New("failed function")
			},
			expErr: errors.New("failed function"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.expBackoff.Retry(tt.f)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func Test_ExponentialBackoff_Next(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		expBackoff           *ExponentialBackoff
		expNextTimeIntervals []time.Duration
	}{
		{
			name:       "With default values",
			expBackoff: NewExponentialBackoff(),
			expNextTimeIntervals: []time.Duration{
				1 * time.Second,
				2 * time.Second,
				4 * time.Second,
				8 * time.Second,
				16 * time.Second,
				32 * time.Second,
				32 * time.Second,
				32 * time.Second,
				32 * time.Second,
				32 * time.Second,
			},
		},
		{
			name: "Custom values",
			expBackoff: &ExponentialBackoff{
				InitialInterval: 2.0 * time.Second,
				MaxBackoffTime:  64 * time.Second,
				Multiplier:      3.0,
				MaxRetryCount:   8,
				currentInterval: 2.0 * time.Second,
			},
			expNextTimeIntervals: []time.Duration{
				2 * time.Second,
				6 * time.Second,
				18 * time.Second,
				54 * time.Second,
				64 * time.Second,
				64 * time.Second,
				64 * time.Second,
				64 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for i := range tt.expBackoff.MaxRetryCount {
				next := tt.expBackoff.next()
				if next < tt.expNextTimeIntervals[i] || next > tt.expNextTimeIntervals[i]+1*time.Second {
					t.Errorf("wrong next time:\n"+
						"actual:%v\n"+
						"expected range:%v-%v\n",
						next, tt.expNextTimeIntervals[i], tt.expNextTimeIntervals[i]+1*time.Second)
				}
			}
		})
	}
}
