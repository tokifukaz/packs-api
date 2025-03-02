package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"packs-api/internal/resources"
	"packs-api/internal/utils"
	"packs-api/mocks"
)

func TestServer_HandleGetAllOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	freezedTime := mocks.NewMockTime(ctrl)
	freezedTime.EXPECT().Now().Return(time.Date(2023, 11, 04, 20, 34, 58, 651387237, time.UTC)).AnyTimes()

	orderOneID := primitive.NewObjectID()
	orderTwoID := primitive.NewObjectID()

	orderOne := &resources.Order{
		ID:        orderOneID,
		Items:     10,
		PackSizes: []int{1, 2, 3},
		PackQuantity: map[int]int{
			1: 1,
			3: 3,
		},
		CreatedAt: freezedTime.Now(),
		UpdatedAt: freezedTime.Now(),
	}

	orderTwo := &resources.Order{
		ID:        orderTwoID,
		Items:     20,
		PackSizes: []int{1, 2, 3},
		PackQuantity: map[int]int{
			2: 1,
			3: 6,
		},
		CreatedAt: freezedTime.Now().AddDate(0, 0, 1),
		UpdatedAt: freezedTime.Now().AddDate(0, 0, 1),
	}

	orders := []*resources.Order{orderOne, orderTwo}

	mongoDB := mocks.NewMockNoSQLStore(ctrl)
	mongoDB.
		EXPECT().
		GetAllOrders(gomock.Any()).
		Return(nil, errors.New("store error")).
		Times(1)
	mongoDB.
		EXPECT().
		GetAllOrders(gomock.Any()).
		Return(orders, nil).
		Times(1)

	logger := utils.NewLogger("test", "packs-api")

	s := new(Server)
	s.Time = freezedTime
	s.Log = logger

	tests := []struct {
		name     string
		status   int
		errorMsg string
	}{
		{"error", 500, "error getting all orders: store error"},
		{"success", 200, ""},
	}

	for _, tt := range tests {
		ctx := context.Background()
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/api/orders", nil)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/api/orders", s.HandleGetAllOrders(mongoDB)).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, tt.status, rr.Code)

		b, _ := io.ReadAll(rr.Body)
		var res map[string]interface{}
		err := json.Unmarshal(b, &res)
		assert.Nil(t, err)

		var expected map[string]interface{}
		if tt.errorMsg == "" {
			expected = map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"createdAt": "2023-11-04T20:34:58.651387237Z",
						"id":        orderOneID.Hex(),
						"items":     float64(10),
						"packQuantity": map[string]interface{}{
							"1": float64(1),
							"3": float64(3),
						},
						"packSizes": []interface{}{float64(1), float64(2), float64(3)},
						"updatedAt": "2023-11-04T20:34:58.651387237Z",
					},
					map[string]interface{}{
						"createdAt": "2023-11-05T20:34:58.651387237Z",
						"id":        orderTwoID.Hex(),
						"items":     float64(20),
						"packQuantity": map[string]interface{}{
							"2": float64(1),
							"3": float64(6),
						},
						"packSizes": []interface{}{float64(1), float64(2), float64(3)},
						"updatedAt": "2023-11-05T20:34:58.651387237Z",
					},
				},
			}
		} else {
			expected = map[string]interface{}{"error": true, "code": float64(rr.Code), "message": tt.errorMsg}
		}
		assert.Equal(t, expected, res)
	}
}

func TestServer_HandleCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	freezedTime := mocks.NewMockTime(ctrl)
	freezedTime.EXPECT().Now().Return(time.Date(2023, 11, 04, 20, 34, 58, 651387237, time.UTC)).AnyTimes()

	orderID := primitive.NewObjectID()

	objectIDGenerator := mocks.NewMockObjectIDGenerator(ctrl)
	objectIDGenerator.
		EXPECT().
		GenerateRandomObjectID().
		Return(orderID).
		Times(2)

	orderOne := &resources.Order{
		ID:        orderID,
		Items:     10,
		PackSizes: []int{1, 2, 3},
		PackQuantity: map[int]int{
			1: 1,
			3: 3,
		},
		CreatedAt: freezedTime.Now(),
		UpdatedAt: freezedTime.Now(),
	}

	mongoDB := mocks.NewMockNoSQLStore(ctrl)
	mongoDB.
		EXPECT().
		CreateOrder(gomock.Any(), orderOne).
		Return(errors.New("store error")).
		Times(1)
	mongoDB.
		EXPECT().
		CreateOrder(gomock.Any(), orderOne).
		Return(nil).
		Times(1)

	logger := utils.NewLogger("test", "packs-api")

	s := new(Server)
	s.ObjectIDGenerator = objectIDGenerator
	s.Time = freezedTime
	s.Log = logger

	tests := []struct {
		name        string
		contentType string
		requestBody []byte
		status      int
		errorMsg    string
	}{
		{"unsupported media type", "application/xml", nil, 415, "Unsupported media type"},
		{"empty request body", "application/json", nil, 400, "unexpected end of JSON input"},
		{"invalid json", "application/json", []byte(`{"items": "a"}`), 400, "json: cannot unmarshal string into Go struct field OrderRequest.items of type int"},
		{"error creating order", "application/json", []byte(`{"items": 10, "packSizes": [1, 2, 3]}`), 500, "error creating order: store error"},
		{"success", "application/json", []byte(`{"items": 10, "packSizes": [1, 2, 3]}`), 201, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/orders", s.HandleCreateOrder(mongoDB)).Methods(http.MethodPost)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.status, rr.Code)

			if tt.errorMsg != "" {
				var res map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &res)
				assert.Nil(t, err)
				expected := map[string]interface{}{"error": true, "code": float64(tt.status), "message": tt.errorMsg}
				assert.Equal(t, expected, res)
			}
		})
	}
}
