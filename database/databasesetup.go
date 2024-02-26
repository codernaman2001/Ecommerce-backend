// package database

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func DBSet() *mongo.Client {
// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://namanbansal:naman123@cluster0.bddogmg.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

// 	defer cancel()

// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Println("Failed to connect to mongodb:")
// 		return nil
// 	}

// 	fmt.Println("Successfully connected to mongodb")

// 	return client
// }

// var Client *mongo.Client = DBSet()

// func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
// 	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
// 	return collection
// }

// func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
// 	var productcollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
// 	return productcollection
// }

package database

//
import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://namanbansal:naman123@cluster0.bqjocfd.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println("Successfully Connected to the mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return collection

}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return productcollection
}
