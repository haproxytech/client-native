package runtime

import (
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v2/models"
)

func TestSingleRuntime_ShowACLS(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	tests := []struct {
		name           string
		fields         fields
		want           models.ACLFiles
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show acl, should return 2 ACLFiles",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: models.ACLFiles{
				&models.ACLFile{
					StorageName: "/etc/acl/blocklist.txt",
					Description: "pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59",
					ID:          "0",
				},
				&models.ACLFile{
					StorageName: "src",
					Description: "acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59",
					ID:          "1",
				},
			},
			socketResponse: map[string]string{
				"show acl\n": ` # id (file) description
				0 (/etc/acl/blocklist.txt) pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				1 () acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowACLS()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowACLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowACLS() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_GetACL(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		nameOrFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.ACLFile
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "get acl with file, should return blocklist.txt",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				nameOrFile: "/etc/acl/blocklist.txt",
			},
			want: &models.ACLFile{
				StorageName: "/etc/acl/blocklist.txt",
				Description: "pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59",
				ID:          "0",
			},
			socketResponse: map[string]string{
				"show acl\n": ` # id (file) description
				0 (/etc/acl/blocklist.txt) pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				1 () acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				`,
			},
		},
		{
			name:   "get acl with name, should return src",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				nameOrFile: "src",
			},
			want: &models.ACLFile{
				StorageName: "src",
				Description: "acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59",
				ID:          "1",
			},
			socketResponse: map[string]string{
				"show acl\n": ` # id (file) description
				0 (/etc/acl/blocklist.txt) pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				1 () acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				`,
			},
		},
		{
			name:   "get acl with not existing name, should return nil and an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				nameOrFile: "something_not_existing",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show acl\n": ` # id (file) description
				0 (/etc/acl/blocklist.txt) pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				1 () acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59
				`,
			},
		},
		{
			name:   "get acl with name, should return nil",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				nameOrFile: "",
			},
			want:           nil,
			wantErr:        true,
			socketResponse: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.GetACL(tt.args.nameOrFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetACL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetACL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowACLFileEntries(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		file string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           models.ACLFilesEntries
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "show acl with file, should 3 ACLFileEntries",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/acl/blocklist.txt",
			},
			want: models.ACLFilesEntries{
				&models.ACLFileEntry{
					ID:    "0x55c476034560",
					Value: "2.178.160.0/20",
				},
				&models.ACLFileEntry{
					ID:    "0x55c475f602f0",
					Value: "2.178.176.0/20",
				},
				&models.ACLFileEntry{
					ID:    "0x55c476035b90",
					Value: "2.178.192.0/20",
				},
			},
			socketResponse: map[string]string{
				"show acl /etc/acl/blocklist.txt\n": ` 0x55c476034560 2.178.160.0/20
				0x55c475f602f0 2.178.176.0/20
				0x55c476035b90 2.178.192.0/20
				`,
			},
		},
		{
			name:   "show acl with wrong file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "src",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show acl src\n": ` Unknown ACL identifier. Please use #<id> or <file>.
				`,
			},
		},
		{
			name:   "show acl with file, should 3 ACLFileEntries",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "",
			},
			want:           nil,
			wantErr:        true,
			socketResponse: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowACLFileEntries(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowACLFileEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowACLS() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_GetACLFileEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		aclID string
		value string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.ACLFileEntry
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "get acl entry, should find it",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.178.160.0",
			},
			want: &models.ACLFileEntry{
				Value: "2.178.160.0/20",
			},
			socketResponse: map[string]string{
				"get acl #0 2.178.160.0\n": ` type=ip, case=sensitive, match=yes, idx=tree, pattern="2.178.160.0/20"
				`,
			},
		},
		{
			name:   "get not existing acl entry, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.0.0.1/24",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"get acl #1 2.0.0.1/24\n": ` type=ip, case=sensitive, match=no
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.GetACLFileEntry(tt.args.aclID, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetACLFileEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetACLFileEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_AddACLFileEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		aclID string
		value string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "add acl entry, should work",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.0.0.1/24",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"get acl #0 2.0.0.1/24\n": ` type=ip, case=sensitive, match=no
				`,
				"add acl #0 2.0.0.1/24\n": `
				`,
			},
		},
		{
			name:   "add invalid acl entry, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.0.0.500/24",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"get acl #0 2.0.0.500/24\n": ` type=ip, case=sensitive, match=no
				`,
				"add acl #0 2.0.0.500/24\n": ` '2.0.0.500/24' is not a valid IPv4 or IPv6 address.
				`,
			},
		},
		{
			name:   "add already existing acl entry, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.178.160.0",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"get acl #0 2.178.160.0\n": ` type=ip, case=sensitive, match=yes, idx=tree, pattern="2.178.160.0/20"
				`,
			},
		},
		{
			name:   "missing aclID, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "",
				value: "2.0.0.500/24",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "missing value, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			if err := s.AddACLFileEntry(tt.args.aclID, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AddACLFileEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_DeleteACLFileEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		aclID string
		value string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "delete acl entry, should work",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "2.0.0.1/24",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"del acl #0 2.0.0.1/24\n": `
				`,
			},
		},
		{
			name:   "delete acl entry not existing, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "10.0.0.0/24",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"del acl #0 10.0.0.0/24\n": ` Key not found.
				`,
			},
		},
		{
			name:   "missing aclID, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "",
				value: "2.0.0.1/24",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "missing value, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				aclID: "0",
				value: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			if err := s.DeleteACLFileEntry(tt.args.aclID, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.DeleteACLFileEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
