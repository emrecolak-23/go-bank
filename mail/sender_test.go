package mail

import (
	"testing"

	"github.com/emrecolak-23/go-bank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	sender := NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)

	subject := "A test mail"
	content := `
		<h1>Hello World!</h1>
		<p>This is a test message</p>
	`
	to := []string{"colakkemre@gmail.com"}
	attachFiles := []string{"../Readme.md"}

	err = sender.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		attachFiles,
	)

	require.NoError(t, err)
}
