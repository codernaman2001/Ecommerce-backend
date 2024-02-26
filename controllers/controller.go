package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codernaman2001/ecommerce/database"
	"github.com/codernaman2001/ecommerce/models"
	generate "github.com/codernaman2001/ecommerce/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Product")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

// func VerifyPassword(pass1 string, pass2 string) bool,string{

// }

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone number alreday registered"})
			return
		}

		// count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		// 	if err != nil {
		// 		log.Panic(err)
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// 		return
		// 	}
		// 	if count > 0 {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		// 	}
		// 	count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		// 	defer cancel()
		// 	if err != nil {
		// 		log.Panic(err)
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// 		return
		// 	}
		// 	if count > 0 {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
		// 		return
		// 	}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserted := UserCollection.InsertOne(ctx, user)

		if inserted != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": inserted.Error()})
		}

		c.JSON(http.StatusAccepted, "SUCCESSFULLY REGISTERED")

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var loginuser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&loginuser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *loginuser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*loginuser.Email, *loginuser.First_Name, *loginuser.Last_Name, loginuser.User_ID)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, loginuser.User_ID)
		c.JSON(http.StatusFound, loginuser)
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		searchdb, err := ProductCollection.Find(ctx, bson.D{})
		if err != nil {
			c.IndentedJSON(404, "there is a error")
			return
		}

		err = searchdb.All(ctx, &productlist)
		if err != nil {
			c.IndentedJSON(404, "error while fetching the data")
			return
		}

		if err := searchdb.Err(); err != nil {
			// Don't forget to log errors. I log them really simple here just
			// to get the point across.
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchdb.Close(ctx)
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.Product_ID = primitive.NewObjectID()
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		_, inserterr := ProductCollection.InsertOne(ctx, product)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusAccepted, product)
	}
}

// func SearchProductByQuery() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var productlist []models.Product
// 		queryparam := c.Query("name")
// 		if queryparam == "" {
// 			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
// 			c.Abort()
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		defer cancel()

// 		searchdb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryparam}})
// 		if err != nil {
// 			c.IndentedJSON(404, "Something went wrong")
// 			return
// 		}

// 		err = searchdb.All(ctx, &productlist)
// 		if err != nil {
// 			c.IndentedJSON(400, "invalid")
// 			return
// 		}

// 		defer searchdb.Close(ctx)
// 		defer cancel()

// 		c.IndentedJSON(200, productlist)
// 	}

// }

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchproducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchproducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}
