package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	shell "github.com/ipfs/go-ipfs-api"
)

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

func main() {
	// IPFS Gateway
	sh := shell.NewShell("localhost:5001") // Ensure IPFS daemon is running

	// Define CIDs
	kmeansCID := "QmX7sQ513XkRnHPsFbRNS2mZt1tZivvhUjHtLvnZMZ3UX5"
	inputCID := "QmRFFNceZPvTMqGAim1adGwQqgfz3ekAaNQGVdtFK6WQuq"

	// Temporary working directory
	tempDir := "./temp"
	os.MkdirAll(tempDir, 0755)

	// Download kmeans.go
	kmeansPath := filepath.Join(tempDir, "kmeans.go")
	if err := downloadFromIPFS(sh, kmeansCID, kmeansPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Download input.txt
	inputPath := filepath.Join(tempDir, "input.txt")
	if err := downloadFromIPFS(sh, inputCID, inputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Execute K-Means code
	if err := executeKMeans(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Upload output.txt to IPFS
	outputPath := filepath.Join(tempDir, "output.txt")
	outputCID, err := uploadToIPFS(sh, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display CIDs
	fmt.Println("Process Complete!")
	fmt.Printf("kmeans.go CID: %s\n", kmeansCID)
	fmt.Printf("input.txt CID: %s\n", inputCID)
	fmt.Printf("output.txt CID: %s\n", outputCID)
}
