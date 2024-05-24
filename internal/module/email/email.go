package email

/*
type email struct {
	logger   slog.Logger
	host     string
	userName string
	password string
	subject  string
}

func Init(logger slog.Logger, host, uName, passwd, sub string) module.Email {
	return &email{
		logger:   logger,
		host:     host,
		userName: uName,
		password: passwd,
		subject:  sub,
	}
}

func (e *email) SendEmail(user db.User, attachmentPath string, wg *sync.WaitGroup) error {
	defer wg.Done()
	// Load email template
	emailTemplate, err := os.ReadFile("./public/email_temp.html")
	if err != nil {
		e.logger.Error("Error reading email template file", err)

		return err
	}

	// Parse email template
	tmpl, err := template.New("email").Parse(string(emailTemplate))
	if err != nil {
		e.logger.Error("Error parsing email template", err)

		return err
	}

	// Prepare email body
	var emailBody bytes.Buffer
	err = tmpl.Execute(&emailBody, struct{ Name string }{Name: user.FirstName})
	if err != nil {
		e.logger.Error("Error executing email template", err)

		return err
	}

	err = e.sendEmail(user.Email, emailBody.String(), attachmentPath)
	if err != nil {
		e.logger.Error("Error sending email", err)

		return err
	}
	if err := os.RemoveAll(attachmentPath); err != nil {
		e.logger.Error("failed to remove attachment file", err)
		return err
	}
	return nil
}

func (e *email) sendEmail(to string, body string, attachmentPath string) error {
	// Initialize SMTP client
	d := gomail.NewDialer(e.host, 587, e.userName, e.password)

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", e.userName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", e.subject)
	m.SetBody("text/html", body)

	// Attach file
	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
*/
