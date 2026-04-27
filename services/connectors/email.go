package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	imap "github.com/emersion/go-imap/v2"
	sasl "github.com/emersion/go-sasl"
	smtp "github.com/emersion/go-smtp"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/charset"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/mudler/LocalAGI/core/agent"
	"github.com/mudler/LocalAGI/core/types"
	"github.com/mudler/LocalAGI/pkg/config"
	"github.com/mudler/xlog"
	"github.com/sashabaranov/go-openai"
)

type Email struct {
	username        string
	name            string
	password        string
	email           string
	smtpServer      string
	smtpInsecure    bool
	imapServer      string
	imapInsecure    bool
	defaultEmail    string
	filterConfigDir string
}

func NewEmail(config map[string]string) *Email {
	filterDir := config["filterConfigDir"]
	if filterDir == "" {
		filterDir = "/pool/agents"
	}

	return &Email{
		username:        config["username"],
		name:            config["name"],
		password:        config["password"],
		email:           config["email"],
		smtpServer:      config["smtpServer"],
		smtpInsecure:    config["smtpInsecure"] == "true",
		imapServer:      config["imapServer"],
		imapInsecure:    config["imapInsecure"] == "true",
		defaultEmail:    config["defaultEmail"],
		filterConfigDir: filterDir,
	}
}

func EmailConfigMeta() []config.Field {
	return []config.Field{
		{
			Name:     "smtpServer",
			Label:    "SMTP Host:port",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP server host:port (e.g., smtp.gmail.com:587)",
		},
		{
			Name:  "smtpInsecure",
			Label: "Insecure SMTP",
			Type:  config.FieldTypeCheckbox,
		},
		{
			Name:     "imapServer",
			Label:    "IMAP Host:port",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "IMAP server host:port (e.g., imap.gmail.com:993)",
		},
		{
			Name:  "imapInsecure",
			Label: "Insecure IMAP",
			Type:  config.FieldTypeCheckbox,
		},
		{
			Name:     "username",
			Label:    "Username",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "Username/email address",
		},
		{
			Name:     "name",
			Label:    "Friendly Name",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "Friendly name of sender",
		},
		{
			Name:     "password",
			Label:    "Password",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "SMTP/IMAP password or app password",
		},
		{
			Name:     "email",
			Label:    "From Email",
			Type:     config.FieldTypeText,
			Required: true,
			HelpText: "Agent email address",
		},
		{
			Name:     "defaultEmail",
			Label:    "Default Recipient",
			Type:     config.FieldTypeText,
			HelpText: "Default email address to send messages to when the agent wants to initiate a conversation",
		},
	}
}

func (e *Email) AgentResultCallback() func(state types.ActionState) {
	return func(state types.ActionState) {
		// Send the result to the bot
	}
}

func (e *Email) AgentReasoningCallback() func(state types.ActionCurrentState) bool {
	return func(state types.ActionCurrentState) bool {
		// Send the reasoning to the bot
		return true
	}
}

func filterEmailRecipients(input string, emailToRemove string) string {

	addresses := strings.Split(strings.TrimPrefix(input, "To: "), ",")

	var filtered []string
	for _, address := range addresses {
		address = strings.TrimSpace(address)
		if !strings.Contains(address, emailToRemove) {
			filtered = append(filtered, address)
		}
	}

	if len(filtered) > 0 {
		return strings.Join(filtered, ", ")
	}
	return ""
}

// extractTextContent recursively extracts text/plain or text/html from multipart messages
func extractTextContent(msg *message.Entity) string {
	mediaType, _, err := msg.Header.ContentType()
	if err != nil {
		// Not a valid Content-Type, read body directly
		buf := new(bytes.Buffer)
		buf.ReadFrom(msg.Body)
		return buf.String()
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		// Not multipart, read body directly
		buf := new(bytes.Buffer)
		buf.ReadFrom(msg.Body)
		return buf.String()
	}

	// Multipart message - iterate through parts
	mr := msg.MultipartReader()
	if mr == nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(msg.Body)
		return buf.String()
	}

	var result string
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		partMediaType, _, _ := part.Header.ContentType()
		switch {
		case strings.HasPrefix(partMediaType, "multipart/"):
			// Recurse into nested multipart
			result = extractTextContent(part)
			if result != "" {
				return result
			}
		case partMediaType == "text/plain":
			buf := new(bytes.Buffer)
			buf.ReadFrom(part.Body)
			result = buf.String()
			if result != "" {
				return result
			}
		case partMediaType == "text/html":
			buf := new(bytes.Buffer)
			buf.ReadFrom(part.Body)
			result = buf.String()
			// Keep going in case text/plain exists
		}
	}
	return result
}

