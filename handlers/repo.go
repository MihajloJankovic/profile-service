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
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	var profilesSlice []protos.ProfileResponse // Create a slice to hold the results

	profileCursor, err := profileCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err // Return nil and the error
	}

	if err = profileCursor.All(ctx, &profilesSlice); err != nil {
		pr.logger.Println(err)
		return nil, err // Return nil and the error
	}

	return &profilesSlice, nil
}

//func (pr *ProfileRepo) GetById(id string) (Patients, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	patientsCollection := pr.getCollection()
//
//	var user protos.ProfileResponse
//	objID, _ := primitive.ObjectIDFromHex(id)
//	err := patientsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
//	if err != nil {
//		pr.logger.Println(err)
//		return nil, err
//	}
//	return &user, nil
//}

//	func (pr *PatientRepo) GetByName(name string) (Patients, error) {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//
//		patientsCollection := pr.getCollection()
//
//		var patients Patients
//		patientsCursor, err := patientsCollection.Find(ctx, bson.M{"name": name})
//		if err != nil {
//			pr.logger.Println(err)
//			return nil, err
//		}
//		if err = patientsCursor.All(ctx, &patients); err != nil {
//			pr.logger.Println(err)
//			return nil, err
//		}
//		return patients, nil
//	}
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

//	func (pr *PatientRepo) Update(id string, patient *Patient) error {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		patientsCollection := pr.getCollection()
//
//		objID, _ := primitive.ObjectIDFromHex(id)
//		filter := bson.M{"_id": objID}
//		update := bson.M{"$set": bson.M{
//			"name":    patient.Name,
//			"surname": patient.Surname,
//		}}
//		result, err := patientsCollection.UpdateOne(ctx, filter, update)
//		pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
//		pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)
//
//		if err != nil {
//			pr.logger.Println(err)
//			return err
//		}
//		return nil
//	}
//
//	func (pr *PatientRepo) AddPhoneNumber(id string, phoneNumber string) error {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		patientsCollection := pr.getCollection()
//
//		objID, _ := primitive.ObjectIDFromHex(id)
//		filter := bson.M{"_id": objID}
//		update := bson.M{"$push": bson.M{
//			"phoneNumbers": phoneNumber,
//		}}
//		result, err := patientsCollection.UpdateOne(ctx, filter, update)
//		pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
//		pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)
//
//		if err != nil {
//			pr.logger.Println(err)
//			return err
//		}
//		return nil
//	}
//
//	func (pr *PatientRepo) Delete(id string) error {
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//		patientsCollection := pr.getCollection()
//
//		objID, _ := primitive.ObjectIDFromHex(id)
//		filter := bson.D{{Key: "_id", Value: objID}}
//		result, err := patientsCollection.DeleteOne(ctx, filter)
//		if err != nil {
//			pr.logger.Println(err)
//			return err
//		}
//		pr.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
//		return nil
//	}
func (pr *ProfileRepo) getCollection() *mongo.Collection {
	patientDatabase := pr.cli.Database("mongoDemo")
	patientsCollection := patientDatabase.Collection("patients")
	return patientsCollection
}
