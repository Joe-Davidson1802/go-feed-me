package app

import (
	"basket/restaurants"
	"basket/templates"
	"basket/types"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/joe-davidson1802/hotwirehandler"
	"github.com/joe-davidson1802/turbo-templ/turbo"
)

var decoder = schema.NewDecoder()

type UpdateBasketHandler struct {
	Config types.Config
}

func (h UpdateBasketHandler) CanHandleModel(m string) bool {
	return m == types.Basket{}.ModelName()
}

func (h UpdateBasketHandler) HandleRequest(w http.ResponseWriter, r *http.Request) (error, hotwirehandler.Model) {
	err := r.ParseForm()

	if err != nil {
		return err, nil
	}

	var restaurant types.Restaurant

	err = decoder.Decode(&restaurant, r.PostForm)

	if err != nil {
		return err, nil
	}

	res := restaurants.RestaurantRepository{Config: h.Config}

	result, err := res.GetRestaurants()

	if err != nil {
		return err, nil
	}

	id, err := strconv.Atoi(restaurant.Id)

	if err != nil {
		return err, nil
	}

	restaurantdata := result[id-1]

	restaurant.Name = restaurantdata.Name

	for itemid, _ := range restaurant.Items {
		fmt.Println("ITEM")
		restaurant.Items[itemid].Price = restaurantdata.Items[itemid].Price
		restaurant.Items[itemid].Name = restaurantdata.Items[itemid].Name
		fmt.Println(restaurant.Items[itemid].Name)
	}

	basket := []types.Restaurant{}

	for _, res := range inmemorybasket {
		basket = append(basket, res)
	}

	basket = append(basket, restaurant)

	w.Header().Add("Cache-Control", "no-cache")

	return nil, types.Basket{
		Restaurants: basket,
	}
}

func (h UpdateBasketHandler) RenderPage(ctx context.Context, m hotwirehandler.Model, w http.ResponseWriter) error {
	mod := m.(types.Basket)

	w.Header().Add("Content-Type", "text/html")

	contents := templates.BasketComponent(mod)

	frame := turbo.TurboFrame(turbo.TurboFrameOptions{
		Id:       "basket",
		Contents: &contents,
	})

	err := frame.Render(ctx, w)

	return err
}

func (h UpdateBasketHandler) RenderStream(ctx context.Context, m hotwirehandler.Model, w http.ResponseWriter) error {
	mod := m.(types.Basket)

	w.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	contents := templates.BasketComponent(mod)

	stream := turbo.TurboStream(turbo.TurboStreamOptions{
		Action:   turbo.UpdateAction,
		Target:   "basket",
		Contents: &contents,
	})

	err := stream.Render(ctx, w)

	return err
}
