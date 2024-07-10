/*
 * Licensed under the Apache License, Version 2.0 (the “License”);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package seed

import (
	"context"
	"log"
	"os"
	"time"

	"crapi.proj/goservice/api/models"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//initialize coupons data
var coupons = []models.Coupon{
	models.Coupon{
		CouponCode: "TRAC075",
		Amount:     "75",
		CreatedAt:  time.Now(),
	},
	models.Coupon{
		CouponCode: "TRAC065",
		Amount:     "65",
		CreatedAt:  time.Now(),
	},
	models.Coupon{
		CouponCode: "TRAC125",
		Amount:     "125",
		CreatedAt:  time.Now(),
	},
}

//initialize Post data
var posts = []models.Post{
	{
		Title:   "Forum rules by Bhupendra@5B",
		Content: "Hello Tuners! Welcome to our car tuner garage forum! To ensure everyone enjoys a positive and productive experience, we've put together a set of community guidelines: Respect Each Other: Treat fellow tuners with respect and courtesy. We're all here to learn and share our passion for car tuning. Stay On Topic: Keep discussions relevant to car tuning, modifications, performance enhancements, and related topics. Off-topic posts may be removed. No Spam or Advertising: Avoid spamming the forum with advertisements or promotional content. If you're unsure whether something qualifies as spam, reach out to a moderator. Be Helpful: Share your knowledge and experiences to help others. Constructive criticism is welcome, but be mindful of others' perspectives. No Personal Attacks or Hate Speech: Offensive language, personal attacks, or discriminatory remarks will not be tolerated. This includes any form of hate speech or bullying. Respect Privacy: Do not share personal or sensitive information about others without their consent. Use Clear and Respectful Language: Communicate clearly and use language that is appropriate for all ages and backgrounds. Report Issues: If you encounter a problem or notice a violation of these guidelines, report it to the moderators. Do not attempt to handle it yourself. Respect Copyright: Do not post copyrighted material without permission. If you're sharing content from other sources, give proper credit. Have Fun and Learn: Above all, enjoy your time here and learn from each other's experiences. We're a community passionate about car tuning! These guidelines are here to ensure our forum remains a friendly and informative space for all members. Let's build a supportive community where everyone can share their love for car tuning. Happy tuning! Best regards, Bhupendra",
	},
}
var emails = [3]string{"bhupendra@5b.com"}

//
func LoadMongoData(mongoClient *mongo.Client, db *gorm.DB) {
	var couponResult interface{}
	var postResult interface{}
	collection := mongoClient.Database(os.Getenv("MONGO_DB_NAME")).Collection("coupons")
	// get a MongoDB document using the FindOne() method
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&couponResult)
	if err != nil {
		for i := range coupons {
			couponData, err := collection.InsertOne(context.TODO(), coupons[i])
			log.Println(couponData, err)
		}
	}
	postCollection := mongoClient.Database(os.Getenv("MONGO_DB_NAME")).Collection("post")
	er := postCollection.FindOne(context.TODO(), bson.D{}).Decode(&postResult)
	if er != nil {
		for j := range posts {
			author, err := models.FindAuthorByEmail(emails[j], db)
			if err != nil {
				log.Println("Error finding author", err)
				continue
			}
			log.Println(author)
			posts[j].Prepare()
			postData, err := models.SavePost(mongoClient, posts[j]) // Assign the returned values to separate variables
			if err != nil {
				log.Println("Error saving post", err)
			}
			log.Println(postData) // Use the returned values as needed
		}
	}
	log.Println("Data seeded successfully")
}
