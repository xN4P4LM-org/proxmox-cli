package main

import "os"

func checkParams(params []string, expectedCommand string, value_required bool) bool {
	shortCommand := "-" + string(expectedCommand[0])
	longCommand := "--" + expectedCommand

	for paramLocation, param := range params {
		if param == shortCommand || param == longCommand {
			if value_required {
				// if the next param is a command, or the last param, return false
				if paramLocation+1 < len(params) {
					if checkValue(params[paramLocation+1]) {
						return true
					}
				}
			} else {
				return true
			}
		}
	}

	return false
}

func checkValue(params string) bool {
	return string(params[0]) != "-"
}

func getParams(params []string, expectedCommand string) string {
	shortCommand := "-" + string(expectedCommand[0])
	longCommand := "--" + expectedCommand

	for paramLocation, param := range params {
		if param == shortCommand || param == longCommand {
			if paramLocation+1 < len(params) {
				return params[paramLocation+1]
			}
		}
	}

	return ""
}

func getEnvVariable(envVar string) string {
	variable, variable_ok := os.LookupEnv(envVar)

	if !variable_ok {
		logger.Fatal("Environment variable ", envVar, " not set")
	}

	return variable
}
