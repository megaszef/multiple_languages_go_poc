package main

import (
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("./lang/active.es.toml")
	bundle.MustLoadMessageFile("./lang/active.ru.toml")

	engine := html.New("./templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		lang := c.Query("lang")
		accept := c.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",
				Other: "Hello {{.Name}}",
			},
			TemplateData: &fiber.Map{
				"Name": "John",
			},
		})

		unreadEmailCount, _ := strconv.ParseInt(c.Query("unread"), 10, 64)

		unreadEmailConfig := &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "MyUnreadEmails",
				One:   "You have {{.PluralCount}} unread email.",
				Other: "You have {{.PluralCount}} unread emails.",
			},
			PluralCount: unreadEmailCount,
		}

		unreadEmails := localizer.MustLocalize(unreadEmailConfig)

		if c.Query("format") == "json" {
			return c.JSON(&fiber.Map{
				"name":          helloPerson,
				"unread_emails": unreadEmails,
			})
		}

		return c.Render("index", fiber.Map{
			"Title":        helloPerson,
			"UnreadEmails": unreadEmails,
		})
	})
	log.Fatal(app.Listen(":3000"))
}
