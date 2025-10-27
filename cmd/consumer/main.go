package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/consumer"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/infra/kmsmock"
	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/config"
	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/tokens"
)

func main() {
	ctx := context.Background()

	tokenFlag := flag.String("token", "", "confirmed token string")
	tokenFile := flag.String("token-file", "", "path to file containing token (alternative to --token)")
	capability := flag.String("capability", "coupons:redeem", "capability expected by the consumer")
	requiredConstraintsFlag := flag.String("require", "", "comma-separated list of constraint keys that must be present")
	domain := flag.String("domain", "coupons", "logical consumer domain label")
	flag.Parse()

	common, err := config.LoadCommon()
	if err != nil {
		log.Fatalf("consumer: load config: %v", err)
	}

	kms := kmsmock.NewService(common.KeyID, []byte(common.Secret))

	rawToken, err := readToken(*tokenFlag, *tokenFile, os.Stdin)
	if err != nil {
		log.Fatalf("consumer: read token: %v", err)
	}

	token, err := tokens.Decode(ctx, kms, rawToken)
	if err != nil {
		log.Fatalf("consumer: decode token: %v", err)
	}
	if err := token.Claims.Validate(time.Now()); err != nil {
		log.Fatalf("consumer: validate claims: %v", err)
	}

	service := consumer.NewService(consumer.Options{Domain: valueOrDefault(domain, "consumer")})

	requiredConstraints := config.ParseList(valueOrDefault(requiredConstraintsFlag, ""))

	result, err := service.Consume(ctx, consumer.ConsumeInput{
		Token:               token,
		Capability:          valueOrDefault(capability, ""),
		RequiredConstraints: requiredConstraints,
	})
	if err != nil {
		log.Fatalf("consumer: consume token: %v", err)
	}

	out := struct {
		Domain            string            `json:"domain"`
		OrderProcessingID string            `json:"order_processing_id"`
		TokenID           string            `json:"token_id"`
		ProcessedAt       time.Time         `json:"processed_at"`
		Constraints       map[string]string `json:"constraints,omitempty"`
	}{
		Domain:            result.Domain,
		OrderProcessingID: result.OrderProcessingID,
		TokenID:           result.TokenID,
		ProcessedAt:       result.ProcessedAt,
		Constraints:       result.Constraints,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		log.Fatalf("consumer: output: %v", err)
	}
}

func readToken(token, filePath string, stdin io.Reader) (string, error) {
	switch {
	case strings.TrimSpace(token) != "":
		return strings.TrimSpace(token), nil
	case strings.TrimSpace(filePath) != "":
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	default:
		line, err := bufio.NewReader(stdin).ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if strings.TrimSpace(line) == "" {
			return "", fmt.Errorf("no token supplied via --token, --token-file, or stdin")
		}
		return strings.TrimSpace(line), nil
	}
}

func valueOrDefault(ptr *string, fallback string) string {
	if ptr == nil || *ptr == "" {
		return fallback
	}
	return *ptr
}
