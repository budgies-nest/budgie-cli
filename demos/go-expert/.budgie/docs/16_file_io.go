/*
METADATA:
Description: Demonstrates Go file I/O operations including reading, writing, file manipulation, and directory operations
Keywords: file, io, read, write, open, close, os, ioutil, bufio, filepath, directory
Category: file-operations
Concepts: file reading/writing, buffered I/O, file permissions, directory traversal, file manipulation
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Function demonstrating basic file operations
func basicFileOperations() {
	fmt.Println("=== BASIC FILE OPERATIONS ===")
	
	filename := "example.txt"
	content := "Hello, File I/O!\nThis is a test file.\nIt has multiple lines."
	
	// Write to file
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	fmt.Printf("Successfully wrote to %s\n", filename)
	
	// Read entire file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	fmt.Printf("File contents:\n%s\n", string(data))
	
	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}
	fmt.Printf("File info: Name=%s, Size=%d bytes, Mode=%v\n", 
		info.Name(), info.Size(), info.Mode())
	
	// Clean up
	os.Remove(filename)
}

// Function demonstrating file opening and closing
func fileOpenClose() {
	fmt.Println("\n=== FILE OPENING AND CLOSING ===")
	
	filename := "open_close_example.txt"
	
	// Create and write to file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close() // Always close files
	
	// Write to file
	content := "This file demonstrates opening and closing.\n"
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
	
	// Write more content
	_, err = file.Write([]byte("This is written as bytes.\n"))
	if err != nil {
		fmt.Printf("Error writing bytes: %v\n", err)
		return
	}
	
	// Sync to ensure data is written
	err = file.Sync()
	if err != nil {
		fmt.Printf("Error syncing file: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully wrote to %s\n", filename)
	
	// Close file explicitly to read it
	file.Close()
	
	// Open file for reading
	file, err = os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file for reading: %v\n", err)
		return
	}
	defer file.Close()
	
	// Read file content
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	fmt.Printf("Read %d bytes:\n%s", n, string(buffer[:n]))
	
	// Clean up
	os.Remove(filename)
}

// Function demonstrating buffered I/O
func bufferedIO() {
	fmt.Println("\n=== BUFFERED I/O ===")
	
	filename := "buffered_example.txt"
	
	// Write using buffered writer
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	defer writer.Flush() // Always flush buffered writer
	
	lines := []string{
		"Line 1: Buffered writing is efficient",
		"Line 2: Especially for multiple writes",
		"Line 3: Buffer reduces system calls",
		"Line 4: Don't forget to flush!",
	}
	
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Error writing line: %v\n", err)
			return
		}
	}
	
	fmt.Printf("Wrote %d lines using buffered writer\n", len(lines))
	writer.Flush() // Ensure data is written
	file.Close()
	
	// Read using buffered reader
	file, err = os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	reader := bufio.NewReader(file)
	lineCount := 0
	
	fmt.Println("Reading with buffered reader:")
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(line) > 0 {
				fmt.Printf("  %s", line)
				lineCount++
			}
			break
		}
		if err != nil {
			fmt.Printf("Error reading line: %v\n", err)
			break
		}
		fmt.Printf("  %s", line)
		lineCount++
	}
	
	fmt.Printf("Read %d lines using buffered reader\n", lineCount)
	
	// Clean up
	os.Remove(filename)
}

// Function demonstrating scanning files
func scanningFiles() {
	fmt.Println("\n=== SCANNING FILES ===")
	
	filename := "scanning_example.txt"
	content := "apple banana\ncherry date elderberry\nfig grape honeydew\n"
	
	// Create file
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	defer os.Remove(filename)
	
	// Open file for scanning
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Scan by lines
	fmt.Println("Scanning by lines:")
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("  Line %d: %s\n", lineNumber, line)
		lineNumber++
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
	}
	
	// Reset file position and scan by words
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	
	fmt.Println("Scanning by words:")
	wordCount := 0
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Printf("  Word %d: %s\n", wordCount+1, word)
		wordCount++
	}
	
	fmt.Printf("Total words: %d\n", wordCount)
}

// Function demonstrating file permissions and attributes
func filePermissions() {
	fmt.Println("\n=== FILE PERMISSIONS ===")
	
	filename := "permissions_example.txt"
	
	// Create file with specific permissions
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	file.WriteString("File with 0600 permissions (owner read/write only)\n")
	file.Close()
	
	// Check file permissions
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}
	fmt.Printf("File permissions: %v\n", info.Mode())
	
	// Change file permissions
	err = os.Chmod(filename, 0644)
	if err != nil {
		fmt.Printf("Error changing permissions: %v\n", err)
		return
	}
	
	// Check updated permissions
	info, err = os.Stat(filename)
	if err != nil {
		fmt.Printf("Error getting updated file info: %v\n", err)
		return
	}
	fmt.Printf("Updated permissions: %v\n", info.Mode())
	
	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File %s exists\n", filename)
	} else if os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", filename)
	}
	
	// Clean up
	os.Remove(filename)
}

// Function demonstrating directory operations
func directoryOperations() {
	fmt.Println("\n=== DIRECTORY OPERATIONS ===")
	
	dirName := "example_directory"
	
	// Create directory
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	fmt.Printf("Created directory: %s\n", dirName)
	
	// Create nested directories
	nestedDir := filepath.Join(dirName, "nested", "deep")
	err = os.MkdirAll(nestedDir, 0755)
	if err != nil {
		fmt.Printf("Error creating nested directories: %v\n", err)
		return
	}
	fmt.Printf("Created nested directories: %s\n", nestedDir)
	
	// Create some files in directories
	files := []string{
		filepath.Join(dirName, "file1.txt"),
		filepath.Join(dirName, "file2.txt"),
		filepath.Join(nestedDir, "deep_file.txt"),
	}
	
	for _, filename := range files {
		err := os.WriteFile(filename, []byte("Sample content"), 0644)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Created file: %s\n", filename)
	}
	
	// List directory contents
	fmt.Printf("\nListing contents of %s:\n", dirName)
	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  [DIR]  %s\n", entry.Name())
		} else {
			fmt.Printf("  [FILE] %s\n", entry.Name())
		}
	}
	
	// Walk directory tree
	fmt.Printf("\nWalking directory tree starting from %s:\n", dirName)
	err = filepath.WalkDir(dirName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		relPath, _ := filepath.Rel(dirName, path)
		if d.IsDir() {
			fmt.Printf("  [DIR]  %s\n", relPath)
		} else {
			fmt.Printf("  [FILE] %s\n", relPath)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
	
	// Clean up
	os.RemoveAll(dirName)
	fmt.Printf("Cleaned up directory: %s\n", dirName)
}

// Function demonstrating file copying
func fileCopying() {
	fmt.Println("\n=== FILE COPYING ===")
	
	sourceFile := "source.txt"
	destFile := "destination.txt"
	content := "This content will be copied from source to destination."
	
	// Create source file
	err := os.WriteFile(sourceFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating source file: %v\n", err)
		return
	}
	defer os.Remove(sourceFile)
	defer os.Remove(destFile)
	
	// Open source file
	src, err := os.Open(sourceFile)
	if err != nil {
		fmt.Printf("Error opening source file: %v\n", err)
		return
	}
	defer src.Close()
	
	// Create destination file
	dst, err := os.Create(destFile)
	if err != nil {
		fmt.Printf("Error creating destination file: %v\n", err)
		return
	}
	defer dst.Close()
	
	// Copy file content
	bytesWritten, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	
	fmt.Printf("Copied %d bytes from %s to %s\n", bytesWritten, sourceFile, destFile)
	
	// Verify copy
	copiedContent, err := os.ReadFile(destFile)
	if err != nil {
		fmt.Printf("Error reading copied file: %v\n", err)
		return
	}
	
	if string(copiedContent) == content {
		fmt.Println("File copied successfully!")
	} else {
		fmt.Println("File copy verification failed!")
	}
}

// Function demonstrating temporary files
func temporaryFiles() {
	fmt.Println("\n=== TEMPORARY FILES ===")
	
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "example_*.txt")
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name()) // Clean up
	defer tmpFile.Close()
	
	fmt.Printf("Created temporary file: %s\n", tmpFile.Name())
	
	// Write to temporary file
	content := "This is temporary content that will be deleted."
	_, err = tmpFile.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to temp file: %v\n", err)
		return
	}
	
	// Read back from temporary file
	tmpFile.Seek(0, 0) // Reset to beginning
	buffer := make([]byte, len(content))
	n, err := tmpFile.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading temp file: %v\n", err)
		return
	}
	
	fmt.Printf("Read from temp file: %s\n", string(buffer[:n]))
	
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "example_dir_*")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir) // Clean up
	
	fmt.Printf("Created temporary directory: %s\n", tmpDir)
	
	// Create file in temporary directory
	tmpDirFile := filepath.Join(tmpDir, "temp_file.txt")
	err = os.WriteFile(tmpDirFile, []byte("File in temporary directory"), 0644)
	if err != nil {
		fmt.Printf("Error creating file in temp dir: %v\n", err)
		return
	}
	
	fmt.Printf("Created file in temp directory: %s\n", tmpDirFile)
}

// Function demonstrating file seeking
func fileSeeking() {
	fmt.Println("\n=== FILE SEEKING ===")
	
	filename := "seeking_example.txt"
	content := "0123456789ABCDEFGHIJ"
	
	// Create file with known content
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer os.Remove(filename)
	
	// Open file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Read from beginning
	buffer := make([]byte, 5)
	n, err := file.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return
	}
	fmt.Printf("Read from beginning: %s\n", string(buffer[:n]))
	
	// Seek to position 10
	offset, err := file.Seek(10, 0) // 0 = from beginning
	if err != nil {
		fmt.Printf("Error seeking: %v\n", err)
		return
	}
	fmt.Printf("Seeked to offset: %d\n", offset)
	
	// Read from new position
	n, err = file.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading after seek: %v\n", err)
		return
	}
	fmt.Printf("Read after seek: %s\n", string(buffer[:n]))
	
	// Seek relative to current position
	offset, err = file.Seek(-5, 1) // 1 = from current position
	if err != nil {
		fmt.Printf("Error seeking relative: %v\n", err)
		return
	}
	fmt.Printf("Seeked relative to offset: %d\n", offset)
	
	// Read again
	n, err = file.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading after relative seek: %v\n", err)
		return
	}
	fmt.Printf("Read after relative seek: %s\n", string(buffer[:n]))
	
	// Seek from end
	offset, err = file.Seek(-5, 2) // 2 = from end
	if err != nil {
		fmt.Printf("Error seeking from end: %v\n", err)
		return
	}
	fmt.Printf("Seeked from end to offset: %d\n", offset)
	
	// Read from end position
	n, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading from end: %v\n", err)
		return
	}
	fmt.Printf("Read from end: %s\n", string(buffer[:n]))
}

// Function demonstrating CSV file handling
func csvHandling() {
	fmt.Println("\n=== CSV FILE HANDLING ===")
	
	filename := "data.csv"
	csvContent := `Name,Age,City
Alice,30,New York
Bob,25,London
Carol,35,Tokyo
Dave,28,Sydney`
	
	// Write CSV content
	err := os.WriteFile(filename, []byte(csvContent), 0644)
	if err != nil {
		fmt.Printf("Error writing CSV file: %v\n", err)
		return
	}
	defer os.Remove(filename)
	
	// Read and parse CSV manually
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		return
	}
	defer file.Close()
	
	fmt.Println("Parsing CSV file:")
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		
		if lineNumber == 0 {
			fmt.Printf("Headers: %v\n", fields)
		} else {
			fmt.Printf("Record %d: Name=%s, Age=%s, City=%s\n", 
				lineNumber, fields[0], fields[1], fields[2])
		}
		lineNumber++
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning CSV: %v\n", err)
	}
}

func main() {
	basicFileOperations()
	fileOpenClose()
	bufferedIO()
	scanningFiles()
	filePermissions()
	directoryOperations()
	fileCopying()
	temporaryFiles()
	fileSeeking()
	csvHandling()
	
	fmt.Println("\n=== FILE I/O BEST PRACTICES ===")
	fmt.Println("1. Always close files with defer")
	fmt.Println("2. Handle errors properly")
	fmt.Println("3. Use buffered I/O for better performance")
	fmt.Println("4. Set appropriate file permissions")
	fmt.Println("5. Use temporary files for sensitive operations")
	fmt.Println("6. Clean up temporary files and directories")
	fmt.Println("7. Use filepath package for cross-platform paths")
	fmt.Println("8. Check file existence before operations")
	
	fmt.Println("\nFile I/O examples completed!")
}