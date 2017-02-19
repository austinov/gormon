package ssh

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// some pieces of code taken from
// https://github.com/rapidloop/rtop
func addKeyAuth(auths []ssh.AuthMethod, keypath string) []ssh.AuthMethod {
	if len(keypath) == 0 {
		return auths
	}

	keypath = expandPath(keypath)

	// read the file
	pemBytes, err := ioutil.ReadFile(keypath)
	if err != nil {
		log.Fatal(err)
	}

	// get first pem block
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Fatalf("no key found in %s", keypath)
	}

	// handle plain and encrypted keyfiles
	if x509.IsEncryptedPEMBlock(block) {
		prompt := fmt.Sprintf("Enter passphrase for key '%s': ", keypath)
		pass, err := getpass(prompt)
		if err != nil {
			return auths
		}
		block.Bytes, err = x509.DecryptPEMBlock(block, []byte(pass))
		if err != nil {
			log.Fatalf("decrypt pem-block error %#v", err)
		}
		key, err := parsePemBlock(block)
		if err != nil {
			log.Fatalf("parse pem-block error %#v", err)
		}
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			log.Fatalf("signer from key error %#v", err)
		}
		return append(auths, ssh.PublicKeys(signer))
	} else {
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			log.Fatalf("private key error %#v", err)
		}
		return append(auths, ssh.PublicKeys(signer))
	}
}

func getpass(prompt string) (pass string, err error) {
	tstate, err := terminal.GetState(0)
	if err != nil {
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		quit := false
		for _ = range sig {
			quit = true
			break
		}
		terminal.Restore(0, tstate)
		if quit {
			fmt.Println()
			os.Exit(2)
		}
	}()
	defer func() {
		signal.Stop(sig)
		close(sig)
	}()

	f := bufio.NewWriter(os.Stdout)
	f.Write([]byte(prompt))
	f.Flush()

	passbytes, err := terminal.ReadPassword(0)
	pass = string(passbytes)

	f.Write([]byte("\n"))
	f.Flush()

	return
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "DSA PRIVATE KEY":
		return ssh.ParseDSAPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type %q", block.Type)
	}
}

func expandPath(path string) string {
	if len(path) < 2 || path[:2] != "~/" {
		return path
	}
	currUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(path, "~", currUser.HomeDir, 1)
}