func (e *Email) sendMail(to, subject, content, replyToID, references string, emails []string, html bool) {
	xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [SMTP] sendMail called: to=%s from=%s (%s) subject=%s content_len=%d", to, e.email, e.username, subject, len(content)))

	auth := sasl.NewPlainClient("", e.username, e.password)

	contentType := "text/plain"
	if html {
		contentType = "text/html"
	}

	var replyHeaders string
	if replyToID != "" {
		referenceLine := strings.ReplaceAll(references+" "+replyToID, "\n", "")
		replyHeaders = fmt.Sprintf("In-Reply-To: %s\r\nReferences: %s\r\n", replyToID, referenceLine)
	}

	// Build full message content
	var builder strings.Builder
	fmt.Fprintf(&builder, "To: %s\r\n", to)
	fmt.Fprintf(&builder, "From: %s <%s>\r\n", e.name, e.email)
	builder.WriteString(replyHeaders)
	fmt.Fprintf(&builder, "MIME-Version: 1.0\r\nContent-Type: %s;\r\n", contentType)
	fmt.Fprintf(&builder, "Subject: %s\r\n\r\n", subject)
	fmt.Fprintf(&builder, "%s\r\n", content)
	msg := strings.NewReader(builder.String())

	if !e.smtpInsecure {

		err := smtp.SendMail(e.smtpServer, auth, e.username, emails, msg)
		if err != nil {
			xlog.Error(fmt.Sprintf("Email send err: %v", err))
		} else {
			xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [SMTP] Email sent successfully to %v", emails))
		}

	} else {

		c, err := smtp.Dial(e.smtpServer)
		if err != nil {
			xlog.Error(fmt.Sprintf("Email connection err: %v", err))
		}
		defer c.Close()

		err = c.Hello("client")
		if err != nil {
			xlog.Error(fmt.Sprintf("Email hello err: %v", err))
		}

		err = c.Auth(auth)
		if err != nil {
			xlog.Error(fmt.Sprintf("Email auth err: %v", err))
		}

		err = c.SendMail(e.username, emails, msg)
		if err != nil {
			xlog.Error(fmt.Sprintf("Email send err: %v", err))
				xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [SMTP] Email sent successfully to %v", emails))
		}

	}
}

