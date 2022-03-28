package core

import (
	"testing"
)

func TestRepo_Get(t *testing.T) {

	InitCache()
	InitPkgMap()

	type fields struct {
		Star     int
		Fork     int
		Watch    int
		ImportBy int
		RepoUrl  string
		Mod      string
		Shared   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{"RepoGet-jenkins-demo", fields{Mod: "github.com/cfanbo/jenkins-demo"}, fields{Star: 0, Fork: 0}},
		//{"RepoGet-delayqueue", fields{Mod: "github.com/cfanbo/delayqueue"}, fields{Star:12, Fork:4}},
		//{"RepoGet-golang", fields{Mod: "github.com/golang/go"}, fields{Star:96980, Fork:14478}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				Star:     tt.fields.Star,
				Fork:     tt.fields.Fork,
				Watch:    tt.fields.Watch,
				ImportBy: tt.fields.ImportBy,
				RepoUrl:  tt.fields.RepoUrl,
				Mod:      tt.fields.Mod,
				Shared:   tt.fields.Shared,
			}
			if err := r.Do(); err != nil {
				t.Errorf("%#v", err)
			}

			if r.Star != tt.want.Star {
				t.Errorf("TestRepo_Get().Star got = %v, want %v", r.Star, tt.want.Star)
			}

			if r.Fork != tt.want.Fork {
				t.Errorf("TestRepo_Get().Fork got = %v, want %v", r.Fork, tt.want.Fork)
			}
		})
	}
}
