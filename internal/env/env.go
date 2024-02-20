package env

import "os"

type RequiredEnvVars struct {
	UploadFilePath      string
	UploadFileRecipient string
}

func GetRequiredEnvVars() (*RequiredEnvVars, error) {
	var (
		vars = new(RequiredEnvVars)
		err  error
	)

	if vars.UploadFilePath, err = getEnvVar("UPLOAD_FILE_PATH"); err != nil {
		return nil, err
	}

	if vars.UploadFileRecipient, err = getEnvVar("UPLOAD_FILE_RECIPIENT"); err != nil {
		return nil, err
	}

	return vars, nil
}

func getEnvVar(envVar string) (string, error) {
	value := os.Getenv(envVar)
	if value == "" {
		return "", NewRequiredEnvVarError(envVar)
	}

	return value, nil
}
