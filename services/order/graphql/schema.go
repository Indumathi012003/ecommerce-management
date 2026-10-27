package graphql

import (
	gql "github.com/graphql-go/graphql"
	"order-management/model"
	"order-management/services/order/service"
)

type Server struct {
	service service.OrderService
	Schema  gql.Schema
}

func NewServer(s service.OrderService) (*Server, error) {
	x := &Server{service: s}
	sc, err := x.buildSchema()
	if err != nil {
		return nil, err
	}
	x.Schema = sc
	return x, nil
}

func (x *Server) buildSchema() (gql.Schema, error) {
	orderType := gql.NewObject(gql.ObjectConfig{
		Name: "Order",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.NewNonNull(gql.Int),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					if o, ok := p.Source.(*model.Order); ok {
						return int(o.ID), nil
					}
					return nil, nil
				},
			},
			"customerId":  &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"productId":   &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"quantity":    &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"totalAmount": &gql.Field{Type: gql.NewNonNull(gql.Float)},
			"status": &gql.Field{
				Type: gql.String,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					if o, ok := p.Source.(*model.Order); ok && o.Status != "" {
						return string(o.Status), nil
					}
					return string(model.StatusCreated), nil
				},
			},
		},
	})

	query := gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"order": &gql.Field{
				Type: orderType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					o, err := x.service.GetOrder(p.Context, id)
					if err != nil || service.IsNotFound(err) {
						return nil, nil
					}
					return o, nil
				},
			},
			"orders": &gql.Field{
				Type: gql.NewList(orderType),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return x.service.ListOrders(p.Context)
				},
			},
		},
	})

	mutation := gql.NewObject(gql.ObjectConfig{
		Name: "Mutation",
		Fields: gql.Fields{
			"createOrder": &gql.Field{
				Type: orderType,
				Args: gql.FieldConfigArgument{
					"customerId": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
					"productId":  &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
					"quantity":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					req := model.CreateOrderRequest{
						CustomerID: p.Args["customerId"].(int),
						ProductID:  p.Args["productId"].(int),
						Quantity:   p.Args["quantity"].(int),
					}
					return x.service.CreateOrder(p.Context, req)
				},
			},
			"updateOrderStatus": &gql.Field{
				Type: orderType,
				Args: gql.FieldConfigArgument{
					"id":     &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
					"status": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					st := model.OrderStatus(p.Args["status"].(string))
					return x.service.UpdateOrderStatus(p.Context, id, st)
				},
			},
		},
	})

	return gql.NewSchema(gql.SchemaConfig{Query: query, Mutation: mutation})
}