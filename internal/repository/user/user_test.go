package user

import (
	"github.com/andReyM228/lib/log"
	"github.com/andReyM228/lib/rabbit"
	"net/http"
	"testing"
)

func TestRepository_Get(t *testing.T) {
	type fields struct {
		log    log.Logger
		client *http.Client
		rabbit rabbit.Rabbit
	}
	type args struct {
		id int64
	}
	rb, err := rabbit.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				log:    log.Init(),
				client: nil,
				rabbit: rb,
			},
			args: args{
				id: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				log:    tt.fields.log,
				client: tt.fields.client,
				rabbit: tt.fields.rabbit,
			}
			got, err := r.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
