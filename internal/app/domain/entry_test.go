package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntry_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantErr     bool
		expectedErr string
		expectedRes *Entry
	}{
		{
			name:    "TimePass",
			data:    []byte(`{"time":"2020-07-26T16:05:27.429286952+03:00"}`),
			wantErr: false,
			expectedRes: &Entry{
				Time:   time.Unix(0, 1595768727429286952),
				Params: []byte(`{"time":"2020-07-26T16:05:27.429286952+03:00"}`),
			},
		},
		{
			name:        "TimeError",
			data:        []byte(`{"time":1595768727429286952}`),
			wantErr:     true,
			expectedErr: "Value looks like Number/Boolean/None, but can't find its end: ',' or '}' symbol",
		},
		{
			name:    "NamespacePass",
			data:    []byte(`{"namespace":"PROD"}`),
			wantErr: false,
			expectedRes: &Entry{
				Namespace: "prod",
				Params:    []byte(`{"namespace":"PROD"}`),
			},
		},
		{
			name:    "SourcePass",
			data:    []byte(`{"source":"app_1"}`),
			wantErr: false,
			expectedRes: &Entry{
				Source: "app_1",
				Params: []byte(`{"source":"app_1"}`),
			},
		},
		{
			name:    "HostPass",
			data:    []byte(`{"host":"127.0.0.1:4222"}`),
			wantErr: false,
			expectedRes: &Entry{
				Host:   "127.0.0.1:4222",
				Params: []byte(`{"host":"127.0.0.1:4222"}`),
			},
		},
		{
			name:    "LevelPass",
			data:    []byte(`{"level":"info"}`),
			wantErr: false,
			expectedRes: &Entry{
				Level:  "info",
				Params: []byte(`{"level":"info"}`),
			},
		},
		{
			name:    "TraceIDPass",
			data:    []byte(`{"trace_id":"qwertyu123456"}`),
			wantErr: false,
			expectedRes: &Entry{
				TraceID: "qwertyu123456",
				Params:  []byte(`{"trace_id":"qwertyu123456"}`),
			},
		},
		{
			name:    "MessagePass",
			data:    []byte(`{"message":"some Message"}`),
			wantErr: false,
			expectedRes: &Entry{
				Message: "some message",
				Params:  []byte(`{"message":"some Message"}`),
			},
		},
		{
			name:    "BuildCommitPass",
			data:    []byte(`{"build_commit":"130362a6fd10cf2f939dd0cfc0ab222cee6a99ec"}`),
			wantErr: false,
			expectedRes: &Entry{
				BuildCommit: "130362a6fd10cf2f939dd0cfc0ab222cee6a99ec",
				Params:      []byte(`{"build_commit":"130362a6fd10cf2f939dd0cfc0ab222cee6a99ec"}`),
			},
		},
		{
			name:    "ConfigHashPass",
			data:    []byte(`{"config_hash":"130362a6fd10cf2f939dd0cfc0ab222cee6a99ec"}`),
			wantErr: false,
			expectedRes: &Entry{
				ConfigHash: "130362a6fd10cf2f939dd0cfc0ab222cee6a99ec",
				Params:     []byte(`{"config_hash":"130362a6fd10cf2f939dd0cfc0ab222cee6a99ec"}`),
			},
		},
		{
			name:    "StringPass",
			data:    []byte(`{"key":["a","b","c","d","e","f"]}`),
			wantErr: false,
			expectedRes: &Entry{
				StringKey: []string{"key", "key", "key", "key", "key", "key"},
				StringVal: []string{"a", "b", "c", "d", "e", "f"},
				Params:    []byte(`{"key":["a","b","c","d","e","f"]}`),
			},
		},
		{
			name:    "FloatPass",
			data:    []byte(`{"key":[1, 2, 3, 4, 5, 6]}`),
			wantErr: false,
			expectedRes: &Entry{
				FloatKey: []string{"key", "key", "key", "key", "key", "key"},
				FloatVal: []float64{1, 2, 3, 4, 5, 6},
				Params:   []byte(`{"key":[1, 2, 3, 4, 5, 6]}`),
			},
		},
		{
			name:    "BigObjectPass",
			data:    []byte(`{"key3":{"key4":{"key5":[11,12,13,"WWW",{"someKey": "someValue"}]}}}`),
			wantErr: false,
			expectedRes: &Entry{
				StringKey: []string{"key5", "somekey"},
				StringVal: []string{"www", "somevalue"},
				FloatKey:  []string{"key5", "key5", "key5"},
				FloatVal:  []float64{11, 12, 13},
				Params:    []byte(`{"key3":{"key4":{"key5":[11,12,13,"WWW",{"someKey": "someValue"}]}}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &Entry{}

			err := entry.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}

			if !tt.wantErr {
				assert.Equal(t, tt.expectedRes, entry)
			}
		})
	}
}

func BenchmarkEntry_UnmarshalJSON(b *testing.B) {
	data := []byte(`{"time":1593709730877594291,"namespace":"prod","source":"app_1","host":"127.100.0.1:50000","level":"debug","trace_id":"1c7fuhpo0ln2dcq","message":"read failed: some message 4","build_commit":"db957a22b3c1d6e508c0828917a5e14c572fb007","config_hash":"130362a6fd10cf2f939dd0cfc0ab222cee6a99ec","key1":[1,2,3,4,5,6],"key2":["a","b","c","d","e","f"],"boolKey":true,"key3":{"key4":{"key5":[11,11,11,"WWW",{"someKey":"someValue"}]}}}`)

	entry := &Entry{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := entry.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
	}
}
