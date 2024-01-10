package utils

import "os"

func AppConfigPath() (string, error) {
	dataPath, err := DataPath()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return dataPath + "/config.json", nil
}

func SQLiteDatabasePath() (string, error) {
	dataPath, err := DataPath()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return dataPath + "/wachat.sqlite3", nil
}

func DataPath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return configPath + "/Wachat", nil
}
