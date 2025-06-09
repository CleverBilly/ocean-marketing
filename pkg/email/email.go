package email

import (
	"crypto/tls"
	"fmt"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// Client 邮件客户端
type Client struct {
	cfg config.EmailConfig
}

// NewClient 创建邮件客户端
func NewClient(cfg config.EmailConfig) *Client {
	return &Client{cfg: cfg}
}

// SendEmail 发送邮件
func (c *Client) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()

	// 设置发件人
	m.SetHeader("From", c.cfg.From)

	// 设置收件人
	m.SetHeader("To", to...)

	// 设置主题
	m.SetHeader("Subject", subject)

	// 设置邮件内容
	m.SetBody("text/html", body)

	// 创建SMTP拨号器
	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.Username, c.cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		logger.Error("发送邮件失败", zap.Error(err), zap.Strings("to", to))
		return err
	}

	logger.Info("邮件发送成功", zap.Strings("to", to), zap.String("subject", subject))
	return nil
}

// SendHTMLEmail 发送HTML邮件
func (c *Client) SendHTMLEmail(to []string, subject, htmlBody string) error {
	return c.SendEmail(to, subject, htmlBody)
}

// SendPlainTextEmail 发送纯文本邮件
func (c *Client) SendPlainTextEmail(to []string, subject, textBody string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", c.cfg.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBody)

	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.Username, c.cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("发送纯文本邮件失败", zap.Error(err), zap.Strings("to", to))
		return err
	}

	logger.Info("纯文本邮件发送成功", zap.Strings("to", to), zap.String("subject", subject))
	return nil
}

// SendEmailWithAttachment 发送带附件的邮件
func (c *Client) SendEmailWithAttachment(to []string, subject, body string, attachments []string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", c.cfg.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// 添加附件
	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.Username, c.cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Error("发送带附件邮件失败", zap.Error(err), zap.Strings("to", to))
		return err
	}

	logger.Info("带附件邮件发送成功", zap.Strings("to", to), zap.String("subject", subject))
	return nil
}

// SendTemplate 发送模板邮件
type TemplateData struct {
	To       []string
	Subject  string
	Template string
	Data     map[string]interface{}
}

func (c *Client) SendTemplate(data TemplateData) error {
	// 这里可以集成模板引擎，如html/template
	// 简化示例，直接替换占位符
	body := data.Template
	for key, value := range data.Data {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		body = fmt.Sprintf(body, placeholder, fmt.Sprintf("%v", value))
	}

	return c.SendEmail(data.To, data.Subject, body)
}
