package pretty

import (
	"bytes"
	gio "github.com/gogo/protobuf/io"
	"github.com/stretchr/testify/require"
	"testing"
	//"github.com/fatih/color"
	"github.com/cucumber/cucumber-messages-go/v2"
	"github.com/cucumber/gherkin-go"
)

func TestPrintsFeatureHeaderWithComments(t *testing.T) {
	src := `# Hello
Feature: Hello
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src, false), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsFeatureAndScenarioHeadersWithComments(t *testing.T) {
	src := `# Hello
Feature: Hello

  # World
  Scenario: World
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src, false), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsAllTheThings(t *testing.T) {
	src := `# A
Feature: A

  # B
  Background: B
    Given b
      | text | number |
      | a    |     10 |
      | bb   |    100 |
      | ccc  |   1000 |

  # C
  Scenario: C
    Given c
      """
      x
       y
        z
      """
    And <c1>
    Then <c2>

    # CE
    @ce
    Examples: CE
      | c1  |   c2 |
      | a   |   10 |
      | bb  |  100 |
      | ccc | 1000 |

  # D
  Rule: D

    # E
    Background: E
      Given e

    # F
    @f @F
    Scenario: F
      Given f
`

	stdout := &bytes.Buffer{}
	ProcessMessages(messageReader(t, src, false), stdout, false)

	require.EqualValues(t,
		src,
		stdout.String())
}

func TestPrintsInResultsMode(t *testing.T) {
	src := `Feature: A

  Scenario: B
    Given passed
    When failed
    Then skipped
`

	// TODO: Add ANSI codes for cursor up (after printing TestCaseStarted)
	out := `Scenario: B
  ✓ Given passed
  ✗ When failed
    Then skipped
`

	stdout := &bytes.Buffer{}
	prettyStdin := messageReader(t, src, true)
	ProcessMessages(prettyStdin, stdout, true)

	require.EqualValues(t,
		out,
		stdout.String())
}

func messageReader(t *testing.T, src string, fakeResults bool) *bytes.Buffer {
	wrapper := &messages.Wrapper{
		Message: &messages.Wrapper_Source{
			Source: &messages.Source{
				Uri:  "features/test.feature",
				Data: src,
				Media: &messages.Media{
					Encoding:    "UTF-8",
					ContentType: "text/x.cucumber.gherkin+plain",
				},
			},
		},
	}

	wrappers := &bytes.Buffer{}
	messageWriter := gio.NewDelimitedWriter(wrappers)
	messageWriter.WriteMsg(wrapper)

	prettyStdin := &bytes.Buffer{}
	_, err := gherkin.Messages(
		nil,
		wrappers,
		"en",
		true,
		true,
		true,
		prettyStdin,
		false,
		fakeResults,
	)
	require.NoError(t, err)
	//prettyStdinWriter := gio.NewDelimitedWriter(prettyStdin)
	//for _, wrapper := range wrappers {
	//	prettyStdinWriter.WriteMsg(&wrapper)
	//}
	return prettyStdin
}
