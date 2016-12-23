package pay

import (
	log "github.com/Sirupsen/logrus"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
)

var ourPlan *stripe.Plan

func InitPayPlan() error {
	var err error
	// Set your secret key: remember to change this to your live secret key in production
	// See your keys here: https://dashboard.stripe.com/account/apikeys
	stripe.Key = "sk_test_bZsi2ACtaj2SgIoKtz3GxObe"

	params := &stripe.PlanParams{
		Name:     "Basic Plan",
		ID:       "basic-monthly",
		Interval: "month",
		Currency: "usd",
		Amount:   0,
	}
	ourPlan, err = plan.New(params)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to get plan")
		return err
	}
	return err
}

func CreatePayingCustomer(email string, ccNumber string,
	expiryMonth string, expiryYear string) (err error,
	customerPayID string) {
	custParams := &stripe.CustomerParams{
		Email: email,
		Source: &stripe.SourceParams{
			Card: &stripe.CardParams{
				Number: ccNumber,
				Month:  expiryMonth,
				Year:   expiryYear,
			},
		},
	}
	customer, err := customer.New(custParams)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to create paying customer")
		return err, ""
	}
	customerPayID = customer.ID
	log.WithFields(log.Fields{"payingcustomer": customer}).
		Debug("Successfully created paying customer")
	return err, customerPayID
}

func StartSubscription(payingCustomerID string) error {

	subParams := &stripe.SubParams{
		Customer: payingCustomerID,
		Plan:     "basic-monthly",
	}
	_, err := sub.New(subParams)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed" +
			" to create paying customer")
		return err
	}
	return err
}
