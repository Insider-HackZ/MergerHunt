#!/bin/bash

if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

if ! command -v python3 &> /dev/null; then
    echo "Python3 is not installed. Installing Python3..."
    sudo apt-get update
    sudo apt-get install python3 -y
else
    echo "Python3 is already installed."
fi

if ! command -v pip3 &> /dev/null; then
    echo "pip3 is not installed. Installing pip3..."
    sudo apt-get install python3-pip -y
else
    echo "pip3 is already installed."
fi

pip3 install beautifulsoup4

if ! command -v wget &> /dev/null; then
    echo "wget is not installed. Installing wget..."
    sudo apt-get install wget -y
else
    echo "wget is already installed."
fi

if ! command -v googler &> /dev/null; then
    echo "googler is not installed. Installing googler..."
    sudo apt-get update
    sudo apt-get install googler -y
else
    echo "googler is already installed."
fi
go build test1.go
sudo mv test1 /usr/local/bin
echo "All requirements are installed and ready to go!"
