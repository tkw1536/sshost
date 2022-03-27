package sshost

import "golang.org/x/crypto/ssh"

// This file contains variables that hold the names of algorithms supported by the ssh package.
// Some of these come from private, some from public constants.

// ssh.CertAlgo* constants
var knownCertAlgos = []string{
	ssh.CertAlgoRSAv01,
	ssh.CertAlgoDSAv01,
	ssh.CertAlgoECDSA256v01,
	ssh.CertAlgoECDSA384v01,
	ssh.CertAlgoECDSA521v01,
	ssh.CertAlgoSKECDSA256v01,
	ssh.CertAlgoED25519v01,
	ssh.CertAlgoSKED25519v01,
}

// ssh.KeyAlgoRSA* constants
var knownKeyAlgos = []string{
	ssh.KeyAlgoRSA,
	ssh.KeyAlgoDSA,
	ssh.KeyAlgoECDSA256,
	ssh.KeyAlgoSKECDSA256,
	ssh.KeyAlgoECDSA384,
	ssh.KeyAlgoECDSA521,
	ssh.KeyAlgoED25519,
	ssh.KeyAlgoSKED25519,
}

// known key exchange algorithms
// contained in private ssh.kexAlgo* constants
var knownKexAlgos = []string{
	"diffie-hellman-group1-sha1",
	"diffie-hellman-group14-sha1",
	"ecdh-sha2-nistp256",
	"ecdh-sha2-nistp384",
	"ecdh-sha2-nistp521",
	"curve25519-sha256@libssh.org",

	// client only
	"diffie-hellman-group-exchange-sha1",
	"diffie-hellman-group-exchange-sha256",
}

// known cipher names
// contained in ssh.supportedCiphers
var knownCiperNames = []string{
	"aes128-ctr", "aes192-ctr", "aes256-ctr",
	"aes128-gcm@openssh.com",
	"chacha20-poly1305@openssh.com",
	"arcfour256", "arcfour128", "arcfour",
	"aes128-cbc",
	"3des-cbc",
}

// knownMACnames
// contained in ssh.supportedMACs
var knownMACNames = []string{
	"hmac-sha2-256-etm@openssh.com", "hmac-sha2-256", "hmac-sha1", "hmac-sha1-96",
}
