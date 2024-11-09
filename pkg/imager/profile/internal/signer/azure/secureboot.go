// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package azure

import (
	"context"
	"crypto"
	"crypto/x509"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/siderolabs/go-pointer"

	"github.com/siderolabs/talos/private/pkg/secureboot/pesign"
)

// SecureBootSigner implements pesign.CertificateSigner interface.
type SecureBootSigner struct {
	keySigner *KeySigner
	cert      *x509.Certificate
}

// Verify interface.
var _ pesign.CertificateSigner = (*SecureBootSigner)(nil)

// Signer returns the signer.
func (s *SecureBootSigner) Signer() crypto.Signer {
	return s.keySigner
}

// Certificate returns the certificate.
func (s *SecureBootSigner) Certificate() *x509.Certificate {
	return s.cert
}

// NewSecureBootSigner creates a new SecureBootSigner.
func NewSecureBootSigner(ctx context.Context, vaultURL, certificateID, certificateVersion string) (*SecureBootSigner, error) {
	certsClient, err := getCertsClient(vaultURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build Azure certificates client: %w", err)
	}

	resp, err := certsClient.GetCertificate(ctx, certificateID, certificateVersion, &azcertificates.GetCertificateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(resp.CER)
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate: %w", err)
	}

	// initialize key signer via existing implementation
	KID := pointer.SafeDeref(resp.KID)

	keySigner, err := NewPCRSigner(ctx, vaultURL, KID.Name(), KID.Version())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize certificate key signer: %w", err)
	}

	return &SecureBootSigner{
		cert:      cert,
		keySigner: keySigner,
	}, nil
}
