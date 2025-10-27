package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/infra/kmsmock"
	"yoshiyoshifujii/go-capability-token-relay-pattern/internal/relay"
	"yoshiyoshifujii/go-capability-token-relay-pattern/pkg/config"
)

func main() {
	ctx := context.Background()

	tokenFlag := flag.String("token", "", "confirmed token string")
	tokenFile := flag.String("token-file", "", "path to file containing token (alternative to --token)")
	flag.Parse()

	common, err := config.LoadCommon()
	if err != nil {
		log.Fatalf("relay: load config: %v", err)
	}

	kms := kmsmock.NewService(common.KeyID, []byte(common.Secret))

	service, err := relay.NewService(kms, nil)
	if err != nil {
		log.Fatalf("relay: init service: %v", err)
	}

	rawToken, err := readToken(*tokenFlag, *tokenFile, os.Stdin)
	if err != nil {
		log.Fatalf("relay: read token: %v", err)
	}

	token, err := service.Verify(ctx, rawToken)
	if err != nil {
		log.Fatalf("relay: verify token: %v", err)
	}

	fmt.Fprintf(os.Stderr, "relay: token %s verified for order %s\n", token.Claims.TokenID, token.Claims.OrderProcessingID)
	fmt.Println(token.Raw)
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
