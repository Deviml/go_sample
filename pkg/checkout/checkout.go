package checkout

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"gorm.io/gorm"
)

const cost = 5500

type Service struct {
	logger log.Logger
	client *sql.DB
	db     *gorm.DB
	sc     *client.API
	s      email.Sender
}

func NewService(logger log.Logger, client *sql.DB, db *gorm.DB, sc *client.API, s email.Sender) *Service {
	return &Service{logger: logger, client: client, db: db, sc: sc, s: s}
}

func (s Service) Charge(amount int64, token string) (*stripe.Charge, error) {
	return s.sc.Charges.New(&stripe.ChargeParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("EQUIPHUNTER PAYMENT"),
		Source:      &stripe.SourceParams{Token: stripe.String(token)},
	})
}

func (s Service) Checkout(ctx context.Context, userID string, request vendorpurchases.CheckoutRequest) error {
	// amount := len(request.Quotes) + len(request.Sublists)
	cost := int(request.Total * 100)

	if request.Total == 0 {
		amount := len(request.Quotes) + len(request.Sublists)

		for _, id := range request.Quotes {
			err := s.AddUserToQuote(ctx, userID, id)
			if err != nil {
				return err
			}
		}
		for _, id := range request.Sublists {
			err := s.AddUserToSublist(ctx, userID, id)
			if err != nil {
				return err
			}
		}
		token, _ := json.Marshal(request.PaymentResponse)
		uid, _ := strconv.Atoi(userID)
		payment := models.Payment{
			Payload:   token,
			WebUserID: uid,
			Total:     request.Total,
			Status:    models.PayedStatus,
			Response:  nil,
		}
		result := s.db.Create(&payment)
		if result.Error != nil {
			s.logger.Log("err", result.Error.Error())
			return result.Error
		}
		var quote models.Quote

		for _, id := range request.Sublists {
			result := s.db.Create(&models.PaymentDetail{
				PaymentID: int(payment.ID),
				SublistID: &id,
			})
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}
		}

		for _, id := range request.Quotes {
			qid, _ := strconv.Atoi(id)
			result := s.db.Create(&models.PaymentDetail{
				PaymentID: int(payment.ID),
				QuoteID:   &qid,
			})
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}

			result = s.db.Model(&quote).Where("id = ?", &qid).Update("status", models.PurchasedStatus)
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}
		}

		err := s.SendReceipt(ctx, userID, request, cost*amount)
		if err != nil {
			return err
		}

		return nil
	} else {
		charge, err := s.Charge(int64(cost), request.Token["id"].(string))
		if err != nil {
			s.logger.Log("err", "checkout", err)
			return err
		}

		for _, id := range request.Quotes {
			err := s.AddUserToQuote(ctx, userID, id)
			if err != nil {
				return err
			}
		}
		for _, id := range request.Sublists {
			err := s.AddUserToSublist(ctx, userID, id)
			if err != nil {
				return err
			}
		}
		token, _ := json.Marshal(request.Token)
		uid, _ := strconv.Atoi(userID)
		payment := models.Payment{
			Payload:    token,
			WebUserID:  uid,
			Email:      request.Email,
			NameOnCard: request.NameOnCard,
			Zipcode:    request.Zipcode,
			Total:      request.Total,
			Status:     models.PayedStatus,
			Response:   charge.LastResponse.RawJSON,
		}
		result := s.db.Create(&payment)
		if result.Error != nil {
			s.logger.Log("err", result.Error.Error())
			return result.Error
		}
		var quote models.Quote

		for _, id := range request.Sublists {
			result := s.db.Create(&models.PaymentDetail{
				PaymentID: int(payment.ID),
				SublistID: &id,
			})
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}
		}

		for _, id := range request.Quotes {
			qid, _ := strconv.Atoi(id)
			result := s.db.Create(&models.PaymentDetail{
				PaymentID: int(payment.ID),
				QuoteID:   &qid,
			})
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}

			result = s.db.Model(&quote).Where("id = ?", &qid).Update("status", models.PurchasedStatus)
			if result.Error != nil {
				s.logger.Log("err", result.Error.Error())
				return result.Error
			}
		}

		err = s.SendReceipt(ctx, userID, request, cost)

		return nil
	}
}

func (s Service) PaymentIntent(ctx context.Context, userID string, request vendorpurchases.PaymentIntentRequest) (interface{}, error) {
	amount := len(request.Quotes) + len(request.Sublists)

	paymentIntent, err := s.sc.PaymentIntents.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * cost)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	})

	return paymentIntent, err
}

