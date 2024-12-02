package sealer

import (
	"context"
	"crypto/rsa"
	"fmt"
	ssv1alpha1 "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	"github.com/bitnami-labs/sealed-secrets/pkg/client/clientset/versioned/scheme"
	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/keyutil"

	"os"
	"strings"
)

type SealerProvider struct {
	certificate       string
	privateKeyContent string
}

func NewSealerProvider(cfg *viper.Viper, logger *zap.Logger) *SealerProvider {
	lg := logger.Named("sealer-provider")
	certificate := strings.TrimSpace(cfg.GetString("app.sealer.certificate"))
	cert, _ := os.CreateTemp("", "*.crt")

	if _, err := cert.Write([]byte(certificate)); err != nil {
		lg.Fatal("failed to write certificate", zap.Error(err))
		return nil
	}

	lg.Info("Created certificate file",
		zap.String("certificate", cert.Name()))
	return &SealerProvider{
		certificate:       cert.Name(),
		privateKeyContent: strings.TrimSpace(cfg.GetString("app.sealer.privateKey")),
	}
}

func (c *SealerProvider) Seal(ctx context.Context, text string) (string, error) {
	cert, err := kubeseal.OpenCert(ctx, nil, "", "", c.certificate)
	if err != nil {
		return "", err
	}
	key, err := kubeseal.ParseKey(cert)
	if err != nil {
		return "", err
	}
	d, err := ssv1alpha1.NewSealedSecret(scheme.Codecs, key, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-project",
			Namespace: "pilot",
		},
		StringData: map[string]string{
			"config": text,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create sealed secret: %v", err)
	}
	return d.Spec.EncryptedData["config"], nil
}

func (c *SealerProvider) Unseal(text string) (string, error) {
	var sealedSecret = &ssv1alpha1.SealedSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config-project",
			Namespace: "pilot",
		},
		Spec: ssv1alpha1.SealedSecretSpec{
			EncryptedData: map[string]string{
				"config": text,
			},
		},
	}

	k, err := keyutil.ParsePrivateKeyPEM([]byte(c.privateKeyContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}
	pk := k.(*rsa.PrivateKey)
	secret, err := sealedSecret.Unseal(scheme.Codecs, map[string]*rsa.PrivateKey{
		"": pk,
	})
	if err != nil {
		return "", fmt.Errorf("failed to unseal secret: %v", err)
	}
	return string(secret.Data["config"]), nil
}
