package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVacancyStatus_UnmarshalText(t *testing.T) {
	test := func(
		data []byte,
		want VacancyStatus,
		wantErr error,
	) func(*testing.T) {
		return func(t *testing.T) {
			var status VacancyStatus
			err := status.UnmarshalText(data)
			require.Exactly(t, wantErr, err)
			require.Exactly(t, want, status)
		}
	}

	tests := []struct {
		name    string
		data    []byte
		want    VacancyStatus
		wantErr error
	}{
		{
			name:    "empty",
			data:    []byte(""),
			want:    VacancyStatusNone,
			wantErr: nil,
		},
		{
			name:    "none",
			data:    []byte("none"),
			want:    VacancyStatusNone,
			wantErr: nil,
		},
		{
			name:    "draft",
			data:    []byte("draft"),
			want:    VacancyStatusDraft,
			wantErr: nil,
		},
		{
			name:    "active",
			data:    []byte("active"),
			want:    VacancyStatusActive,
			wantErr: nil,
		},
		{
			name:    "inactive",
			data:    []byte("inactive"),
			want:    VacancyStatusInactive,
			wantErr: nil,
		},
		{
			name:    "invalid",
			data:    []byte("foo"),
			want:    VacancyStatusNone,
			wantErr: ErrInvalidVacancyStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, test(tt.data, tt.want, tt.wantErr))
	}
}
