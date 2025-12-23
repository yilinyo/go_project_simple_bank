package mail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yilinyo/project_bank/util"
)

func TestSendGMail(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	fmt.Printf("%+v\n", config)
	emailSender := NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)

	subject := "TestSendGMail"

	content := `<h1>This is a test email</h1>`

	to := []string{"1322780122@qq.com"}
	attachFiles := []string{"../README.md"}

	err = emailSender.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		attachFiles,
	)
	require.NoError(t, err)
}
