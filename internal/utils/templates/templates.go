package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// EmailTemplate represents an email template with its subject and body
type EmailTemplate struct {
	Subject string
	Body    string
}

// GetEmailTemplate returns the requested email template
func GetEmailTemplate(templateName string) (*EmailTemplate, error) {
	templates := map[string]EmailTemplate{
		"order_confirmation": {
			Subject: "Order Confirmation",
			Body: `
<!DOCTYPE html>
<html>
<head>
    <title>Order Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f8f8; padding: 10px; text-align: center; }
        .content { padding: 20px; }
        .footer { margin-top: 20px; font-size: 0.8em; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Thank you for your order!</h1>
        </div>
        <div class="content">
            <p>Hello {{.Order.Customer.FirstName}},</p>
            <p>Your order <strong>#{{.Order.ID}}</strong> has been received.</p>
            
            <h2>Order Summary</h2>
            <p><strong>Total:</strong> {{printf "%.2f" .Order.Total}} {{.Config.Currency}}</p>
            <p><strong>Status:</strong> {{.Order.Status}}</p>
            
            <p>We'll notify you when your order status changes.</p>
        </div>
        <div class="footer">
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`,
		},
		"admin_notification": {
			Subject: "New Order Notification",
			Body: `
<!DOCTYPE html>
<html>
<head>
    <title>New Order #{{.Order.ID}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f8f8; padding: 10px; text-align: center; }
        .content { padding: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Order Received</h1>
        </div>
        <div class="content">
            <h2>Customer Details</h2>
            <p><strong>Name:</strong> {{.Order.Customer.FirstName}} {{.Order.Customer.LastName}}</p>
            <p><strong>Email:</strong> {{.Order.Customer.Email}}</p>
            <p><strong>Phone:</strong> {{.Order.Customer.Phone}}</p>
            
            <h2>Order Details</h2>
            <p><strong>Order ID:</strong> {{.Order.ID}}</p>
            <p><strong>Total:</strong> {{printf "%.2f" .Order.Total}} {{.Config.Currency}}</p>
        </div>
    </div>
</body>
</html>
`,
		},
		"status_update": {
			Subject: "Order Status Update",
			Body: `
<!DOCTYPE html>
<html>
<head>
    <title>Order Status Update</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f8f8; padding: 10px; text-align: center; }
        .content { padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Status Updated</h1>
        </div>
        <div class="content">
            <p>Hello {{.Order.Customer.FirstName}},</p>
            <p>The status of your order <strong>#{{.Order.ID}}</strong> has been updated to:</p>
            <p style="font-size: 1.2em; font-weight: bold; color: #2c3e50;">{{.Order.Status}}</p>
            <p>Thank you for shopping with us!</p>
        </div>
    </div>
</body>
</html>
`,
		},
	}

	if template, ok := templates[templateName]; ok {
		return &template, nil
	}
	return nil, fmt.Errorf("template not found: %s", templateName)
}

// ParseTemplate parses the template body with the given data
func ParseTemplate(templateBody string, data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(templateBody)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}
