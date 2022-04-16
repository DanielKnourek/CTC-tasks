package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URI = "mongodb://db_mongo:27017"
const MONGO_DB_NAME = "test"
const MONGO_COLLECTION_NAME = "products"

type Product struct {
	Id     	primitive.ObjectID 	`bson:"_id,omitempty" json:"id,omitempty"`
	Name    string  			`bson:"name,omitempty"`
	Price   float64 			`bson:"price"`
	Ammount int     			`bson:"ammount"`
}

func errorPage(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"error\": \"%v\"}\n", err)
}

var collection *mongo.Collection
var ctx context.Context

func initializeConnection() (error) {
	ctx, ctx_cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctx_cancel()
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// fmt.Println("Connected to MongoDB!")
	collection = client.Database(MONGO_DB_NAME).Collection(MONGO_COLLECTION_NAME)
	return nil
}

func Put(w http.ResponseWriter, r *http.Request) {
	err := initializeConnection()
	if err != nil {
		errorPage(w, err)
		return
	}
	
	// create product from request
	var Product_req Product
	err = json.NewDecoder(r.Body).Decode(&Product_req)
	if err != nil {
		errorPage(w, err)
		return
	}
	if Product_req.Name == "" {
		errorPage(w, fmt.Errorf("missing required fields: [name]"))
		return
	}	

	// insert the product
	insertResult, err := collection.InsertOne(ctx, Product_req)
	if err != nil {
		errorPage(w, err)
		return
	}

	res, err := json.Marshal(insertResult)
	if err != nil {
		errorPage(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

func List(w http.ResponseWriter, r *http.Request) {
	err := initializeConnection()
	if err != nil {
		errorPage(w, err)
		return
	}
	
	filter := bson.D{{}}
	results, err := collection.Find(ctx, filter)
	if err != nil {
		errorPage(w, err)
		return
	}
	defer results.Close(ctx)

	var result []Product
	if err := results.All(ctx, &result); err != nil {
		errorPage(w, err)
		return
	}
	res, err := json.Marshal(result)
	if err != nil {
		errorPage(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

func Get(w http.ResponseWriter, r *http.Request) {
	err := initializeConnection()
	if err != nil {
		errorPage(w, err)
		return
	}
	vars := mux.Vars(r)

	var result Product

	hex, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return
	}
	filter := bson.D{{Key: "_id", Value: hex}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		errorPage(w, err)
		return
	}
	
	res, err := json.Marshal(result)
	if err != nil {
		errorPage(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := initializeConnection()
	if err != nil {
		errorPage(w, err)
		return
	}
	vars := mux.Vars(r)

	var result Product

	hex, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return
	}
	filter := bson.D{{Key: "_id", Value: hex}}
	err = collection.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		errorPage(w, err)
		return
	}
	
	res, err := json.Marshal(result)
	if err != nil {
		errorPage(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

func Update(w http.ResponseWriter, r *http.Request) {
	err := initializeConnection()
	if err != nil {
		errorPage(w, err)
		return
	}
	vars := mux.Vars(r)

	var Product_req Product
	err = json.NewDecoder(r.Body).Decode(&Product_req)
	if err != nil {
		errorPage(w, err)
		return
	}

	var result Product

	hex, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{{Key: "_id",Value: hex}}
	update_fields := bson.D{{Key: "$set", Value: Product_req}}
	err = collection.FindOneAndUpdate(ctx, filter, update_fields, opts).Decode(&result)
	if err != nil {
		errorPage(w, err)
		return
	}
	
	res, err := json.Marshal(result)
	if err != nil {
		errorPage(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}