package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {

	host := "20.126.71.5"
	port := "22"
	user := ""
	pass := ""
	cmd  := "ps"

	hostKey := getHostKey(host)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout:         5 * time.Second,
	}

	// connect
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// start session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// setup standard out and error
	// uses writer interface
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}

}

func getHostKey(host string) ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

  fmt.Println(file)
  

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
