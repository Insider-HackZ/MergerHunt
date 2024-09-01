package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"strconv"
)

func main() {
	banner()
	if len(os.Args) < 2 {
		fmt.Println("Please provide an organization name using '--org'")
		return
	}

	if os.Args[1] == "--org" && len(os.Args) > 2 {
		searchTerm := os.Args[2]
		command(searchTerm)
	} else {
		fmt.Println("Usage: go run tool.go --org <organization_name>")
		return
	}

	err := runPythonScript()
	if err != nil {
		fmt.Println("Error running Python script:", err)
	} else {
		fmt.Println("Script executed successfully. Check f_output.txt for results.")
	}
	Rm_extra()
}

func command(searchTerm string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("googler -n 5 %s acquisition wikipedia --json > test.txt", searchTerm))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running googler:", err)
		return
	}

	fmt.Println("Working on it....")

	cmd2 := exec.Command("bash", "-c", `cat test.txt | grep -oP '"url": *"\K[^"]+' | tee ot.txt`)
	errorChecker(cmd2.Run())

	cmd3 := exec.Command("bash", "-c", `rm test.txt`)
	errorChecker(cmd3.Run())

	urls, e := readURLsFromFile("ot.txt")
	errorChecker(e)

	urlKeywords := []string{"list", "mergers", "acquisitions"}
	keywordGroups := [][]string{
		{"Number"},
		{"Acquisition date", "Acquisition data", "Date"},
		{"Company"},
		{"Business"},
		{"Country"},
		{"Price", "Value (USD)", "Acquired for (USD)"},
		{"Used as or integrated with", "Derived products"},
		{"Refs.", "References"},
		{"Acquired on"},
		{"Acquisition type", "Acquisition status"},
		{"Deal size"},
		{"Transaction type"},
		{"Description"},
	}

	for i, url := range urls {
		if containsAny(url, urlKeywords) && checkForHalfGroups(url, keywordGroups) {
			fmt.Printf("Downloading: %s\n", url)
			fileName := fmt.Sprintf("html_data-tool%d.txt", i+1)
			cmd4 := exec.Command("bash", "-c", fmt.Sprintf("wget -O %s %s", fileName, url))
			errorChecker(cmd4.Run())

			cmd5 := exec.Command("bash", "-c", fmt.Sprintf("paste -s -d ' ' %s >> output.txt", fileName))
			errorChecker(cmd5.Run())

			cmd6 := exec.Command("bash", "-c", fmt.Sprintf("rm %s", fileName))
			errorChecker(cmd6.Run())
		}
	}
	cmd7 := exec.Command("bash", "-c", "grep -oP '<td>.*?</td>' output.txt > all_td_tags.txt")
	errorChecker(cmd7.Run())
}

func readURLsFromFile(filePath string) ([]string, error) {
	var urls []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func checkForHalfGroups(url string, keywordGroups [][]string) bool {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading content from URL %s: %v\n", url, err)
		return false
	}

	contentStr := strings.ToLower(string(content))

	totalGroups := len(keywordGroups)
	foundCount := 0

	for _, group := range keywordGroups {
		if containsAny(contentStr, group) {
			foundCount++
		}
	}

	return foundCount >= (totalGroups+1)/2
}

func containsAny(contentStr string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(contentStr, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func isNumeric(str string) bool {
	str = strings.TrimSpace(str)
	_, err := strconv.Atoi(str)
	return err == nil
}

func errorChecker(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func runPythonScript() error {
	pythonCode := `
from bs4 import BeautifulSoup

# Read HTML content from a file
with open('all_td_tags.txt', 'r') as file:
    html = file.read()

soup = BeautifulSoup(html, 'html.parser')
rows = soup.find_all('td')

# Open an output file
with open('f_output.txt', 'w') as output_file:
    # Iterate through all <td> tags
    for i in range(len(rows)):
        # Check if the <td> contains a number
        if rows[i].get_text(strip=True).isdigit():
            number = rows[i].get_text(strip=True)
            # Ensure there are at least two more <td> elements following the current one
            if i + 2 < len(rows):
                date = rows[i + 1].get_text(strip=True)
                company = rows[i + 2].get_text(strip=True)
                # Write each entry to the output file, each on a new line
                output_file.write(f"Number: {number}, Date: {date}, Company: {company}\n")

	`

	tempFile, err := os.CreateTemp("", "*.py")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the Python code to the temporary file
	_, err = tempFile.WriteString(pythonCode)
	if err != nil {
		return fmt.Errorf("error writing to temp file: %w", err)
	}

	// Execute the Python code
	cmd := exec.Command("python3", tempFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing Python code: %w\nOutput: %s", err, output)
	}


	return nil
}

func Rm_extra(){


	
	cmd8 := exec.Command("bash", "-c", "rm all_td_tags.txt output.txt ot.txt ")
	errorChecker(cmd8.Run())
}

func banner() {
	fmt.Printf(`
    _                   _   ____                  _ 
   / \   ___ __ _ _   _(_) / ___|  ___ ___  _   _| |_
  / _ \ / __/ _' | | | | | \___ \ / __/ _ \| | | | __|
 / ___ \ (_| (_| | |_| | |  ___) | (_| (_) | |_| | |_ 
/_/   \_\___\__, |\__,_|_| |____/ \___\___/ \__,_|\__|
               |_| developed by: harshj054
`)
}
