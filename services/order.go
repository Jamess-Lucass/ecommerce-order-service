package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Jamess-Lucass/ecommerce-order-service/emails"
	"github.com/Jamess-Lucass/ecommerce-order-service/middleware"
	"github.com/Jamess-Lucass/ecommerce-order-service/models"
	"github.com/Jamess-Lucass/ecommerce-order-service/requests"
	"github.com/Jamess-Lucass/ecommerce-order-service/utils"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService struct {
	db *mongo.Database
}

func NewOrderService(db *mongo.Database) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) List(ctx *fasthttp.RequestCtx, user *middleware.Claim) ([]*models.Order, error) {
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "userId", Value: uuid.MustParse(user.Subject)}},
			bson.D{{Key: "email", Value: user.Email}},
		}},
	}

	var orders []*models.Order
	cur, err := s.db.Collection("orders").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var order models.Order
		err := cur.Decode(&order)
		if err != nil {
			return orders, err
		}

		orders = append(orders, &order)
	}

	if err := cur.Err(); err != nil {
		return orders, err
	}

	cur.Close(ctx)

	if len(orders) == 0 {
		return []*models.Order{}, nil
	}

	return orders, nil
}

func (s *OrderService) Get(ctx *fasthttp.RequestCtx, user *middleware.Claim, id primitive.ObjectID) (*models.Order, error) {
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "userId", Value: uuid.MustParse(user.Subject)}},
			bson.D{{Key: "email", Value: user.Email}},
		}},
		{Key: "_id", Value: id},
	}

	var order models.Order
	if err := s.db.Collection("orders").FindOne(ctx, filter).Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) Create(ctx context.Context, order *models.Order) error {
	_, err := s.db.Collection("orders").InsertOne(ctx, order)

	return err
}

type CatalogResponse struct {
	Name string `json:"name"`
}

type UserResponse struct {
	Email     string `json:"email"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

var tmpl = template.Must(template.New("").Parse(emails.OrderPlacedEmail))

func (s *OrderService) SendPurchasedEmail(order *models.Order) error {
	if (order.Email == "" || order.Name == "") && order.UserId != uuid.Nil {
		// get the email from the user service
		uri := fmt.Sprintf("%s/api/v1/users/%s", os.Getenv("USER_SERVICE_BASE_URL"), order.UserId)
		body, err := utils.HttpGet(uri)
		if err != nil {
			return err
		}

		var user UserResponse
		if err := json.NewDecoder(body).Decode(&user); err != nil {
			return err
		}

		order.Email = user.Email
		order.Name = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	}

	purchaseEmail := emails.OrderPlaced{Name: order.Name, ID: order.ID.Hex(), Address: order.Address}
	for _, item := range order.Items {
		uri := fmt.Sprintf("%s/api/v1/catalog/%s", os.Getenv("CATALOG_SERVICE_BASE_URL"), item.CatalogId)
		body, err := utils.HttpGet(uri)
		if err != nil {
			return err
		}

		var catalog CatalogResponse
		if err := json.NewDecoder(body).Decode(&catalog); err != nil {
			return err
		}

		purchaseEmailItem := emails.OrderPlacedItem{
			Name:     catalog.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		}
		purchaseEmail.Items = append(purchaseEmail.Items, purchaseEmailItem)
	}

	var html bytes.Buffer
	if err := tmpl.Execute(&html, purchaseEmail); err != nil {
		return err
	}

	email := requests.Email{To: []string{order.Email}, From: "no-reply@ecommerce-order-service.test", Subject: "Your order is being processed", Body: html.String()}

	uri := fmt.Sprintf("%s/api/v1/emails", os.Getenv("EMAIL_SERVICE_BASE_URL"))

	body, err := json.Marshal(email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}

	return nil
}
