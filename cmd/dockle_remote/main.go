package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	l "log"
	"github.com/Portshift/dockle/pkg/log"
	"github.com/spf13/viper"
	"github.com/Portshift/klar/docker"
	"github.com/Portshift/klar/docker/token"
	dockle_config "github.com/Portshift/dockle/config"
	dockle_run "github.com/Portshift/dockle/pkg"
	dockle_types "github.com/Portshift/dockle/pkg/types"
	"github.com/containers/image/v5/docker/reference"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"
	"context"
)


func createDockleConfig(scanConfig *dockerfileScanConfig, imageName string) *dockle_config.Config {
	return &dockle_config.Config{
		Debug: scanConfig.verbose,
		Timeout: scanConfig.timeoutSec,
		ImageName: imageName,
		Username: scanConfig.username,
		Password: scanConfig.password,
	}
}

type dockerfileScanConfig struct {
	timeoutSec time.Duration
	verbose    bool
	scanUUID   string
	resultPath string
	username   string
	password   string
}

func loadScanConfig() (*dockerfileScanConfig, error) {
	var err error
	config := &dockerfileScanConfig{}

	config.verbose = viper.GetBool(verboseVar)
	config.scanUUID = viper.GetString(scanUUID)
	if len(config.scanUUID) == 0 {
		return nil, fmt.Errorf("no scan UUID (%v)", scanUUID)
	}
	config.resultPath = viper.GetString(resultPath)
	if len(config.resultPath) == 0 {
		return nil, fmt.Errorf("no result path (%v)", resultPath)
	}

	timeout, err := strconv.Atoi(viper.GetString(timeoutSecVar))
	if err != nil {
		return nil, fmt.Errorf("failed to convert timeout to int. value=%v", viper.GetString(timeoutSecVar))
	}

	config.timeoutSec = time.Duration(timeout) * time.Second

	return config, nil
}

func run(imageName string) {
	scanConfig, err := loadScanConfig()
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}
	if err = log.InitLogger(scanConfig.verbose); err != nil {
		l.Fatal(err)
	}

	ref, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		log.Logger.Fatalf("Failed to parse image name. name=%v: %v", imageName, err)
	}

	// strip tag if image has digest and tag
	ref = docker.ImageNameWithDigestOrTag(ref)
	// add default tag "latest"
	ref = reference.TagNameOnly(ref)
	adjustedImageName := ref.String()

	credExtractor := token.CreateCredExtractor()
	if scanConfig.username, scanConfig.password, err = credExtractor.GetCredentials(context.Background(), ref); err != nil {
		log.Logger.Fatalf("Failed to get credentials. image name=%v: %v", adjustedImageName, err)
	}

	scanResults := &dockle_types.ImageAssessment{
		Image:      imageName,
		ScanUUID:   scanConfig.scanUUID,
	}

	log.Logger.Infof("Scanning image %v", ref)

	assessmentMap, err := dockle_run.RunFromConfig(createDockleConfig(scanConfig, adjustedImageName))
	if err != nil {
		errMsg := fmt.Errorf("failed to run dockle: %w", err)
		log.Logger.Error(errMsg)
		scanResults.ScanErr = dockle_types.ConvertError(errMsg)
		scanResults.Success = false
	} else {
		scanResults.Success = true
		scanResults.Assessment = assessmentMap
		log.Logger.Infof("Image was successfully scanned")
	}

	err = sendScanResults(scanConfig.resultPath, scanResults)
	if err != nil {
		log.Logger.Fatalf("Failed send scan results: %v", err)
	}
}

func sendScanResults(resultServicePath string, scanResults *dockle_types.ImageAssessment) error {
	scanResultsB, err := json.Marshal(scanResults)
	if err != nil {
		return fmt.Errorf("failed marshal results: %v", err)
	}
	log.Logger.Infof("Sending results. results=%s", scanResultsB)

	req, err := http.NewRequest("POST", resultServicePath, bytes.NewBuffer(scanResultsB))
	if err != nil {
		return fmt.Errorf("failed to forward scan results: %v", err)
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			log.Logger.Errorf("Response Status: %s", resp.Status)
		}
		return err
	}
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Logger.Warnf("Failed to dump response: %v", err)
	} else {
		log.Logger.Debugf("Dumping response: %s", respDump)
	}

	return nil
}

const (
	verboseVar = "VERBOSE"
	timeoutSecVar = "TIMEOUT_SEC"
	scanUUID = "SCAN_UUID"
	resultPath = "RESULT_SERVICE_PATH"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Logger.Fatalf("Image name is required")
	}

	viper.SetDefault(verboseVar, "false")
	viper.SetDefault(timeoutSecVar, "90")
	viper.AutomaticEnv()

	run(os.Args[1])
}
