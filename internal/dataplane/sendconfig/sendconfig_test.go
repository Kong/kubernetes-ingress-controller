package sendconfig

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/kong/deck/file"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

func Test_renderConfigWithCustomEntities(t *testing.T) {
	type args struct {
		state                   *file.Content
		customEntitiesJSONBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "basic sanity test for fast-path",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: nil,
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "does not break with random bytes in the custom entities",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte("random-bytes"),
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "custom entities cannot hijack core entities",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte(`{"services":[{"host":"rogue.example.com","name":"rogue"}]}`),
			},
			want:    []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
		{
			name: "custom entities can be populated",
			args: args{
				state: &file.Content{
					FormatVersion: "1.1",
					Services: []file.FService{
						{
							Service: kong.Service{
								Name: kong.String("foo"),
								Host: kong.String("example.com"),
							},
						},
					},
				},
				customEntitiesJSONBytes: []byte(`{"my-custom-dao-name":` +
					`[{"name":"custom1","key1":"value1"},` +
					`{"name":"custom2","dumb":"test-value","boring-test-value-name":"really?"}]}`),
			},
			want: []byte(`{"_format_version":"1.1",` +
				`"my-custom-dao-name":[{"key1":"value1","name":"custom1"},` +
				`{"boring-test-value-name":"really?","dumb":"test-value","name":"custom2"}]` +
				`,"services":[{"host":"example.com","name":"foo"}]}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderConfigWithCustomEntities(logrus.New(), tt.args.state, tt.args.customEntitiesJSONBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderConfigWithCustomEntities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderConfigWithCustomEntities() = %v, want %v",
					string(got), string(tt.want))
			}
		})
	}
}

func Test_updateReportingUtilities(t *testing.T) {
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
}

func Test_pushFailureReason(t *testing.T) {
	apiConflictErr := kong.NewAPIError(http.StatusConflict, "conflict in configuration")

	testCases := []struct {
		name           string
		err            error
		expectedReason string
	}{
		{
			name:           "generic_error",
			err:            errors.New("some generic error"),
			expectedReason: metrics.FailureReasonOther,
		},
		{
			name:           "api_conflict_error",
			err:            apiConflictErr,
			expectedReason: metrics.FailureReasonConflict,
		},
		{
			name:           "api_conflict_error_wrapped",
			err:            fmt.Errorf("wrapped conflict api err: %w", apiConflictErr),
			expectedReason: metrics.FailureReasonConflict,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			reason := pushFailureReason(tc.err)
			require.Equal(t, tc.expectedReason, reason)
		})
	}
}
