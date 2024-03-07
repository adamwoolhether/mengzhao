package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"

	"mengzhao/db"
	"mengzhao/view/credits"
)

func CreditsIndex(w http.ResponseWriter, r *http.Request) error {

	return render(w, r, credits.Index())
}

func StripeCheckout(w http.ResponseWriter, r *http.Request) error {
	productID := chi.URLParam(r, "productID")
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	checkoutParams := stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:7331/checkout/success/{CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:7331/checkout/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(productID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	sesh, err := session.New(&checkoutParams)
	if err != nil {
		return err
	}

	return htmxRedirect(w, r, sesh.URL)
}

func StripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	sessionID := chi.URLParam(r, "sessionID")
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	sesh, err := session.Get(sessionID, nil)
	if err != nil {
		return err
	}

	lineItemParams := stripe.CheckoutSessionListLineItemsParams{Session: stripe.String(sesh.ID)}
	iter := session.ListLineItems(&lineItemParams)
	iter.Next()
	item := iter.LineItem()
	priceID := item.Price.ID

	switch priceID {
	case os.Getenv("100_CREDITS"):
		user.Account.Credits += 100
	case os.Getenv("250_CREDITS"):

		user.Account.Credits += 250
	case os.Getenv("600_CREDITS"):
		user.Account.Credits += 600
	default:
		return fmt.Errorf("invalid price ID: %s", priceID)
	}

	if err := db.UpdateAccount(r.Context(), &user.Account); err != nil {
		return err
	}

	fmt.Println("updated account", user.Account)

	http.Redirect(w, r, "/generate", http.StatusSeeOther)

	return nil
}

func StripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {

	return nil
}