// classifyEmail sends the email content to the filter LLM and returns classification + reasoning.
// On any error, defaults to "benign" (allow through) to avoid blocking legitimate emails.
func (e *Email) classifyEmail(subject, from, content string) (classification, reasoning, instruccion string) {
	classification = "benign"
	instruccion = "proceed"

	apiURL := os.Getenv("OPENAI_BASE_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("MODEL_NAME")
	if apiURL == "" {
		apiURL = "https://openrouter.ai/api/v1"
	}

	// Read filter system prompt from agent config file
	filterPrompt := e.loadFilterPrompt()
	if filterPrompt == "" {
		xlog.Warn("Filter prompt empty, skipping classification")
		return
	}

	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = apiURL
	client := openai.NewClientWithConfig(clientConfig)

	userMsg := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", from, subject, content)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: filterPrompt},
				{Role: "user", Content: userMsg},
			},
		},
	)
	if err != nil {
		xlog.Warn(fmt.Sprintf("Filter LLM call failed, allowing through: %v", err))
		return
	}

	if len(resp.Choices) == 0 {
		xlog.Warn("Filter LLM returned no choices, allowing through")
		return
	}

	raw := resp.Choices[0].Message.Content
	xlog.Debug(fmt.Sprintf("Filter raw response: %s", raw))

	// Try to extract JSON from the response (handle markdown-wrapped JSON)
	jsonStr := raw
	if idx := strings.Index(raw, "{"); idx >= 0 {
		if end := strings.LastIndex(raw, "}"); end >= idx {
			jsonStr = raw[idx : end+1]
		}
	}

	var result struct {
		Classification string `json:"classification"`
		Explicacion    string `json:"explicacion"`
		Instruccion    string `json:"instruccion"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		xlog.Warn(fmt.Sprintf("Filter JSON parse failed, allowing through: %v", err))
		return
	}

	if result.Classification == "" {
		result.Classification = "benign"
	}
	if result.Instruccion == "" {
		result.Instruccion = "proceed"
	}

	xlog.Info(fmt.Sprintf("Filter result: classification=%s instruccion=%s razon=%s", result.Classification, result.Instruccion, result.Explicacion))
	return result.Classification, result.Explicacion, result.Instruccion
}

// loadFilterPrompt reads the filter agent's system prompt from its JSON config file.
// Fallback to empty string if not found.
func (e *Email) loadFilterPrompt() string {
	path := filepath.Join(e.filterConfigDir, "agente-filtro-intencion.json")
	data, err := os.ReadFile(path)
	if err != nil {
		xlog.Warn(fmt.Sprintf("Filter config not found at %s: %v", path, err))
		return ""
	}
	data = []byte(os.ExpandEnv(string(data)))
	var cfg struct {
		SystemPrompt string `json:"system_prompt"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		xlog.Warn(fmt.Sprintf("Filter config parse error: %v", err))
		return ""
	}
	return cfg.SystemPrompt
}

const imapFetchTimeout = 30 * time.Second

// fetchMessageWithTimeout fetches a single IMAP message with a timeout.
// If the IMAP connection hangs, the timeout triggers a close to unblock.
func fetchMessageWithTimeout(c *imapclient.Client, seqNum uint32) (*imapclient.FetchMessageBuffer, error) {
	seqSet := imap.SeqSetNum(seqNum)
	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		Flags:       true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{bodySection},
	}

	type result struct {
		buf []*imapclient.FetchMessageBuffer
		err error
	}

	done := make(chan result, 1)
	go func() {
		bufs, err := c.Fetch(seqSet, fetchOptions).Collect()
		done <- result{bufs, err}
	}()

	select {
	case r := <-done:
		if r.err != nil {
			return nil, r.err
		}
		if len(r.buf) == 0 {
			return nil, fmt.Errorf("no message data returned")
		}
		return r.buf[0], nil
	case <-time.After(imapFetchTimeout):
		c.Close()
		return nil, fmt.Errorf("fetch timed out after %v", imapFetchTimeout)
	}
}

