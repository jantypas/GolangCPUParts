package RemoteLogging

import "testing"

func TestLogEvent(t *testing.T) {
	type args struct {
		level  string
		source string
		msg    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogEvent(tt.args.level, tt.args.source, tt.args.msg)
		})
	}
}

func TestLogInit(t *testing.T) {
	type args struct {
		appname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LogInit(tt.args.appname); (err != nil) != tt.wantErr {
				t.Errorf("LogInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetLogginActive(t *testing.T) {
	type args struct {
		state bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogginActive(tt.args.state)
		})
	}
}
