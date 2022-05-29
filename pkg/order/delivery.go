package order

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"html/template"
	"l0/pkg/domain"
	"log"
	"net/http"
	"strconv"
)

type OUsecase interface {
	GetOrderByID(id int) (*domain.Order, error)
	AddOrder(order *domain.Order) error
}

type API struct {
	ou OUsecase
}

func NewAPI(ou OUsecase) *API {
	return &API{ou}
}

var tmpl = template.Must(template.New("order").Parse(
	`
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Orders</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js" integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2" crossorigin="anonymous"></script>
</head>
<body>
    <div class="container">
        <div class="navbar">
            <h1 class="navbar-brand">Orders</h1>
        </div>

        <form method="get">
            <div class="mb-3">
                <label for="orderid" class="form-label">Input order id</label>
                <input type="text" class="form-control" name="orderid" id="orderid">
            </div>
            <button type="submit" class="btn btn-primary">Submit</button>
        </form>

        <h3>
            Order id: {{ .OrderUID }}
        </h3>
        <div>{{ . }}</div>
    </div>
</body>
</html>
`))

func (a *API) InputOrderIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("orderid")
	if id == "" {
		tmpl.Execute(w, nil)

		return
	}

	iid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	order, err := a.ou.GetOrderByID(iid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	tmpl.Execute(w, order)
}

func (a *API) SubscribeToOrders(m *stan.Msg) {
	ord, err := orderFromJSON(m.Data)
	if err != nil {
		log.Println("error:", err)

		return
	}

	if err = a.ou.AddOrder(ord); err != nil {
		log.Println(err)
	}
}

func orderFromJSON(data []byte) (*domain.Order, error) {
	order := &domain.Order{}

	if err := json.Unmarshal(data, order); err != nil {
		return nil, fmt.Errorf("error unmarshalling order: %w", err)
	}

	return order, nil
}
