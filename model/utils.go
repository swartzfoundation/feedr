package model

import (
	"encoding/base32"
	"strings"
	"unicode"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ContextKey string
type UUID []byte

var caser = cases.Title(language.BritishEnglish)

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769").WithPadding(base32.NoPadding)

// NewID  is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// without the padding.
func NewIDEncode() (string, error) {
	v4, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return encoding.EncodeToString(UUID(v4[:])), nil
}

func NewBase32ID() string {
	return xid.New().String()
}

func NewID() string {
	// zero, o, 1, i, and I are removed to avoid confusion with 0 and 1.
	alpha := "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjklmnpqrstuvwxyz"
	id, err := gonanoid.Generate(alpha, 9)
	if err != nil {
		return NewBase32ID()
	}
	return id
}

func HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	// As of today go still uses 10 as default cost. Next versions might bump to 12 or higher
	cost := max(bcrypt.DefaultCost, 12)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, cost)
	return string(hash), err
}

func CheckPassword(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

// filterBlocklist returns `r` if it is not in the blocklist, otherwise drop (-1).
// Blocklist is taken from https://www.w3.org/TR/unicode-xml/#Charlist
func filterBlocklist(r rune) rune {
	const drop = -1
	switch r {
	case '\u0340', '\u0341': // clones of grave and acute; deprecated in Unicode `
		return drop
	case '\u17A3', '\u17D3': // obsolete characters for Khmer; deprecated in Unicode
		return drop
	case '\u2028', '\u2029': // line and paragraph separator
		return drop
	case '\u202A', '\u202B', '\u202C', '\u202D', '\u202E': // BIDI embedding controls
		return drop
	case '\u206A', '\u206B': // activate/inhibit symmetric swapping; deprecated in Unicode
		return drop
	case '\u206C', '\u206D': // activate/inhibit Arabic form shaping; deprecated in Unicode
		return drop
	case '\u206E', '\u206F': // activate/inhibit national digit shapes; deprecated in Unicode
		return drop
	case '\uFFF9', '\uFFFA', '\uFFFB': // interlinear annotation characters
		return drop
	case '\uFEFF': // byte order mark
		return drop
	case '\uFFFC': // object replacement character
		return drop
	}

	// Scoping for musical notation
	if r >= 0x0001D173 && r <= 0x0001D17A {
		return drop
	}

	// Language tag code points
	if r >= 0x000E0000 && r <= 0x000E007F {
		return drop
	}
	return r
}

// SanitizeUnicode will remove undesirable Unicode characters from a string.
func SanitizeUnicode(s string) string {
	return strings.Map(filterBlocklist, s)
}

func ToTitleCase(v string) string {
	return caser.String(v)
}

func ToSentenceCase(str string) string {
	if len(str) == 0 {
		return str
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
