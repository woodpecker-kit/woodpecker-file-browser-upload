package file_browser_upload_test

import (
	"github.com/sinlov-go/unittest-kit/env_kit"
	"github.com/stretchr/testify/assert"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/file_browser_upload"
	"testing"
)

func TestHttpHttpBestLinkIgnoreRetry(t *testing.T) {

	testUrls := env_kit.FetchOsEnvStringSlice(keyLinkSpeedTestUrls)
	if len(testUrls) == 0 {
		t.Logf("skip test by not set env: %s", keyLinkSpeedTestUrls)
		return
	}

	// mock HttpBestLinkIgnoreRetry
	type args struct {
		testUrls      []string
		timeoutSecond uint
		retries       uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "sample",
			args: args{
				testUrls:      testUrls,
				timeoutSecond: 5,
				retries:       3,
			},
		},
		{
			name: "with-error-link",
			args: args{
				testUrls:      append(testUrls, "http://192.168.30.20"),
				timeoutSecond: 5,
				retries:       3,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			linkSpeed := file_browser_upload.NewLinkSpeed(tc.args.timeoutSecond, tc.args.retries)

			// do HttpBestLinkIgnoreRetry
			gotResult, gotErr := linkSpeed.BestLinkIgnoreRetry(tc.args.testUrls)

			// verify HttpBestLinkIgnoreRetry
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
			t.Logf("gotResult: %v", gotResult)
		})
	}
}

func TestHttpLinkSpeed(t *testing.T) {

	testUrls := env_kit.FetchOsEnvStringSlice(keyLinkSpeedTestUrls)
	if len(testUrls) == 0 {
		t.Logf("skip test by not set env: %s", keyLinkSpeedTestUrls)
		return
	}

	// mock HttpLinkSpeed
	type args struct {
		//
		testUrls      []string
		timeoutSecond uint
		retries       uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "sample",
			args: args{
				testUrls:      testUrls,
				timeoutSecond: 5,
				retries:       3,
			},
		},
		{
			name: "error-link",
			args: args{
				testUrls: []string{
					"http://192.168.30.200",
					"http://127.0.0.1:50001",
					"http://127.0.0.1:60001",
				},
				timeoutSecond: 5,
				retries:       3,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// do HttpLinkSpeed
			linkSpeed := file_browser_upload.NewLinkSpeed(tc.args.timeoutSecond, tc.args.retries)
			gotResult, gotErr := linkSpeed.DoUrls(tc.args.testUrls)

			// verify HttpLinkSpeed
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
			//wd_log.VerboseJsonf(gotResult, "gotResult")
			t.Logf("gotResult: %v", gotResult)
		})
	}
}
