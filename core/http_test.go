package core

import "testing"

func Test_parseStar(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"0", args{`data-pjax-replace="true" title="0" data-view-component="true" class="Counter js-social-count">1</span>`}, 0, false},
		{"<10", args{`data-pjax-replace="true" title="1" data-view-component="true" class="Counter js-social-count">1</span>`}, 1, false},
		{"<1K", args{`data-pjax-replace="true" title="123" data-view-component="true" class="Counter js-social-count">1</span>`}, 123, false},
		{">=1K", args{`data-pjax-replace="true" title="1,234" data-view-component="true" class="Counter js-social-count">1</span>`}, 1234, false},
		{"google/go", args{`<span id="repo-stars-counter-star" aria-label="96979 users starred this repository" data-singular-suffix="user starred this repository" data-plural-suffix="users starred this repository" data-pjax-replace="true" title="96,979" data-view-component="true" class="Counter js-social-count">97k</span>`}, 96979, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStar(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseStar() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFork(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"0", args{`<span id="repo-network-counter" data-pjax-replace="true" title="0" data-view-component="true" class="Counter">0</span>`}, 0, false},
		{"<10", args{`<span id="repo-network-counter" data-pjax-replace="true" title="1" data-view-component="true" class="Counter">0</span>`}, 1, false},
		{"<1K", args{`<span id="repo-network-counter" data-pjax-replace="true" title="123" data-view-component="true" class="Counter">0</span>`}, 123, false},
		{">=1K", args{`<span id="repo-network-counter" data-pjax-replace="true" title="1,234" data-view-component="true" class="Counter">0</span>`}, 1234, false},
		{"golang/go", args{`<span id="repo-network-counter" data-pjax-replace="true" title="14,478" data-view-component="true" class="Counter">14.5k</span>`}, 14478, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFork(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseFork() got = %v, want %v", got, tt.want)
			}
		})
	}
}
