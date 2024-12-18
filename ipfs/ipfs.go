package ipfs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSManager struct {
	Shell *shell.Shell
	TempDir string
}

// NewIPFSManager creates a new IPFS manager instance
func NewIPFSManager(ipfsURL string) (*IPFSManager, error) {
	sh := shell.NewShell(ipfsURL)
	tempDir := "./temp/ipfs"
	
	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	return &IPFSManager{
		Shell: sh,
		TempDir: tempDir,
	}, nil
}

// UploadAlgorithm uploads an algorithm file to IPFS and returns its CID and hash
func (im *IPFSManager) UploadAlgorithm(algorithmPath string) (string, string, error) {
	// Upload to IPFS
	cid, err := uploadToIPFS(im.Shell, algorithmPath)
	if err != nil {
		return "", "", err
	}

	// Calculate hash of the algorithm
	hash, err := calculateFileHash(algorithmPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to calculate algorithm hash: %v", err)
	}

	return cid, hash, nil
}

// UploadData uploads input data to IPFS and returns its CID and hash
func (im *IPFSManager) UploadData(dataPath string) (string, string, error) {
	// Upload to IPFS
	cid, err := uploadToIPFS(im.Shell, dataPath)
	if err != nil {
		return "", "", err
	}

	// Calculate hash of the data
	hash, err := calculateFileHash(dataPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to calculate data hash: %v", err)
	}

	return cid, hash, nil
}

// ExecuteAlgorithm downloads algorithm and data from IPFS, executes it, and returns output hash
func (im *IPFSManager) ExecuteAlgorithm(algorithmCID, dataCID string) (string, error) {
	// Create temporary paths for downloaded files
	algorithmPath := filepath.Join(im.TempDir, "algorithm.go")
	dataPath := filepath.Join(im.TempDir, "input.data")
	
	// Download algorithm and data
	if err := downloadFromIPFS(im.Shell, algorithmCID, algorithmPath); err != nil {
		return "", err
	}
	if err := downloadFromIPFS(im.Shell, dataCID, dataPath); err != nil {
		return "", err
	}

	// Execute the algorithm
	outputPath := filepath.Join(im.TempDir, "output.data")
	if err := executeAlgorithm(algorithmPath, dataPath, outputPath); err != nil {
		return "", err
	}

	// Calculate hash of the output
	outputHash, err := calculateFileHash(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to calculate output hash: %v", err)
	}

	// Cleanup temporary files
	os.Remove(algorithmPath)
	os.Remove(dataPath)
	os.Remove(outputPath)

	return outputHash, nil
}

// VerifyOutput verifies if the calculated output hash matches the expected hash
func (im *IPFSManager) VerifyOutput(calculatedHash, expectedHash string) bool {
	return calculatedHash == expectedHash
}

// Helper function to calculate file hash
func calculateFileHash(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// Helper function to execute the algorithm
func executeAlgorithm(algorithmPath, inputPath, outputPath string) error {
	// Compile the algorithm
	execPath := filepath.Join(filepath.Dir(algorithmPath), "algorithm.exe")
	cmd := exec.Command("go", "build", "-o", execPath, algorithmPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile algorithm: %v", err)
	}
	defer os.Remove(execPath)

	// Execute the compiled algorithm
	cmd = exec.Command(execPath, inputPath, outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to execute algorithm: %v\nOutput: %s", err, output)
	}

	return nil
}

func downloadFromIPFS(sh *shell.Shell, cid, outputPath string) error {
	fmt.Printf("Downloading %s from IPFS...\n", cid)
	data, err := sh.Cat(cid)
	if err != nil {
		return fmt.Errorf("failed to download file from IPFS: %v", err)
	}
	defer data.Close()

	content, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	err = ioutil.WriteFile(outputPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("File downloaded to %s\n", outputPath)
	return nil
}

func uploadToIPFS(sh *shell.Shell, filePath string) (string, error) {
	fmt.Printf("Uploading %s to IPFS...\n", filePath)

	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Upload to IPFS
	cid, err := sh.Add(file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to IPFS: %v", err)
	}

	fmt.Printf("File uploaded with CID: %s\n", cid)
	return cid, nil
}

func executeKMeans() error {
	fmt.Println("Compiling kmeans.go...")

	// Ensure the 'temp' directory exists
	if err := os.MkdirAll("./temp", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Compile kmeans.go into an executable file named 'kmeans.exe' on Windows
	cmd := exec.Command("go", "build", "-o", "./temp/kmeans.exe", "kmeans.go")
	cmd.Dir = "./temp"
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile kmeans.go: %v\n%s", err, string(cmdOutput))
	}

	// Verify that the executable exists in the temp directory
	executablePath := filepath.Join("./temp", "kmeans.exe")
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		return fmt.Errorf("executable file not found at %s", executablePath)
	}

	// Run the compiled executable from the temp directory
	fmt.Println("Running kmeans...")
	cmd = exec.Command(executablePath)
	cmd.Dir = "./temp"
	cmdOutput, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute kmeans.exe: %v\n%s", err, string(cmdOutput))
	}

	fmt.Printf("K-Means executed successfully.\nOutput:\n%s\n", string(cmdOutput))
	return nil
}
