package database

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tUUID = reflect.TypeOf(uuid.UUID{})

func uuidEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	u := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(u[:], bson.TypeBinaryUUID)
}

func uuidDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var err error

	switch vrType := vr.Type(); vrType {
	case bson.TypeBinary:
		data, _, err = vr.ReadBinary()
		if err != nil {
			return err
		}
	case bson.TypeNull:
		if err = vr.ReadNull(); err != nil {
			return err
		}
		val.Set(reflect.ValueOf(uuid.UUID{}))
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if len(data) != 16 {
		return fmt.Errorf("invalid UUID length: %d", len(data))
	}

	var u uuid.UUID
	copy(u[:], data)
	val.Set(reflect.ValueOf(u))
	return nil
}

type MongoConfig struct {
	URI         string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
	MaxIdleTime time.Duration
}

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoConnection(cfg MongoConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Register UUID codec to handle MongoDB Binary UUID (subtype 4)
	reg := bson.NewRegistry()
	reg.RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue))
	reg.RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue))

	clientOptions := options.Client().ApplyURI(cfg.URI).SetRegistry(reg)

	if cfg.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	}
	if cfg.MaxIdleTime > 0 {
		clientOptions.SetMaxConnIdleTime(cfg.MaxIdleTime)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}
