package ticket

/*
type ticket struct {
	log           slog.Logger
	successUrl    string
	cancelUrl     string
	errorUrl      string
	notifyUrl     string
	itemPrice     float64
	accountNumber string
	bank          string
	amount        float64
	sessionUrl    string
	apiKey        string
	expireDate    time.Duration
	storage       storage.User
}

func Init(
	log slog.Logger, cUrl, eUrl, nUrl, accno, bank, sessUrl, successUrl, key string,
	price, amount float64, storage storage.Ticket, expdate time.Duration) module.Ticket {
	return &ticket{
		log:           log,
		successUrl:    successUrl,
		cancelUrl:     cUrl,
		errorUrl:      eUrl,
		notifyUrl:     nUrl,
		itemPrice:     price,
		accountNumber: accno,
		amount:        amount,
		sessionUrl:    sessUrl,
		apiKey:        key,
		storage:       storage,
		expireDate:    expdate,
		bank:          bank,
	}
}

func (t *ticket) CreateCheckoutSession(c *gin.Context, user model.User) error {

	if strings.HasPrefix(user.Phone, "09") {
		user.Phone = "251" + strings.TrimPrefix(user.Phone, "0")
	}

	nonce := uuid.NewString()
	requestBody := model.PaymentRequest{
		CancelURL:      t.cancelUrl,
		Nonce:          nonce,
		Phone:          user.Phone,
		Email:          user.Email,
		SuccessURL:     fmt.Sprintf("%s/%s", t.successUrl, nonce),
		ErrorURL:       t.errorUrl,
		NotifyURL:      t.notifyUrl,
		PaymentMethods: []string{"TELEBIRR"},
		ExpireDate:     time.Now().Add(t.expireDate * time.Hour).Format("2006-01-02T15:04:05"),
		Items: []interface{}{
			map[string]interface{}{
				"name":        "ticket",
				"quantity":    1,
				"price":       t.itemPrice,
				"description": "Ticket for grand event at Gion Hotel",
				"image":       "",
			},
		},
		Beneficiaries: []model.Beneficiary{
			{
				AccountNumber: t.accountNumber,
				Bank:          t.bank,
				Amount:        t.amount,
			},
		},
		Lang: "EN",
	}

	requestByte, err := json.Marshal(requestBody)
	fmt.Println(string(requestByte))
	if err != nil {
		t.log.Error("failed to marshal the request body of checkout session")
		return err
	}
	request, err := http.NewRequest(http.MethodPost, t.sessionUrl, bytes.NewBuffer(requestByte))
	if err != nil {
		t.log.Error("failed to create request struct for checkout session")
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-arifpay-key", t.apiKey)
	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		t.log.Error("failed to do checkout request", err)
		fmt.Println("err", err)
		return err
	}

	if resp.StatusCode != 200 {
		t.log.Warn("checkout request not successful", resp.StatusCode)
		return fmt.Errorf("checkout session request failed ")
	}

	var checkout model.CheckoutResponse
	err = json.NewDecoder(resp.Body).Decode(&checkout)
	if err != nil {
		t.log.Error("failed to decode checkout response body", err)
		return err
	}
	err = t.storage.RegisterUserToDb(user, checkout.Data.SessionId, nonce)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, checkout.Data.PaymentUrl)
	return nil

}
func (t *ticket) UpdatePaymentStatus(status, sid string) (db.User, error) {
	user, err := t.storage.UpdatePaymentStatus(status, sid)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// GeneratePDFTicket generates a PDF ticket with user information and a QR code
func (t *ticket) GeneratePDFTicket(userData db.User) (*gopdf.GoPdf, error) {
	// Create PDF
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 396, H: 150}})
	pdf.AddPage()
	pdf.SetLineWidth(1)
	pdf.SetFillColor(0, 0, 0)

	// Set font
	if err := pdf.AddTTFFont("Arial", "./public/font/arial.ttf"); err != nil {
		return nil, err
	}
	pdf.SetFont("Arial", "", 8)
	// Draw QR code
	qrFileName := fmt.Sprintf("./public/image/qr_%d.png", userData.ID) // Use unique ID as file name
	qrURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=%d", userData.ID)
	qrImage, err := http.Get(qrURL)
	if err != nil {
		return nil, err
	}
	defer qrImage.Body.Close()

	// Create a file to store the QR code image
	qrFile, err := os.Create(qrFileName)
	if err != nil {
		return nil, err
	}
	defer qrFile.Close()

	// Write the QR code image to the file
	_, err = io.Copy(qrFile, qrImage.Body)
	if err != nil {
		return nil, err
	}
	err = pdf.Image("./public/image/ticket.png", 0, 0, &gopdf.Rect{
		H: 150,
		W: 390,
	})
	if err != nil {
		t.log.Debug("failed to read png")
		return nil, err
	}
	// Embed the QR code image into the PDF
	err = pdf.Image(qrFileName, 310, 0, &gopdf.Rect{
		H: 90,
		W: 80,
	})
	pdf.SetX(310)
	pdf.SetY(125)
	pdf.Cell(&gopdf.Rect{W: 5, H: 5}, fmt.Sprintf("Name: %s %s", userData.FirstName, userData.LastName))

	pdf.Br(8)
	pdf.SetX(310)
	pdf.Cell(&gopdf.Rect{W: 5, H: 5}, fmt.Sprintf("Phone: %s", userData.Phone))

	pdf.Br(8)
	pdf.SetX(310)
	pdf.Cell(&gopdf.Rect{W: 5, H: 5}, fmt.Sprintf("TicketNo: %d", userData.ID))
	if err != nil {
		return nil, err
	}
	if err := os.RemoveAll(qrFileName); err != nil {
		return nil, err
	}
	return &pdf, nil
}


*/
