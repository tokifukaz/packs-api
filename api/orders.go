package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"packs-api/internal/resources"
	"packs-api/internal/services"
	"packs-api/internal/store"
)

func (s *Server) HandleCreateOrder(mongoDB store.NoSQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !s.HasContentType(r, "application/json") {
			s.WriteJSONError(w, http.StatusUnsupportedMediaType, "Unsupported media type")
			return
		}

		if r.Body == nil {
			s.WriteJSONError(w, http.StatusBadRequest, "Request body is empty")
			return
		}

		b, err := io.ReadAll(r.Body)
		defer func() {
			_ = r.Body.Close()
		}()
		if err != nil {
			s.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		var orderRequest *resources.OrderRequest
		err = json.Unmarshal(b, &orderRequest)
		if err != nil {
			s.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		packs := services.GetPacks(orderRequest.Items, orderRequest.PackSizes)

		if len(packs) < 1 {
			s.Log.Error("no packs to ship")
			s.WriteJSONError(w, http.StatusBadRequest, "no packs to ship")
			return
		}

		var order resources.Order
		now := s.Time.Now()
		order.ID = s.ObjectIDGenerator.GenerateRandomObjectID()
		order.Items = orderRequest.Items
		order.PackSizes = orderRequest.PackSizes
		order.PackQuantity = packs
		order.CreatedAt = now
		order.UpdatedAt = now

		err = mongoDB.CreateOrder(ctx, &order)
		if err != nil {
			err := fmt.Errorf("error creating order: %w", err)
			s.Log.WithField("error", err.Error()).Error("failed to create order")
			s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

}

func (s *Server) HandleGetAllOrders(mongoDB store.NoSQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		s.Log.Info("getting all orders")

		orders, err := mongoDB.GetAllOrders(ctx)
		if err != nil {
			err := fmt.Errorf("error getting all orders: %w", err)
			s.Log.WithField("error", err.Error()).Error("failed to get all orders")
			s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		s.Log.Info("orders retrieved successfully")

		data := map[string]interface{}{
			"data": orders,
		}

		jsonResponse, err := json.Marshal(data)
		if err != nil {
			s.Log.WithField("error", err.Error()).Error("invalid response body")
			s.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResponse)
	}
}