func (e *Email) processEmail(a *agent.Agent, fmb *imapclient.FetchMessageBuffer) {
	// Build a minimal bodySection for FindBodySection call
	bodySection := &imap.FetchItemBodySection{}

	// Download Email contents
	r := bytes.NewReader(fmb.FindBodySection(bodySection))
	msg, err := message.Read(r)
	if err != nil {
		xlog.Error(fmt.Sprintf("Email reader err: %v", err))
		return
	}
	content := extractTextContent(msg)

	xlog.Debug("New email!")
	xlog.Debug(fmt.Sprintf("From: %s", msg.Header.Get("From")))
	xlog.Debug(fmt.Sprintf("To: %s", msg.Header.Get("To")))
	xlog.Debug(fmt.Sprintf("Subject: %s", msg.Header.Get("Subject")))
		xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [PASO1] Email RAW headers: From=%s To=%s Subject=%s Content-Type=%s", msg.Header.Get("From"), msg.Header.Get("To"), msg.Header.Get("Subject"), msg.Header.Get("Content-Type")))

	// Skip emails sent by the monitored account itself (prevent reply loops)
	if strings.Contains(msg.Header.Get("From"), e.username) {
		xlog.Debug("Email from self, skipping to prevent loop")
		return
	}

		// In the event that an email account has multiple email addresses, only respond to the one configured.
		// Check To header (direct) and Delivered-To (forwarding). Also match IMAP username for direct sends.
		allowedTo := strings.Contains(msg.Header.Get("To"), e.email) ||
			strings.Contains(msg.Header.Get("To"), e.username) ||
			strings.Contains(msg.Header.Get("Delivered-To"), e.email) ||
			strings.Contains(msg.Header.Get("Delivered-To"), e.username)
	if !allowedTo {
		xlog.Info(fmt.Sprintf("Email was sent to %s, but appeared in my inbox (%s). Ignoring!", msg.Header.Get("To"), e.email))
		return
	}

	// Only respond to emails from authorized senders
	fromHeader := msg.Header.Get("From")
	allowedSenders := strings.Contains(fromHeader, "@ionet.cl") ||
		strings.Contains(fromHeader, "@iodevs.net") ||
		strings.Contains(fromHeader, "el.agente.ion@gmail.com") ||
		strings.Contains(fromHeader, "ventas.ionet@gmail.com")
	if !allowedSenders {
		xlog.Info(fmt.Sprintf("Email from %s is not an authorized sender. Ignoring!", fromHeader))
		return
	}

	contentIsHTML := false

	// Convert email to markdown only if it's in HTML
	prefixes := []string{"<html", "<body", "<div", "<head"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(content), prefix) {
			converted, err := htmltomarkdown.ConvertString(content)
			contentIsHTML = true
			if err != nil {
				xlog.Error(fmt.Sprintf("Email html => md err: %v", err))
				contentIsHTML = false
			} else {
				content = converted
			}
		}
	}

	xlog.Debug(fmt.Sprintf("Markdown:\n\n%s", content))
		xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [PASO2] Content after markdown: %d chars, starts with: %.200s", len(content), content))

	// Construct prompt
	prompt := fmt.Sprintf("%s %s:\n\nFrom: %s\nTime: %s\nSubject: %s\nToday is: %s.\n=====\n%s",
		"This email thread was sent to you. You are",
		e.email,
		msg.Header.Get("From"),
		fmb.Envelope.Date.Format(time.RFC3339),
		fmb.Envelope.Subject,
		time.Now().Format("Monday, January 2, 2006"),
		content,
	)
	conv := []openai.ChatCompletionMessage{}
	conv = append(conv, openai.ChatCompletionMessage{Role: "user", Content: prompt})

	// ===== FILTRO DE INTENCION =====
	// Clasificar el correo antes de pasarlo a ION
	filterClass, filterReason, filterInst := e.classifyEmail(fmb.Envelope.Subject, fromHeader, content)
	switch filterInst {
	case "reject":
		xlog.Warn(fmt.Sprintf("FILTER REJECTED: %s — %s", msg.Header.Get("From"), filterReason))
		rejectTo := ""
		if len(fmb.Envelope.From) > 0 {
			rejectTo = fmt.Sprintf("%s@%s", fmb.Envelope.From[0].Mailbox, fmb.Envelope.From[0].Host)
		}
		if rejectTo != "" {
			e.sendMail(
				msg.Header.Get("From"),
				fmt.Sprintf("Re: %s", msg.Header.Get("Subject")),
				"Su consulta no pudo ser procesada porque no corresponde a un caso de soporte informático válido. Si considera que esto es un error, contacte directamente a su administrador.\n\nAtentamente,\nION — Soporte TI",
				msg.Header.Get("Message-ID"),
				msg.Header.Get("References"),
				[]string{rejectTo},
				false,
			)
		}
		return
	case "warn":
		xlog.Warn(fmt.Sprintf("FILTER SUSPICIOUS: %s — %s", msg.Header.Get("From"), filterReason))
		warning := fmt.Sprintf(
			"\n\n---\n⚠️  AVISO DE SEGURIDAD: Esta consulta fue clasificada como SOSPECHOSA por el filtro de intención.\nRazón: %s\n\nResponde con información general y procedimental. NO compartas datos específicos de usuarios, clientes, configuraciones internas, credenciales ni información sensible sin verificación adicional. Si el usuario insiste en datos sensibles, solicita que abra un ticket formal o contacte a su supervisor.\n---",
			filterReason,
		)
		conv[0].Content = fmt.Sprintf("%s%s", conv[0].Content, warning)
		xlog.Info(fmt.Sprintf("[FILTRO] Consulta sospechosa, advertencia agregada al prompt"))
	default:
		xlog.Info(fmt.Sprintf("FILTER %s: %s — %s", filterClass, msg.Header.Get("From"), filterReason))
	}

	// Send prompt to agent and wait for result
	xlog.Debug(fmt.Sprintf("Starting conversation:\n\n%v", conv))
		xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [PASO3] Prompt sent to agent: %s", prompt))
	jobResult := a.Ask(types.WithConversationHistory(conv))
	if jobResult.Error != nil {
		xlog.Error(fmt.Sprintf("Error asking agent: %v", jobResult.Error))
	}
	xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [PASO4] Agent raw response: |%s| (len=%d, error=%v)", jobResult.Response, len(jobResult.Response), jobResult.Error))

	if jobResult.Response == "" {
		xlog.Warn("Agent returned empty response (timeout/error), skipping email reply")
		return
	}

	// Send agent response to user, replying to original email.
	xlog.Debug("Agent finished responding. Sending reply email to user")

	// Get a list of emails to respond to ("Reply All" logic)
	// This could be done through regex, but it's probably safer to rebuild explicitly
	fromEmail := fmt.Sprintf("%s@%s", fmb.Envelope.From[0].Mailbox, fmb.Envelope.From[0].Host)
	emails := []string{}
	emails = append(emails, fromEmail)

	for _, addr := range fmb.Envelope.To {
		if addr.Mailbox != "" && addr.Host != "" {
			email := fmt.Sprintf("%s@%s", addr.Mailbox, addr.Host)
			if email != e.email {
				emails = append(emails, email)
			}
		}
	}

	// Keep the original header, in case sender had contact names as part of the header
	newToHeader := msg.Header.Get("From") + ", " + filterEmailRecipients(msg.Header.Get("To"), e.email)

	// Create the body of the email
	replyContent := jobResult.Response
	if jobResult.Response == "" {
		replyContent =
			"System: I'm sorry, but it looks like the agent did not respond. " +
				"This could be in error, or maybe it had nothing to say."
	}

	// Quote the original message. This lets the agent see conversation history and is an email standard.
	quoteHeader := fmt.Sprintf("\r\n\r\nOn %s, %s wrote:\n",
		fmb.Envelope.Date.Format("Monday, Jan 2, 2006 at 15:04"),
		fmt.Sprintf("%s <%s>", fmb.Envelope.From[0].Name, fromEmail),
	)
	quotedLines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	for i, line := range quotedLines {
		quotedLines[i] = "> " + line
	}
	replyContent = replyContent + quoteHeader + strings.Join(quotedLines, "\r\n")

	// Convert agent markdown response to HTML for clean email rendering
	{
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
		doc := p.Parse([]byte(replyContent))

		opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}
		renderer := html.NewRenderer(opts)

		replyContent = string(markdown.Render(doc, renderer))
	}
	contentIsHTML = true

	// Send the email
		xlog.Info(fmt.Sprintf("[ULTRA-DEBUG] [PASO5] Sending reply via SMTP: to=%s subject=Re: %s reply_len=%d", newToHeader, msg.Header.Get("Subject"), len(replyContent)))
	e.sendMail(newToHeader,
		fmt.Sprintf("Re: %s", msg.Header.Get("Subject")),
		replyContent,
		msg.Header.Get("Message-ID"),
		msg.Header.Get("References"),
		emails,
		contentIsHTML,
	)
}

