package utils

import (
	"fmt"
	"os"

	"../models"
)

func DeleteFiles(folderPath string, res *models.ResponseObj) {
	err := os.RemoveAll(folderPath)

	if err != nil {
		fmt.Printf("Could not remove folder %s. Error: %s\n", folderPath, err.Error())
		res.ConsoleOut += "\nCould not remove your files from server, please notify. Error: " + err.Error()
	}
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
