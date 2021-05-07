package errors

import (
	"testing"
)

func Test_ErrorsType(t *testing.T) {
	codes := []int{
		EXIT_CODE_CONFIG_ERROR,
		EXIT_CODE_EXECUTION_FAILURE,
		EXIT_CODE_APP_DEV_ERROR,
		EXIT_CODE_BAD_PORT,
		EXIT_CODE_TEMPLATE_FILENAME_EMPTY,
		EXIT_CODE_TPL_NOT_FOUND,
		EXIT_CODE_TPL_ERROR,
		EXIT_CODE_BAD_MAPPING_FILE,
		EXIT_CODE_INVALID_LOGLEVEL,
		EXIT_METRICS_ISSUE,
	}

	for code := range codes {
		if code < 0 {
			t.Errorf("Invalid code, not of type Int")
		}
	}
}