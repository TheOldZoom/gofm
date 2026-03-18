package verbose

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Enabled() bool {
	return viper.GetBool("verbose")
}

func Printf(format string, args ...any) {
	if !Enabled() {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func PrintJSON(label string, v any) {
	if !Enabled() {
		return
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s (json marshal failed: %v)\n", label, err)
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n%s\n", label, data)
}

func PrintJSONBytes(label string, data []byte) {
	if !Enabled() {
		return
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n%s\n", label, data)
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n%s\n", label, pretty.Bytes())
}
