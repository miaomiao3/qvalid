package qvalid

import (
	"regexp"
)

// Basic regular expressions for validating strings
const (
	Email          string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	UpperAlpha     string = "^[A-Z]+$"
	LowerAlpha     string = "^[a-z]+$"
	Alpha          string = "^[a-zA-Z]+$"
	Alphanumeric   string = "^[a-zA-Z0-9]+$"
	Numeric        string = "^[0-9]+$"
	Int            string = "^(?:[-+]?(?:0|[1-9][0-9]*))$" // +- 0|1..
	Float          string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	Hexadecimal    string = "^[0-9a-fA-F]+$"
	ASCII          string = "^[\x00-\x7F]+$"
	MultiByte      string = "[^\x00-\x7F]"
	Base64         string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	PrintableASCII string = "^[\x20-\x7E]+$"
	DNSName        string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	IP             string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema      string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername    string = `(\S+(:\S*)?@)`
	URLPath        string = `((\/|\?|#)[^\s]*)`
	URLPort        string = `(:(\d{1,5}))`
	URLIP          string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	URLSubdomain   string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL            string = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	Semver         string = "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$"
)

const (
	validTag string = "valid"
)

const (
	StringTypeEmail        = "email"
	StringTypeAlpha        = "alpha"
	StringTypeUpperAlpha   = "upper_alpha"
	StringTypeLowerAlpha   = "lower_alpha"
	StringTypeAlphaNumeric = "alpha_numeric"
	StringTypeNumeric      = "numeric"
	StringTypeInt          = "int"
	StringTypeFloat        = "float"
	StringTypeHex          = "hex"
	StringTypeAscii        = "ascii"
	StringTypeVisibleAscii = "visible_ascii"
	StringTypeBytes        = "bytes"
	StringTypeBase64       = "base64"
	StringTypeDNS          = "dns"
	StringTypeVersion      = "version"
	StringTypeIp           = "ip"
	StringTypePort         = "port"
	StringTypeURL          = "url"
)

var stringRegexMap = map[string]*regexp.Regexp{
	StringTypeEmail:        regexp.MustCompile(Email),
	StringTypeAlpha:        regexp.MustCompile(Alpha),
	StringTypeUpperAlpha:   regexp.MustCompile(UpperAlpha),
	StringTypeLowerAlpha:   regexp.MustCompile(LowerAlpha),
	StringTypeAlphaNumeric: regexp.MustCompile(Alphanumeric),
	StringTypeNumeric:      regexp.MustCompile(Numeric),
	StringTypeInt:          regexp.MustCompile(Int),
	StringTypeFloat:        regexp.MustCompile(Float),
	StringTypeHex:          regexp.MustCompile(Hexadecimal),
	StringTypeAscii:        regexp.MustCompile(ASCII),
	StringTypeVisibleAscii: regexp.MustCompile(PrintableASCII),
	StringTypeBytes:        regexp.MustCompile(MultiByte),
	StringTypeBase64:       regexp.MustCompile(Base64),
	StringTypeDNS:          regexp.MustCompile(DNSName),
	StringTypeVersion:      regexp.MustCompile(Semver),
	StringTypeIp:           regexp.MustCompile(URLIP),
	StringTypePort:         regexp.MustCompile(URLPort),
	StringTypeURL:          regexp.MustCompile(URL),
}
