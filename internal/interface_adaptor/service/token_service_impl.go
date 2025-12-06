package service

import (
	"context"
	"fmt"
	"strings"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/domain"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/service"
)

type tokenServiceImpl struct{}

func NewTokenService() service.TokenService {
	return &tokenServiceImpl{}
}

func (s *tokenServiceImpl) IssuePaymentToken(ctx context.Context, input service.IssuePaymentTokenInput) (service.SignedToken, error) {
	return service.SignedToken{Value: fmt.Sprintf("payment-token:%s:%s", input.OrderProcessingID, input.PaymentMethod)}, nil
}

func (s *tokenServiceImpl) RelayTokens(ctx context.Context, input service.RelayTokensInput) (map[string]service.SignedToken, error) {
	return map[string]service.SignedToken{
		"cart":    {Value: input.CartToken},
		"payment": {Value: input.PaymentToken},
	}, nil
}

func (s *tokenServiceImpl) ConfirmCartToken(ctx context.Context, input service.ConfirmCartTokenInput) (service.SignedToken, error) {
	items := make([]string, len(input.Cart.Items))
	for i, item := range input.Cart.Items {
		items[i] = string(item)
	}
	return service.SignedToken{
		Value: fmt.Sprintf(
			"cart-token:%s:%s:%s",
			input.Cart.BusinessID,
			input.Cart.CartID,
			strings.Join(items, "|"),
		),
	}, nil
}

func (s *tokenServiceImpl) ParseCartToken(ctx context.Context, token service.SignedToken) (domain.Cart, error) {
	if len(token.Value) == 0 {
		return domain.Cart{}, fmt.Errorf("cart token is empty")
	}
	if !strings.HasPrefix(token.Value, "cart-token:") {
		return domain.Cart{}, fmt.Errorf("invalid cart token")
	}

	parts := strings.SplitN(token.Value, ":", 4)
	if len(parts) != 4 {
		return domain.Cart{}, fmt.Errorf("invalid cart token format")
	}

	businessID := domain.BusinessID(parts[1])
	cartID := domain.CartID(parts[2])
	itemPart := parts[3]
	if len(itemPart) == 0 {
		return domain.Cart{}, fmt.Errorf("invalid cart token: empty items")
	}

	itemStrings := strings.Split(itemPart, "|")
	items := make([]domain.ItemID, 0, len(itemStrings))
	for _, item := range itemStrings {
		if len(item) == 0 {
			continue
		}
		items = append(items, domain.ItemID(item))
	}

	cartItems := domain.CartItems(items)
	cart := domain.Cart{
		BusinessID: businessID,
		CartID:     cartID,
		Items:      cartItems,
	}

	return cart, cart.Validate()
}
