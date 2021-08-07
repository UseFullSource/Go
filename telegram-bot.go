package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tb "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	BOT_TOKEN  string = ""
	DATE_URL   string = "http://www.persiancalapi.ir/jalali"
	ROOT_ID    int    = //root user id
	MY_CHAT_ID int64  = //your chat id
)

type Date struct {
	IsHoliday bool     `json:"is_holiday"`
	Events    []Events `json:"events"`
}

type Events struct {
	Description           string `json:"description"`
	IsRelegious           bool   `json:"is_relegious"`
	AdditionalDescription string `json:"additional_description"`
}

var botInlineKeyboard = tb.NewInlineKeyboardMarkup(
	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("add", "⚡ اضافه کردن ادمین جدید ⚡"),
		tb.NewInlineKeyboardButtonData("️group️", "⚡ نمایش لیست گروه ها ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("️disadd", "⚡ حذف از ادمین ها ⚡"),
		tb.NewInlineKeyboardButtonData("setting", "⚡ نمایش پنل تنظیمات کاربر ریشه ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("schedule", "⚡ زمانبندی رخدادهای گروه ⚡"),
		tb.NewInlineKeyboardButtonData("+", "⚡ تشکر از کاربر ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("admins", "⚡ نمایش لیست ادمین ها ⚡"),
		tb.NewInlineKeyboardButtonData("️channel️", "⚡ نمایش لیست کانال ها ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("mute", "⚡ قرار دادن کاربر در حالت سکوت ⚡"),
		tb.NewInlineKeyboardButtonData("!@", "⚡ دادن اخطار عدم منشن کردن ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("date", "⚡ نمایش رخداد ها در تقویم رسمی ⚡"),
		tb.NewInlineKeyboardButtonData("about", "⚡ درباره ما ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("unmute", "⚡ خارج کردن کاربر از حالت سکوت ⚡"),
		tb.NewInlineKeyboardButtonData("ban", "⚡ اخراج کاربر از انجمن ⚡"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("command", "⚡ نمایش تمام دستورات ربات ⚡"),
	),
)

var captchaKeyboard = tb.NewInlineKeyboardMarkup(
	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("اینجا کلیک کنید", "اینجا کلیک کنید"),
	),
)

var settingKeyboard = tb.NewInlineKeyboardMarkup(
	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("lockGP", "lockGP"),
		tb.NewInlineKeyboardButtonData("unlockGP", "unlockGP"),
	),

	tb.NewInlineKeyboardRow(
		tb.NewInlineKeyboardButtonData("changeTitle", "changeTitle"),
		tb.NewInlineKeyboardButtonData("deleteMsg", "deleteMsg"),
	),
)

var botKeyboard = tb.NewReplyKeyboard(
	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("⚡ارسال پیام برای انجمن⚡"),
	),

	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("🔰لینک ها🔰"),
	),

	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("💠درباره انجمن💠"),
	),

	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("🔧تنظیمات🔧"),
	),

	tb.NewKeyboardButtonRow(
		tb.NewKeyboardButton("🔥درباره سازنده ربات🔥"),
	),
)

