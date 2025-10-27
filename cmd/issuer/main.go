package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/infra/kmsmock"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/issuer"
	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/config"
)

func main() {
	ctx := context.Background()

	subject := flag.String("subject", "user_123", "subject (sub) claim")
	orderID := flag.String("order-id", "op_demo", "order processing identifier")
	issuerID := flag.String("issuer", "coupons.svc", "issuer (iss) claim")
	audience := flag.String("audience", "order-processing", "audience (aud) claim")
	capabilitiesFlag := flag.String("capabilities", "coupons:redeem", "comma-separated capabilities to embed")
	constraintsFlag := flag.String("constraints", "", "comma-separated key=value constraints")
	kidOverride := flag.String("kid", "", "override KMS key identifier")
	ttl := flag.Duration("ttl", 5*time.Minute, "token time-to-live (e.g. 2m)")
	flag.Parse()

	common, err := config.LoadCommon()
	if err != nil {
		log.Fatalf("issuer: load config: %v", err)
	}

	keyID := common.KeyID
	if kidOverride != nil && *kidOverride != "" {
		keyID = *kidOverride
	}

	kms := kmsmock.NewService(keyID, []byte(common.Secret))

	service, err := issuer.NewService(kms, issuer.Options{
		Issuer:     *issuerID,
		Audience:   *audience,
		KeyID:      keyID,
		DefaultTTL: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("issuer: init service: %v", err)
	}

	capabilities := config.ParseList(valueOrDefault(capabilitiesFlag, ""))
	if len(capabilities) == 0 {
		log.Fatal("issuer: at least one capability is required (set --capabilities)")
	}

	constraints, err := config.ParseAssignments(valueOrDefault(constraintsFlag, ""))
	if err != nil {
		log.Fatalf("issuer: parse constraints: %v", err)
	}

	token, err := service.IssueConfirmedToken(ctx, issuer.IssueRequest{
		Subject:           valueOrDefault(subject, ""),
		OrderProcessingID: valueOrDefault(orderID, ""),
		Capabilities:      capabilities,
		Constraints:       constraints,
		TTL:               durationOrZero(ttl),
	})
	if err != nil {
		log.Fatalf("issuer: issue token: %v", err)
	}

	fmt.Println(token.Raw)
}

func valueOrDefault(ptr *string, fallback string) string {
	if ptr == nil || *ptr == "" {
		return fallback
	}
	return *ptr
}

func durationOrZero(ptr *time.Duration) time.Duration {
	if ptr == nil {
		return 0
	}
	return *ptr
}
