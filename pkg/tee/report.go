package tee

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	// RemoteAttestationProtocol is the protocol used for remote attestation
	RemoteAttestationProtocol       = "/remote_attestation_v1"
	RemoteAttestationProtocolCert   = RemoteAttestationProtocol + "/cert"
	RemoteAttestationProtocolReport = RemoteAttestationProtocol + "/report"
)

func VerifyNode(ctx context.Context, peerID peer.ID, node host.Host, signer []byte, production bool) error {
	// Open a stream
	streamCert, err := node.NewStream(ctx, peerID, RemoteAttestationProtocolCert)
	if err != nil {
		return fmt.Errorf("failed to open stream to get certificate: %w", err)
	}
	defer streamCert.Close()

	certificateBytes := []byte{}
	_, err = streamCert.Read(certificateBytes)
	if err != nil {
		return fmt.Errorf("failed to read certificate from stream: %w", err)
	}

	streamReport, err := node.NewStream(ctx, peerID, RemoteAttestationProtocolReport)
	if err != nil {
		return fmt.Errorf("failed to open stream to get report: %w", err)
	}
	defer streamReport.Close()

	reportBytes := []byte{}
	_, err = streamReport.Read(reportBytes)
	if err != nil {
		return fmt.Errorf("failed to read report from stream: %w", err)
	}

	return verifyReport(reportBytes, certificateBytes, signer, production)
}

func VerifierWrapper(ctx context.Context, node host.Host, signer []byte, production bool, streamHandler func(network.Stream)) network.StreamHandler {
	return func(s network.Stream) {
		if err := VerifyNode(ctx, s.Conn().RemotePeer(), node, signer, production); err != nil {
			fmt.Println("Failed to verify node:", err)
			s.Close()
			return
		}
		streamHandler(s)
	}
}

// RegisterRemoteAttestation registers the remote attestation protocol
// Using TEE and Intel SGX Enclaves with EGO ( https://github.com/edgelesssys/ego )
func RegisterRemoteAttestation(node host.Host) {
	node.SetStreamHandler(RemoteAttestationProtocolCert, func(s network.Stream) {
		defer s.Close()

		// Get the certificate for the remote node
		cert, err := getCertForNode(node)
		if err != nil {
			fmt.Println("Failed to get certificate for node:", err)
			return
		}
		// Send the report and certificate to the remote node
		_, err = s.Write(cert)
		if err != nil {
			fmt.Println("Failed to send report and certificate to node:", err)
			return
		}
	})
	node.SetStreamHandler(RemoteAttestationProtocolReport, func(s network.Stream) {
		defer s.Close()
		// Get the report for the remote node
		report, err := getReportForNode(node)
		if err != nil {
			fmt.Println("Failed to get report for node:", err)
			return
		}

		// Send the report and certificate to the remote node
		_, err = s.Write(report)
		if err != nil {
			fmt.Println("Failed to send report and certificate to node:", err)
			return
		}
	})
}

func getReportForNode(node host.Host) ([]byte, error) {
	pubkey, err := node.ID().ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to extract public key from p2p identity: %w", err)
	}
	rawKey, err := pubkey.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to extract public key from p2p identity: %w", err)
	}
	hash := sha256.Sum256(rawKey)
	return enclave.GetRemoteReport(hash[:])
}

func getCertForNode(node host.Host) ([]byte, error) {
	// If it's a random identity, get the pubkey from Libp2p
	// and convert these to Ethereum public Hex types
	pubkey, err := node.ID().ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to extract public key from p2p identity: %w", err)
	}
	return pubkey.Raw()
}

func verifyReport(reportBytes, certBytes, signer []byte, production bool) error {
	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
		// XXX: We'll ignore this issue for now. For an app that should run in production, you must decide which of the different TCBStatus values are acceptable for you to continue.")
		if production {
			return errors.New("TCB level is invalid")
		}
	} else if err != nil {
		return err
	}

	hash := sha256.Sum256(certBytes)
	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
		return errors.New("report data does not match the certificate's hash")
	}

	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).
	if report.SecurityVersion < 2 {
		return errors.New("invalid security version")
	}
	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
		return errors.New("invalid product")
	}
	if !bytes.Equal(report.SignerID, signer) {
		return errors.New("invalid signer")
	}

	// For production, you must also verify that report.Debug == false
	if production && report.Debug {
		return errors.New("debug is true")
	}

	return nil
}