func dateFetcher(year, month, day int) (string, error) {
	response, err := http.Get(DATE_URL + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/" + strconv.Itoa(day))

	if err != nil {

		return "", err
	}

	holiday := "نیست"

	r := &Date{}

	err = json.NewDecoder(response.Body).Decode(r)

	if r.IsHoliday {
		holiday = "است"
	}

	var m = ""

	switch month {
	case 1:
		m = "فروردین"
		break
	case 2:
		m = "اردیبهشت"
		break
	case 3:
		m = "خرداد"
		break
	case 4:
		m = "تیر"
		break
	case 5:
		m = "مرداد"
		break
	case 6:
		m = "شهریور"
		break
	case 7:
		m = "مهر"
		break
	case 8:
		m = "آبان"
		break
	case 9:
		m = "آذر"
		break
	case 10:
		m = "دی"
		break
	case 11:
		m = "بهمن"
		break
	case 12:
		m = "اسفند"
		break
	}

	msg := "💡 این روز تعطیل " + holiday + " 💡\n\n"
	title := "💠<b>رخداد های مهم روز " + strconv.Itoa(day) + " " + m + " سال " + strconv.Itoa(year) + " : </b>\n\n"
	events := ""

	for _, date := range r.Events {
		events += "🎈 " + date.Description + "\n"

		if date.IsRelegious {
			events += "🕌 این یک رخداد مذهبی است\n"
		}

		if date.AdditionalDescription != "" {
			events += "🗽 این رخداد در تاریخ میلادی " + date.AdditionalDescription + " به وقوع پیوسته است\n\n"
		}
	}

	if events == "" {
		events = msg + title + "⚓ هیچ رخداد مهم ثبت شده ایی در این روز وجود ندارد"
	} else {

		events = msg + title + events
	}

	if response.Status != "200 OK" {
		events = "bad response"
	}

	return events, err
}

func main() {
	bot, err := tb.NewBotAPI(BOT_TOKEN)

	if err != nil {
		log.Println(err)
	}

	bot.Debug = false

	end := true
	lastUpdateID := 0

	u := tb.NewUpdate(lastUpdateID)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	updates.Clear()

	if err != nil {
		log.Println(err)
	}

	for update := range updates { //HANDLE UPDATE

		if update.CallbackQuery != nil { //HANDLE CALLBACK QUERY UPDATES
			if update.CallbackQuery.Data == "اینجا کلیک کنید" { //HANEL USER JOINED MESSAGE'S BUTTON
				canSendMessage := true
				restrict := tb.RestrictChatMemberConfig{
					ChatMemberConfig: tb.ChatMemberConfig{
						UserID: update.CallbackQuery.From.ID,
						ChatID: update.CallbackQuery.Message.Chat.ID,
					},
					CanSendMessages:       &canSendMessage,
					CanSendOtherMessages:  &canSendMessage,
					CanSendMediaMessages:  &canSendMessage,
					CanAddWebPagePreviews: &canSendMessage,
				}
				bot.RestrictChatMember(restrict)
				continue
			} else { //HANDLE SETTING MESSSAGE'S BUTTON
				if update.CallbackQuery.From.ID == ROOT_ID {
					switch update.CallbackQuery.Data {
					case "lockGP":
						msg := "🔰 گروه به درخواست *کاربر ریشه* قفل شد"
						new_msg := tb.NewMessage(update.CallbackQuery.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						bot.Send(new_msg)
						continue
					case "unlockGP":
						msg := "🔰 گروه به درخواست *کاربر ریشه* باز شد"
						new_msg := tb.NewMessage(update.CallbackQuery.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						bot.Send(new_msg)
						continue
					case "changeTitle":
						continue
					case "deleteMsg":
						continue
					}
				}
				//HANDLE GROUP COMMANDS CALLBACK QUERY
				callback := tb.CallbackConfig{
					ShowAlert:       true,
					Text:            update.CallbackQuery.Data,
					CallbackQueryID: update.CallbackQuery.ID,
					CacheTime:       update.CallbackQuery.Message.Date,
				}

				bot.AnswerCallbackQuery(callback)
				continue

			}
		}

		if update.ChannelPost != nil { //HANDLE CHANNEL POSTS
			if strings.Contains(update.ChannelPost.Text, "#ارسال") {
				new_msg := tb.NewForward(MY_CHAT_ID, update.ChannelPost.Chat.ID, update.ChannelPost.MessageID)
				bot.Send(new_msg)
			}
			continue
		}

		if update.Message == nil { //SKIP NON MESSAGE UPDATES
			continue
		}

		if update.Message.Chat.IsPrivate() == false { //HANDLE GROUP COMMANDS
			//HANDLE NEW MEMBER JOINED
			if update.Message.NewChatMembers != nil {
				if len(*update.Message.NewChatMembers) != 0 {
					canSendMessage := false

					for _, newMember := range *update.Message.NewChatMembers {
						msg := "🍀 سلام [" + newMember.FirstName + "](tg://user?id=" + strconv.Itoa(newMember.ID) + ") عزیز به گروه *" + update.Message.Chat.Title + "* خوش آمدید\n\n🔆 جهت اطمینان از *ربات نبودن* شما روی دکمه زیر کلیک کنید تا محدودیت ارسال پیام برای شما رفع شود\n\n🔰 مدت زمان باقی مانده برای تایید شما ۲ ساعت میباشد"
						new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						new_msg.ReplyMarkup = captchaKeyboard
						new_msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(new_msg)

						restrict := tb.RestrictChatMemberConfig{
							ChatMemberConfig: tb.ChatMemberConfig{
								UserID: newMember.ID,
								ChatID: update.Message.Chat.ID,
							},
							CanSendMessages:       &canSendMessage,
							CanSendOtherMessages:  &canSendMessage,
							CanSendMediaMessages:  &canSendMessage,
							CanAddWebPagePreviews: &canSendMessage,
						}

						bot.RestrictChatMember(restrict)
					}
					continue
				}
			}

			if end == false {
				if update.Message.From.ID != ROOT_ID {
					new_delete_msg := tb.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					bot.DeleteMessage(new_delete_msg)
				}
				continue
			}
			//HANDLE DATE
			cmd := strings.Fields(update.Message.Text)
			if len(cmd) == 4 {
				if strings.ToLower(cmd[0]) == "date" {
					year, _ := strconv.ParseInt(cmd[1], 10, 64)
					month, _ := strconv.ParseInt(cmd[2], 10, 64)
					day, _ := strconv.ParseInt(cmd[3], 10, 64)

					events, err := dateFetcher(int(year), int(month), int(day))

					if err != nil {
						log.Println(err)
					}

					if events != "bad response" {
						new_msg := tb.NewMessage(update.Message.Chat.ID, events)
						new_msg.ParseMode = "HTML"
						new_msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(new_msg)
						continue
					}
				}
			}

			switch update.Message.Text { //HANDLE GOROP COMMANDS
			case "command":
				botCommandsMessage := "🔆<b>راهنمای کامل دستورات ربات</b>🔆\n\n🔱 برای مشاهده راهنمای هر دستور روی آن کلیک کنید\n\n❗ کاربر با سطح دسترسی بالاتر میتواند دستورات سطوح پایین خود را اجرا کند"
				new_msg := tb.NewMessage(update.Message.Chat.ID, botCommandsMessage)
				new_msg.ParseMode = "HTML"
				new_msg.ReplyToMessageID = update.Message.MessageID
				new_msg.ReplyMarkup = botInlineKeyboard
				bot.Send(new_msg)
				continue
			case "mute":
				if update.Message.ReplyToMessage != nil {
					chatConf := tb.ChatConfig{
						ChatID: MY_CHAT_ID,
					}

					admins, _ := bot.GetChatAdministrators(chatConf)

					for _, admin := range admins {
						if update.Message.From.ID == admin.User.ID {
							resUntil := time.Now().Add(time.Hour * 2).Unix()

							canSendMessage := false
							restrict := tb.RestrictChatMemberConfig{
								ChatMemberConfig: tb.ChatMemberConfig{
									UserID: update.Message.ReplyToMessage.From.ID,
									ChatID: update.Message.Chat.ID,
								},
								CanSendMessages:       &canSendMessage,
								CanSendOtherMessages:  &canSendMessage,
								CanSendMediaMessages:  &canSendMessage,
								CanAddWebPagePreviews: &canSendMessage,
								UntilDate:             resUntil,
							}

							res, _ := bot.RestrictChatMember(restrict)

							if res.Ok {
								user_id := update.Message.ReplyToMessage.From.ID
								user_name := update.Message.ReplyToMessage.From.FirstName
								admin_id := update.Message.From.ID
								admin_name := update.Message.From.FirstName
								msg := "🔰 کاربر [" + user_name + "](tg://user?id=" + strconv.Itoa(user_id) + ") به درخواست [" + admin_name + "](" + "tg://user?id=" + strconv.Itoa(admin_id) + ") به مدت ۲ ساعت در حالت سکوت قرار گرفت."

								new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
								new_msg.ParseMode = "Markdown"
								new_msg.ReplyToMessageID = update.Message.MessageID
								bot.Send(new_msg)
							}

							break
						}
					}
					continue
				}
				continue
			case "unmute":
				if update.Message.ReplyToMessage != nil {
					chatConf := tb.ChatConfig{
						ChatID: MY_CHAT_ID,
					}

					admins, _ := bot.GetChatAdministrators(chatConf)

					for _, admin := range admins {
						if update.Message.From.ID == admin.User.ID {
							resUntil := time.Now().Add(time.Hour * 2).Unix()

							canSendMessage := true
							restrict := tb.RestrictChatMemberConfig{
								ChatMemberConfig: tb.ChatMemberConfig{
									UserID: update.Message.ReplyToMessage.From.ID,
									ChatID: update.Message.Chat.ID,
								},
								CanSendMessages:       &canSendMessage,
								CanSendOtherMessages:  &canSendMessage,
								CanSendMediaMessages:  &canSendMessage,
								CanAddWebPagePreviews: &canSendMessage,
								UntilDate:             resUntil,
							}

							res, _ := bot.RestrictChatMember(restrict)

							if res.Ok {
								user_id := update.Message.ReplyToMessage.From.ID
								user_name := update.Message.ReplyToMessage.From.FirstName
								admin_id := update.Message.From.ID
								admin_name := update.Message.From.FirstName
								msg := "🔰 کاربر [" + user_name + "](tg://user?id=" + strconv.Itoa(user_id) + ") به درخواست [" + admin_name + "](" + "tg://user?id=" + strconv.Itoa(admin_id) + ") از حالت سکوت خارج شد."

								new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
								new_msg.ParseMode = "Markdown"
								new_msg.ReplyToMessageID = update.Message.MessageID
								bot.Send(new_msg)
							}

							break
						}
					}
					continue
				}

				continue

			case "add":
				if update.Message.ReplyToMessage != nil {
					if update.Message.From.ID == ROOT_ID {
						user_id := update.Message.ReplyToMessage.From.ID
						user_name := update.Message.ReplyToMessage.From.FirstName
						admin_id := ROOT_ID
						admin_name := update.Message.From.FirstName

						canAdminDo := true

						PromoteMember := tb.PromoteChatMemberConfig{
							ChatMemberConfig: tb.ChatMemberConfig{
								ChatID:             update.Message.Chat.ID,
								UserID:             user_id,
								SuperGroupUsername: "@RaziCE",
							},
							CanDeleteMessages:  &canAdminDo,
							CanInviteUsers:     &canAdminDo,
							CanRestrictMembers: &canAdminDo,
							CanPinMessages:     &canAdminDo,
						}

						res, err := bot.PromoteChatMember(PromoteMember)

						if err != nil {
							log.Println(err)
						}

						if res.Ok {
							msg := "🚀 کاربر [" + user_name + "](tg://user?id=" + strconv.Itoa(user_id) + ") توسط [" + admin_name + "](tg://user?id=" + strconv.Itoa(admin_id) + ") به سطح *ادمین* ارتقا یافت"
							new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
							new_msg.ParseMode = "Markdown"
							new_msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(new_msg)
							continue
						}
					} else {
						msg := "💡 *کاربر گرامی شما مجوز اجرای این دستور را *ندارید"
						new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						new_msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(new_msg)
						continue
					}
				}
				continue
			case "!@":
				if update.Message.ReplyToMessage != nil {
					msg := "💡 کاربر گرامی لطفا از منشن کردن دیگران که باعث ایجاد مزاحمت میشود خودداری کنید"
					new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
					new_msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
					bot.Send(new_msg)
					continue
				}
				continue
			case "disadd":
				if update.Message.ReplyToMessage != nil {
					if update.Message.From.ID == ROOT_ID {
						user_id := update.Message.ReplyToMessage.From.ID
						user_name := update.Message.ReplyToMessage.From.FirstName
						admin_id := ROOT_ID
						admin_name := update.Message.From.FirstName

						canAdminDo := false

						PromoteMember := tb.PromoteChatMemberConfig{
							ChatMemberConfig: tb.ChatMemberConfig{
								ChatID:             update.Message.Chat.ID,
								UserID:             user_id,
								SuperGroupUsername: "@RaziCE",
							},
							CanRestrictMembers: &canAdminDo,
							CanEditMessages:    &canAdminDo,
							CanInviteUsers:     &canAdminDo,
							CanPinMessages:     &canAdminDo,
						}

						res, err := bot.PromoteChatMember(PromoteMember)

						if err != nil {
							log.Println(err)
						}

						if res.Ok {
							msg := "🚀 کاربر [" + user_name + "](tg://user?id=" + strconv.Itoa(user_id) + ") توسط [" + admin_name + "](tg://user?id=" + strconv.Itoa(admin_id) + ") از سطح *ادمین* به سطح *کاربر عادی* انتقال یافت"
							new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
							new_msg.ParseMode = "Markdown"
							new_msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(new_msg)
							continue
						}
					} else {
						msg := "💡 *کاربر گرامی شما مجوز اجرای این دستور را *ندارید"
						new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						new_msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(new_msg)
						continue
					}
				}
				continue
			case "setting":
				if update.Message.From.ID == ROOT_ID {
					end = false
					msg := "🔰 *کاربر ریشه* در حال انجام تنظیمات مدیریتی گروه است. تا اتمام این فرایند امکان ارسال پیام برای *کاربران عادی* و *ادمین ها* وجود ندارد"
					msg += "\n\n💡 *کاربر ریشه ی گرامی لطفا برای اتمام فرایند تنظیمات کلمه ی end را با حروف کوچک بفرستید*"
					new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
					new_msg.ParseMode = "Markdown"
					new_msg.DisableNotification = false
					new_msg.ReplyMarkup = settingKeyboard
					new_msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(new_msg)
					continue
				} else {
					if end == false {
						msg := "💡 *کاربر گرامی شما مجوز اجرای این دستور را *ندارید"
						new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
						new_msg.ParseMode = "Markdown"
						new_msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(new_msg)
					}
					continue
				}
			case "end":
				if update.Message.From.ID == ROOT_ID {
					end = true
					msg := "🔰 پایان تنظیمات توسط *کاربر ریشه* اعلام گردید. هم اکنون ارسال پیام در گروه امکان پذیر است"
					new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
					new_msg.ParseMode = "Markdown"
					bot.Send(new_msg)
				}
				continue
			}

		} else { //HANDLE PRIVATE MESSAGES

			if update.Message.IsCommand() { //HANDLE PRIVATE COMMANDS
				command := update.Message.Command()
				switch command {
				case "start":
					msg := "🔆 <b>از منوی زیر انتخاب کنید</b> 🔆"
					new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
					new_msg.ParseMode = "HTML"
					new_msg.ReplyMarkup = botKeyboard
					bot.Send(new_msg)
					break
				}

				continue
			} else { //HANDLE REGULAR PRIVATE MESSAGES
				switch update.Message.Text {
				case "⚡ارسال پیام برای انجمن⚡":
					break
				case "🔰لینک ها🔰":
					break
				case "💠درباره انجمن💠":
					break
				case "🔧تنظیمات🔧":
					if update.Message.From.ID == ROOT_ID {

					} else {
						new_msg := tb.NewMessage(update.Message.Chat.ID, "💡 شما مجوز انجام تنظیمات را ندارید 💡")
						bot.Send(new_msg)
					}
					break
				case "🔥درباره سازنده ربات🔥":
					msg := "🎯 این ربات توسط [| D3741L |](tg://user?id=" + strconv.Itoa(ROOT_ID) + ") ساخته شده است.\n\n🍀 برای سفارش ساخت ربات میتوانید با سازنده ربات در ارتباط باشید"
					new_msg := tb.NewMessage(update.Message.Chat.ID, msg)
					new_msg.ParseMode = "Markdown"
					bot.Send(new_msg)
					break
				}
				continue
			}
		}
	}
}
