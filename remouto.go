package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	log "github.com/TNTest/logrus"
)

func main() {
	ostype := runtime.GOOS

	var hostsPath string
	var err error
	switch ostype {
	default:
		hostsPath = "/etc/hosts"
	case "windows":
		hostsPath = os.Getenv("windir") + "\\System32\\drivers\\etc\\hosts"
	}
	log.Infof("OS type: %v", ostype)

	imoutoLocal := hostsPath + ".imouto"

	if _, err = os.Stat(hostsPath); err == nil {
		log.Infof("Found hosts file: %v", hostsPath)

		//get remote imouto hosts file
		remoteImoutoUrl := "https://raw.githubusercontent.com/phoenixlzx/imouto.host/master/imouto.host.txt"
		log.Infof("Begin to fetch imouto hosts from: %v", remoteImoutoUrl)
		resp, err := http.Get(remoteImoutoUrl)
		if err != nil {
			log.Errorf("getting imouto hosts file fail! %v", err)
			readyToExit(false)
		}
		defer resp.Body.Close()
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("getting imouto hosts content from %v fail! %v", remoteImoutoUrl, err)
			readyToExit(false)
		}
		log.Infof("Got imouto hosts file from: %v", remoteImoutoUrl)

		//read hosts current content
		hostsContent, err := ioutil.ReadFile(hostsPath)
		if err != nil {
			log.Errorf("reading hosts file fail!  %v", err)
			readyToExit(false)
		}
		log.Info("Read hosts file content.")

		hostsBase := hostsPath + ".base"
		var baseContent []byte
		if _, err := os.Stat(hostsBase); os.IsNotExist(err) {
			log.Info("hosts base file is not exists, creating it...")
			baseContent = hostsContent
			err = ioutil.WriteFile(hostsBase, baseContent, 0666)
			if err != nil {
				log.Errorf("Cannot write hosts base file. %v", err)
				readyToExit(false)
			}
			log.Info("hosts base file is caeated it.")

		} else {
			baseContent, err = ioutil.ReadFile(hostsBase)
			if err != nil {
				log.Errorf("Cannot read hosts base file. %v", err)
				readyToExit(false)
			}
			log.Info("Got hosts base file content.")

		}

		if _, err := os.Stat(imoutoLocal); os.IsNotExist(err) {
			log.Infof("no imouto local file: %s. Creating it.", imoutoLocal)

			err = ioutil.WriteFile(imoutoLocal, body, 0666)
			if err != nil {
				log.Errorf("Cannot write local imouto file. %v", err)
				readyToExit(false)
			}
			log.Info("local imouto file created.")

		} else {

			//read local imouto file
			var currentContent []byte
			currentContent, err = ioutil.ReadFile(imoutoLocal)
			if err != nil {
				log.Errorf("reading imouto hosts file fail! Exit program. %v", err)
				readyToExit(false)
			}
			log.Info("Got local imouto file content.")

			//compare local and remote imouto
			isEq := sliceEq(body, currentContent)
			if isEq {
				log.Info("no new version found. Exit.")
				os.Exit(0)
				//readyToExit(true)
			} else {
				log.Info("new version of imouto found. writing to local imouto file...")

				err = ioutil.WriteFile(imoutoLocal, body, 0666)
				if err != nil {
					log.Errorf("Cannot write local imouto file. %v", err)
					readyToExit(false)
				}
				log.Info("local imouto updated to new version.")

			}
		}

		hostsBKPath := hostsPath + ".bk"
		err = ioutil.WriteFile(hostsBKPath, hostsContent, 0666)
		if err != nil {
			log.Errorf("Cannot write hosts bk file. %v", err)
			readyToExit(false)
		}
		log.Info("Backup current hosts file finish.")

		//write hosts file with base content
		/*_, err = hostsFile.Write(baseContent)
		if err != nil {
			log.Errorf("Cannot write  base content to hosts file. %v", err)
			readyToExit(false)
		}*/
		/*err = ioutil.WriteFile(hostsPath, baseContent, 0644)
		if err != nil {
			log.Errorf("Cannot write hosts file with base content. %v", err)
			readyToExit(false)
		}
		log.Info("Writed to hosts file with base content.")*/

		//appand imouto file content to hosts file
		/*var hostsFile *os.File
		hostsFile, err = os.OpenFile(hostsPath, os.O_APPEND, 0644)
		if err != nil {
			log.Errorf("Cannot open hosts file. %v", err)
			readyToExit(false)
		}
		hostsFile.Close()
		log.Info("Open hosts file with appand model successful.")

		_, err = hostsFile.Write(body)
		if err != nil {
			log.Errorf("Cannot appand imouto content to hosts file. %v", err)
			readyToExit(false)
		}*/
		fullBody := concat(baseContent, body)
		err = ioutil.WriteFile(hostsPath, fullBody, 0644)
		if err != nil {
			log.Errorf("Cannot write hosts file with base content. %v", err)
			readyToExit(false)
		}
		log.Info("Update new version of imouto hosts content to hosts file.")
		//readyToExit(true)

	} else {
		log.Error("hosts file(%v) is not exists or cannot access! %v", hostsPath, err)
		readyToExit(false)
	}

}

func readyToExit(noErr bool) {
	log.Info("Press ** ENTER ** to exit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	if noErr {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func sliceEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func concat(head, tail []byte) []byte {
	c := make([]byte, len(head)+len(tail))
	copy(c[:len(head)], head)
	copy(c[len(head):], tail)
	return c
}
