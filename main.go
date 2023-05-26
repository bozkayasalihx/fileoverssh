package main

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
  key, err := ioutil.ReadFile("/home/malware/maker_key.pem")
  if err != nil {
    log.Fatalf("unable to read private key: %v", err)
  }

  signer ,err := ssh.ParsePrivateKey(key) 
  if err != nil {
    log.Fatalf("unable to parse private key: %v" ,err) 
  }

  config := &ssh.ClientConfig{
    User: "azureuser", 
    Auth: []ssh.AuthMethod{
      ssh.PublicKeys(signer),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    
  } 


  client , err := ssh.Dial("tcp", "20.126.71.5:22", config)
  if err != nil {
    log.Fatalf("unable to connect : %v", err)
  }

  session, err := client.NewSession()
  if err != nil {
    log.Fatalf("failed to create session: %v", err)
  }
  
  defer session.Close()

  session.Stdout = os.Stdout 
  session.Stdin  = os.Stdin 
  session.Stderr = os.Stderr

  if err := session.Shell(); err != nil {
    log.Fatalf("couldn't start shell on remote host: %v", err)
  }
  
  if err := session.Wait(); err != nil {
    log.Fatalf("failed to wait for shell: %v", err)
  }
}
