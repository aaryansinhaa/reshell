package scripts

import (
	"reflect"
	"testing"
)

func TestParseParameters(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{
			name:     "no parameters",
			code:     "echo 'hello world'",
			expected: nil,
		},
		{
			name: "positional parameters only",
			code: `echo $1
echo "argument $2 and $9"`,
			expected: []string{"Arg1", "Arg2", "Arg9"},
		},
		{
			name: "declarative param tags only",
			code: `# @param ProjectName
# @param Version
echo "deploying $ProjectName version $Version"`,
			expected: []string{"ProjectName", "Version"},
		},
		{
			name: "mixed positional and declarative tags",
			code: `# @param OutputFile
cp $1 $OutputFile
echo "copied to $OutputFile"`,
			expected: []string{"OutputFile", "Arg1"},
		},
		{
			name: "duplicate parameters are ignored",
			code: `# @param Out
echo $1
echo $1
# @param Out`,
			expected: []string{"Out", "Arg1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseParameters(tt.code)
			if len(got) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseParameters() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