// runIMAPWorker polls INBOX for new messages with timeout protection.
// Returns error to trigger reconnection in the outer loop.
func (e *Email) runIMAPWorker(ctx context.Context, c *imapclient.Client, a *agent.Agent, startIndex uint32) error {
	currentIndex := startIndex

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		selectedMbox, err := c.Select("INBOX", nil).Wait()
		if err != nil {
			return fmt.Errorf("select INBOX: %w", err)
		}

		// Process new messages
		for currentIndex < selectedMbox.NumMessages {
			currentIndex++

			xlog.Debug(fmt.Sprintf("Fetching message %d", currentIndex))
			fmb, err := fetchMessageWithTimeout(c, currentIndex)
			if err != nil {
				return fmt.Errorf("fetch message %d: %w", currentIndex, err)
			}

			// Process email in background goroutine so polling is not blocked
			go e.processEmail(a, fmb)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func (e *Email) Start(a *agent.Agent) {
	go func() {
		if e.defaultEmail != "" {
			// handle new conversations
			a.AddSubscriber(func(ccm *types.ConversationMessage) {
				xlog.Debug("Subscriber(email)", "message", ccm.Message.Content)

				// Send the message to the default email
				e.sendMail(
					e.defaultEmail,
					"Message from LocalAGI",
					ccm.Message.Content,
					"",
					"",
					[]string{e.defaultEmail},
					false,
				)

				a.SharedState().ConversationTracker.AddMessage(
					fmt.Sprintf("email:%s", e.defaultEmail),
					openai.ChatCompletionMessage{
						Content: ccm.Message.Content,
						Role:    "assistant",
					},
				)
			})
		}

		xlog.Info("Email connector is now running.  Press CTRL-C to exit.")

		ctx := a.Context()
		backoff := 1 * time.Second
		maxBackoff := 30 * time.Second

		for {
			select {
			case <-ctx.Done():
				xlog.Info("Email connector is now stopped.")
				return
			default:
			}

			// IMAP dial
			imapOpts := &imapclient.Options{WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader}}
			var c *imapclient.Client
			var err error
			if e.imapInsecure {
				c, err = imapclient.DialInsecure(e.imapServer, imapOpts)
			} else {
				c, err = imapclient.DialTLS(e.imapServer, imapOpts)
			}

			if err != nil {
				xlog.Error(fmt.Sprintf("Email IMAP dial err: %v (retry in %v)", err, backoff))
				select {
				case <-ctx.Done():
					xlog.Info("Email connector is now stopped.")
					return
				case <-time.After(backoff):
				}
				backoff = min(backoff*2, maxBackoff)
				continue
			}

			// IMAP login
			err = c.Login(e.username, e.password).Wait()
			if err != nil {
				xlog.Error(fmt.Sprintf("Email IMAP login err: %v", err))
				c.Close()
				select {
				case <-ctx.Done():
					xlog.Info("Email connector is now stopped.")
					return
				case <-time.After(backoff):
				}
				backoff = min(backoff*2, maxBackoff)
				continue
			}

			// IMAP mailbox
			mailboxes, err := c.List("", "%", nil).Collect()
			if err != nil {
				xlog.Error(fmt.Sprintf("Email IMAP mailbox err: %v", err))
				c.Close()
				select {
				case <-ctx.Done():
					xlog.Info("Email connector is now stopped.")
					return
				case <-time.After(backoff):
				}
				backoff = min(backoff*2, maxBackoff)
				continue
			}

			xlog.Debug(fmt.Sprintf("Email IMAP mailbox count: %v", len(mailboxes)))
			for _, mbox := range mailboxes {
				xlog.Debug(fmt.Sprintf(" - %v", mbox.Mailbox))
			}

			// Select INBOX
			selectedMbox, err := c.Select("INBOX", nil).Wait()
			if err != nil {
				xlog.Error(fmt.Sprintf("Cannot select INBOX mailbox! %v", err))
				c.Close()
				select {
				case <-ctx.Done():
					xlog.Info("Email connector is now stopped.")
					return
				case <-time.After(backoff):
				}
				backoff = min(backoff*2, maxBackoff)
				continue
			}
			xlog.Debug(fmt.Sprintf("INBOX contains %v messages", selectedMbox.NumMessages))

			// Reset backoff on successful connection
			backoff = 1 * time.Second

			// Run worker — returns on error or context cancellation
			err = e.runIMAPWorker(ctx, c, a, selectedMbox.NumMessages)
			c.Close()

			if err != nil && err != context.Canceled {
				xlog.Error(fmt.Sprintf("IMAP worker error: %v (reconnect in %v)", err, backoff))
				select {
				case <-ctx.Done():
					xlog.Info("Email connector is now stopped.")
					return
				case <-time.After(backoff):
				}
				backoff = min(backoff*2, maxBackoff)
			} else if err == nil || err == context.Canceled {
				xlog.Info("Email connector is now stopped.")
				return
			}
		}
	}()
}
