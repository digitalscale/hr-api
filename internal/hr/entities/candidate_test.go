package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGender_UnmarshalText(t *testing.T) {
	test := func(data []byte, want Gender, wantErr error) func(*testing.T) {
		return func(t *testing.T) {
			var gender Gender
			err := gender.UnmarshalText(data)
			require.Exactly(t, wantErr, err)
			require.Exactly(t, want, gender)
		}
	}

	tests := []struct {
		name    string
		data    []byte
		want    Gender
		wantErr error
	}{
		{
			name:    "empty",
			data:    []byte(""),
			want:    GenderNone,
			wantErr: nil,
		},
		{
			name:    "none",
			data:    []byte("none"),
			want:    GenderNone,
			wantErr: nil,
		},
		{
			name:    "male",
			data:    []byte("male"),
			want:    GenderMale,
			wantErr: nil,
		},
		{
			name:    "female",
			data:    []byte("female"),
			want:    GenderFemale,
			wantErr: nil,
		},
		{
			name:    "invalid",
			data:    []byte("foo"),
			want:    GenderNone,
			wantErr: ErrInvalidGender,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, test(tt.data, tt.want, tt.wantErr))
	}
}

func TestEducationLevel_UnmarshalText(t *testing.T) {
	test := func(data []byte, want EducationLevel, wantErr error) func(*testing.T) {
		return func(t *testing.T) {
			var gender EducationLevel
			err := gender.UnmarshalText(data)
			require.Exactly(t, wantErr, err)
			require.Exactly(t, want, gender)
		}
	}

	tests := []struct {
		name    string
		data    []byte
		want    EducationLevel
		wantErr error
	}{
		{
			name:    "empty",
			data:    []byte(""),
			want:    EducationLevelNone,
			wantErr: nil,
		},
		{
			name:    "none",
			data:    []byte("none"),
			want:    EducationLevelNone,
			wantErr: nil,
		},
		{
			name:    "secondary",
			data:    []byte("secondary"),
			want:    EducationLevelSecondary,
			wantErr: nil,
		},
		{
			name:    "specialSecondary",
			data:    []byte("specialSecondary"),
			want:    EducationLevelSpecialSecondary,
			wantErr: nil,
		},
		{
			name:    "unfinishedHigher",
			data:    []byte("unfinishedHigher"),
			want:    EducationLevelUnfinishedHigher,
			wantErr: nil,
		},
		{
			name:    "higher",
			data:    []byte("higher"),
			want:    EducationLevelHigher,
			wantErr: nil,
		},
		{
			name:    "bachelor",
			data:    []byte("bachelor"),
			want:    EducationLevelBachelor,
			wantErr: nil,
		},
		{
			name:    "master",
			data:    []byte("master"),
			want:    EducationLevelMaster,
			wantErr: nil,
		},
		{
			name:    "candidate",
			data:    []byte("candidate"),
			want:    EducationLevelCandidate,
			wantErr: nil,
		},
		{
			name:    "doctor",
			data:    []byte("doctor"),
			want:    EducationLevelDoctor,
			wantErr: nil,
		},
		{
			name:    "invalid",
			data:    []byte("foo"),
			want:    EducationLevelNone,
			wantErr: ErrInvalidEducationLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, test(tt.data, tt.want, tt.wantErr))
	}
}
