package rotate

import (
	"testing"
)

func TestRotateWriter_Rotate(t *testing.T) {

	conf := &Conf{
		Directory: "/tmp/",
		Prefix:    "log",
		Size:      10,
		Count:     5,
	}

	w, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		writer *RotateWriter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{writer: w},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fields.writer
			if err := w.Rotate(); (err != nil) != tt.wantErr {
				t.Errorf("RotateWriter.Rotate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRotateWriter_Write(t *testing.T) {

	data := []byte("hola mundo, hola mundo\n")
	conf := &Conf{
		Directory: "/tmp/",
		Prefix:    "log",
		Size:      32,
		Count:     5,
	}

	w, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		writer *RotateWriter
	}

	type args struct {
		output []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{writer: w},
			args: args{
				output: data,
			},
			wantErr: false,
			want:    len(data),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := w.Write(tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("RotateWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RotateWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
