# MergerHunt
This tool fetches and processes acquisition data related to a specified organization. It uses Go for the main logic, shell commands for data extraction, and Python for HTML parsing.

## Features

- Fetches top 5 Google search results for an organization related to acquisitions.
- Filters and downloads relevant web pages.
- Extracts table data from the downloaded web pages.
- Processes the data to retrieve useful information about acquisitions.
- Outputs the results to a file named `f_output.txt`.

## Requirements

- **Go**: Ensure Go is installed on your system.
- **Python 3**: Required for the HTML parsing script.
- **pip3**: Required for installing Python packages.
- **BeautifulSoup4**: Python package for parsing HTML.
- **wget**: Required to download web pages.
- **googler**: A command-line tool to search Google from the terminal.

## Setup Instructions

### Step 1: Download and Install Dependencies

1. Clone or download this repository to your local machine.
 ```
 git clone https://github.com/Byte-BloggerBase/MergerHunt.git
 ```

3. Run the `setup.sh` script to install all necessary dependencies:

 ```
   ./setup.sh
 ```

This will check for and install the following:

- Python3
- pip3
- BeautifulSoup4
- wget
- googler

### Step 2: Run the Tool

To use the tool, execute the following command in your terminal:

```
MergerHunt --org <organization_name>
```

Replace `<organization_name>` with the name of the organization you want to search for.

### Step 3: View the Results

After running the tool, you can check the results in the `f_output.txt` file. This file contains the processed acquisition data related to the specified organization.

> If anyone would like to contribute to the development of Insider-HackZ/MergerHunt, please send an email to official@bytebloggerbase.com.