func (s Service) ConfirmPaymentIntentAndCheckout(ctx context.Context, userID string, request vendorpurchases.CheckoutRequest) error {

	amount := len(request.Quotes) + len(request.Sublists)

	// pi, err := s.sc.PaymentIntents.Confirm(
	// 	request.PaymentResponse["paymentIntent"].(map[string]interface{})["id"].(string),
	// 	&stripe.PaymentIntentConfirmParams{
	// 		PaymentMethod: stripe.String(request.PaymentResponse["paymentIntent"].(map[string]interface{})["payment_method"].(string)),
	// 	},
	// )
	// if err != nil {
	// 	s.logger.Log("err", "checkout", err)
	// 	return err
	// }

	// if pi.Status == stripe.PaymentIntentStatusSucceeded {

	for _, id := range request.Quotes {
		err := s.AddUserToQuote(ctx, userID, id)
		if err != nil {
			return err
		}
	}
	for _, id := range request.Sublists {
		err := s.AddUserToSublist(ctx, userID, id)
		if err != nil {
			return err
		}
	}
	token, _ := json.Marshal(request.PaymentResponse)
	uid, _ := strconv.Atoi(userID)
	payment := models.Payment{
		Payload:   token,
		WebUserID: uid,
		Total:     request.Total,
		Status:    models.PayedStatus,
		Response:  nil,
	}
	result := s.db.Create(&payment)
	if result.Error != nil {
		s.logger.Log("err", result.Error.Error())
		return result.Error
	}
	var quote models.Quote

	for _, id := range request.Sublists {
		result := s.db.Create(&models.PaymentDetail{
			PaymentID: int(payment.ID),
			SublistID: &id,
		})
		if result.Error != nil {
			s.logger.Log("err", result.Error.Error())
			return result.Error
		}
	}

	for _, id := range request.Quotes {
		qid, _ := strconv.Atoi(id)
		result := s.db.Create(&models.PaymentDetail{
			PaymentID: int(payment.ID),
			QuoteID:   &qid,
		})
		if result.Error != nil {
			s.logger.Log("err", result.Error.Error())
			return result.Error
		}

		result = s.db.Model(&quote).Where("id = ?", &qid).Update("status", models.PurchasedStatus)
		if result.Error != nil {
			s.logger.Log("err", result.Error.Error())
			return result.Error
		}
	}

	err := s.SendReceipt(ctx, userID, request, cost*amount)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) SendReceipt(ctx context.Context, userID string, request vendorpurchases.CheckoutRequest, amount int) error {
	amountf := float32(amount) / 100
	formattedAmount := fmt.Sprintf("%.2f", amountf)
	var ws models.WebUser
	result := s.db.Find(&ws, "id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	purchases := make([]templates.Purchase, 0)
	for _, id := range request.Quotes {
		var q models.Quote
		result := s.db.Preload("EquipmentRequest.Equipment").Preload("SupplyRequest.Supply").Find(&q, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		description := ""
		if q.EquipmentRequest != nil {
			description = q.EquipmentRequest.Equipment.Name
		} else {
			description = q.SupplyRequest.Supply.Name
		}
		purchases = append(purchases, templates.Purchase{
			Description: description,
			Type:        "Quote",
		})
	}

	for _, id := range request.Sublists {
		var sb models.Sublist
		result := s.db.Find(&sb, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		purchases = append(purchases, templates.Purchase{
			Description: sb.ProjectName,
			Type:        "Sublist",
		})
	}

	return s.s.SendReceipt(ctx, ws.Username, ws.FullName, purchases, formattedAmount)
}

func (s Service) AddUserToSublist(ctx context.Context, userID string, sublistID int) error {
	query, args, err := sq.Insert("web_user_sublists").Columns("sublist_id", "web_user_id").Values(sublistID, userID).ToSql()
	if err != nil {
		return err
	}
	_, err = s.client.Exec(query, args...)
	return err
}

func (s Service) AddUserToQuote(ctx context.Context, userID string, quoteID string) error {
	var wu models.WebUser
	result := s.db.Find(&wu, "id = ?", userID)
	if result != nil && result.Error != nil {
		return result.Error
	}
	var user models.VendorRental
	id, _ := strconv.Atoi(userID)
	result = s.db.Find(&user, "id = ?", wu.ProfileID)
	if result.Error != nil {
		return result.Error
	}
	id, _ = strconv.Atoi(quoteID)
	err := s.db.Model(&user).Association("Quotes").Append(&models.Quote{
		Model: gorm.Model{
			ID: uint(id),
		},
	})
	return err
}
