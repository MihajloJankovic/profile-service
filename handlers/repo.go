package handlers

import (
	"context"
	"fmt"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"log"
	"os"
	"time"

	// NoSQL: module containing Mongo api client
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NoSQL: ProductRepo struct encapsulating Mongo api client
type ProfileRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*ProfileRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &ProfileRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (pr *ProfileRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *ProfileRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}
func (pr *ProfileRepo) GetAll() (*[]protos.ProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()
	var profilesSlice []protos.ProfileResponse

	profileCursor, err := profileCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = profileCursor.All(ctx, &profilesSlice); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return &profilesSlice, nil
}

func (pr *ProfileRepo) GetById(emaila string) (*protos.ProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()
	var profile protos.ProfileResponse

	err := profileCollection.FindOne(ctx, bson.M{"email": emaila}).Decode(&profile)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}

	return &profile, nil
}
func (pr *ProfileRepo) Create(profile *protos.ProfileResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	profileCollection := pr.getCollection()

	result, err := profileCollection.InsertOne(ctx, &profile)
	if err != nil {
		pr.logger.Println(err)
		return err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (pr *ProfileRepo) getCollection() *mongo.Collection {
	profileDatabase := pr.cli.Database("mongoDemo")
	profileCollection := profileDatabase.Collection("profiles")
	return profileCollection
}
