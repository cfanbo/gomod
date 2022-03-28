package core

import "testing"

func Test_convertGoModUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"github.com/nsf", args{url: "github.com/nsf"}, "github.com/nsf", true},
		{"github.com/nsf/", args{url: "github.com/nsf"}, "github.com/nsf", true},
		{"https://github.com/nsf", args{url: "github.com/nsf"}, "github.com/nsf", true},
		{"https://github.com/", args{url: "github.com/"}, "github.com/", true},
		{"github.com", args{url: "github.com"}, "github.com", false},
		{"github.com/nsf/termbox-go", args{url: "github.com/nsf/termbox-go"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"https://github.com/nsf/termbox-go/", args{url: "github.com/nsf/termbox-go/"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"github.com/nsf/termbox-go/a/b", args{url: "github.com/nsf/termbox-go/a/b"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"github.com/nsf/termbox-go/a/b/", args{url: "github.com/nsf/termbox-go/a/b/"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"https://github.com/nsf/termbox-go/a/b/c/d", args{url: "github.com/nsf/termbox-go/a/b/c/d"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"github.com/nsf/termbox-go/a", args{url: "github.com/nsf/termbox-go/a"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"https://github.com/nsf/termbox-go/a/", args{url: "github.com/nsf/termbox-go/a/"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
		{"github.com/nsf/termbox-go/a/", args{url: "github.com/nsf/termbox-go/a/"}, "https://raw.githubusercontent.com/nsf/termbox-go/master/go.mod", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGoModURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertGoModUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertGoModUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
