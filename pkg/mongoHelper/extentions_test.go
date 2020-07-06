package mongoHelper

import (
	"encoding/hex"
	"github.com/alperhankendi/go-rest-api/internal/config"
	"github.com/alperhankendi/go-rest-api/pkg/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

func TestMongoObject_CastToId(t *testing.T) {


		tests := []struct {
			name string
			args string
			want bson.M
			wantErr error
		}{
			{name: "WithValid", args: "5d2399ef96fb765873a24bae", want: bson.M{"_id":primitive.ObjectID{93,35,153,239,150,251,118,88,115,162,75,174}},wantErr: nil},
			{name: "WithInvalidLen", args: "ABC", want: bson.M{"_id":primitive.ObjectID{93,35,153,239,150,251,118,88,115,162,75,174}},wantErr: hex.ErrLength},//primitive.ErrInvalidHex},
			{name: "WithInvalidLen", args: "ABCDEF", want: bson.M{"_id":primitive.ObjectID{93,35,153,239,150,251,118,88,115,162,75,174}},wantErr: primitive.ErrInvalidHex},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got,err := CastToId(tt.args);
				t.Logf("Err:%v",err)
				t.Logf("Got:%v",got)

				if (err != nil) &&  err != tt.wantErr {
					t.Errorf("CastToId() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err==nil && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CastToId() = %v, want %v", got, tt.want)
				}
			})
		}

}

func TestConnectDb(t *testing.T) {

	tests := []struct {
		name    string
		args    *config.MongoSettings
		wantErr bool
	}{
		{name: "WithValidConfiguration_ShouldBeConnected", args: &config.MongoSettings{
			Uri: "mongodb://root:example@127.0.0.1:27017/admin",
		} ,  wantErr:false},
		{name: "WithInValidDatabase_ShouldBeFailed", args: &config.MongoSettings{
			Uri: "mongodb://root:example@127.0.0.1:27017/admin1",
		} ,  wantErr:true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := ConnectDb(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectDb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSetFindOptions(t *testing.T) {

	tests := []struct {
		name            string
		args            *pagination.Pages
		wantFindOptions *options.FindOptions
	}{
		{name: "WithValidLimitPageOptions", args: pagination.New(1,5,50) , wantFindOptions: options.Find().SetLimit(5),},
		{name: "WithValidSkipAndLimitPageOptions", args: pagination.New(2,5,50) , wantFindOptions: options.Find().SetSkip(5).SetLimit(5),},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFindOptions := SetFindOptions(tt.args); !reflect.DeepEqual(gotFindOptions, tt.wantFindOptions) {
				t.Errorf("SetFindOptions() = %v, want %v", gotFindOptions, tt.wantFindOptions)
			}
		})
	}
}

func TestBuildQuery(t *testing.T) {

	_idTest,_ :=primitive.ObjectIDFromHex("5d2399ef96fb765873a24bae")
	bsonTest :=bson.M{}
	bsonTest["_id"] = _idTest
	tests := []struct {
		name      string
		args      map[string]string
		wantQuery bson.M
	}{
		{name: "WithEmptyParameter", args:nil , wantQuery:bson.M{} },
		{name:"WithOneParameter", args: map[string]string{"Key":"Test"}, wantQuery: bson.M{"Key":"Test"}},
		{name:"WithOneParameterAndValueIsEmpty", args: map[string]string{"Key":""}, wantQuery: bson.M{}},
		{name:"WithIdParameter", args: map[string]string{"_id":"5d2399ef96fb765873a24bae"}, wantQuery: bsonTest },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotQuery := BuildQuery(tt.args); !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Errorf("BuildQuery() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}