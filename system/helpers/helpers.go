package helpers

import "fmt"

var NotAdmin string = "You are not an admin"
var BotNotAdmin string = "Bot not admin"
var NotGroup string = "Not a group"
var NotOwner string = "Owner only"
var NotRegisteredNum string = "The number is not registered on WhatsApp"
var Wait = "⏱️"
var Success string = "✅"
var Warning string = "⚠️"
var Failed string = "❌"
var InputQuery = "Query is required"

func ExampleUse(prefix, command string) string {
	ex := fmt.Sprintf("Example of use:\n%s%s", prefix, command)
	return ex
}
