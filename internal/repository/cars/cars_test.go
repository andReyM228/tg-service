package cars

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"testing"
	"tg_service/internal/domain"
)

func TestRepository_Get(t *testing.T) {
	type fields struct {
		log    *logrus.Logger
		client *http.Client
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Car
		wantErr bool
	}{{
		name: "success",
		fields: fields{
			log:    logrus.New(),
			client: http.DefaultClient,
		},
		args: args{
			id: 1,
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				log:    tt.fields.log,
				client: tt.fields.client,
			}
			got, err := r.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
