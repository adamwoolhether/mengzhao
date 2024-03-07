package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"

	"mengzhao/view/credits"
)

func CreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return render(w, r, credits.Index())
}

func StripeCheckout(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(chi.URLParam(r, "productID"))
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	checkoutParams := stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("___"),
				Quantity: stripe.Int64(1),
			},
		},
	}

	sesh, err := session.New(&checkoutParams)
	if err != nil {
		return err
	}

	http.Redirect(w, r, sesh.URL, http.StatusSeeOther)

	return nil
}
