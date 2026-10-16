package render

import "slices"

// ConsentMessages maps each consent.* translation key to its English text. It
// is the contract between host-manager's per-key fallbacks and nexus-manager's
// bundled translation.
var ConsentMessages = map[string]string{
	"consent.page-title":         "Allow access?",
	"consent.title":              "Allow access to your account?",
	"consent.signed-in-as":       "Signed in as %s.",
	"consent.web.sends-to":       "Allowing sends access to:",
	"consent.web.website":        "Website: %s",
	"consent.device.sends-to":    "Allowing sends access to %s (%s).",
	"consent.device.application": "an application on this device",
	"consent.full-address":       "Full address: %s",
	"consent.client-name":        "The application calls itself %s. This name is provided by the application and has not been verified.",
	"consent.permissions":        "It will be able to act on your account with your permissions.",
	"consent.permissions-scopes": "It will be able to act on your account with your permissions and read: %s.",
	"consent.warning":            "Only allow this if you started connecting this application yourself.",
	"consent.allow":              "Allow",
	"consent.deny":               "Deny",
}

// ConsentMessageKeys returns the keys of ConsentMessages, sorted.
func ConsentMessageKeys() []string {
	keys := make([]string, 0, len(ConsentMessages))
	for k := range ConsentMessages {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
