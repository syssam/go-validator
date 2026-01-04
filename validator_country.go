package validator

import "strings"

// ISO 3166-1 alpha-2 country codes (2-letter)
var countryAlpha2 = map[string]bool{
	"AD": true, "AE": true, "AF": true, "AG": true, "AI": true, "AL": true, "AM": true, "AO": true,
	"AQ": true, "AR": true, "AS": true, "AT": true, "AU": true, "AW": true, "AX": true, "AZ": true,
	"BA": true, "BB": true, "BD": true, "BE": true, "BF": true, "BG": true, "BH": true, "BI": true,
	"BJ": true, "BL": true, "BM": true, "BN": true, "BO": true, "BQ": true, "BR": true, "BS": true,
	"BT": true, "BV": true, "BW": true, "BY": true, "BZ": true, "CA": true, "CC": true, "CD": true,
	"CF": true, "CG": true, "CH": true, "CI": true, "CK": true, "CL": true, "CM": true, "CN": true,
	"CO": true, "CR": true, "CU": true, "CV": true, "CW": true, "CX": true, "CY": true, "CZ": true,
	"DE": true, "DJ": true, "DK": true, "DM": true, "DO": true, "DZ": true, "EC": true, "EE": true,
	"EG": true, "EH": true, "ER": true, "ES": true, "ET": true, "FI": true, "FJ": true, "FK": true,
	"FM": true, "FO": true, "FR": true, "GA": true, "GB": true, "GD": true, "GE": true, "GF": true,
	"GG": true, "GH": true, "GI": true, "GL": true, "GM": true, "GN": true, "GP": true, "GQ": true,
	"GR": true, "GS": true, "GT": true, "GU": true, "GW": true, "GY": true, "HK": true, "HM": true,
	"HN": true, "HR": true, "HT": true, "HU": true, "ID": true, "IE": true, "IL": true, "IM": true,
	"IN": true, "IO": true, "IQ": true, "IR": true, "IS": true, "IT": true, "JE": true, "JM": true,
	"JO": true, "JP": true, "KE": true, "KG": true, "KH": true, "KI": true, "KM": true, "KN": true,
	"KP": true, "KR": true, "KW": true, "KY": true, "KZ": true, "LA": true, "LB": true, "LC": true,
	"LI": true, "LK": true, "LR": true, "LS": true, "LT": true, "LU": true, "LV": true, "LY": true,
	"MA": true, "MC": true, "MD": true, "ME": true, "MF": true, "MG": true, "MH": true, "MK": true,
	"ML": true, "MM": true, "MN": true, "MO": true, "MP": true, "MQ": true, "MR": true, "MS": true,
	"MT": true, "MU": true, "MV": true, "MW": true, "MX": true, "MY": true, "MZ": true, "NA": true,
	"NC": true, "NE": true, "NF": true, "NG": true, "NI": true, "NL": true, "NO": true, "NP": true,
	"NR": true, "NU": true, "NZ": true, "OM": true, "PA": true, "PE": true, "PF": true, "PG": true,
	"PH": true, "PK": true, "PL": true, "PM": true, "PN": true, "PR": true, "PS": true, "PT": true,
	"PW": true, "PY": true, "QA": true, "RE": true, "RO": true, "RS": true, "RU": true, "RW": true,
	"SA": true, "SB": true, "SC": true, "SD": true, "SE": true, "SG": true, "SH": true, "SI": true,
	"SJ": true, "SK": true, "SL": true, "SM": true, "SN": true, "SO": true, "SR": true, "SS": true,
	"ST": true, "SV": true, "SX": true, "SY": true, "SZ": true, "TC": true, "TD": true, "TF": true,
	"TG": true, "TH": true, "TJ": true, "TK": true, "TL": true, "TM": true, "TN": true, "TO": true,
	"TR": true, "TT": true, "TV": true, "TW": true, "TZ": true, "UA": true, "UG": true, "UM": true,
	"US": true, "UY": true, "UZ": true, "VA": true, "VC": true, "VE": true, "VG": true, "VI": true,
	"VN": true, "VU": true, "WF": true, "WS": true, "YE": true, "YT": true, "ZA": true, "ZM": true,
	"ZW": true,
}

// ISO 3166-1 alpha-3 country codes (3-letter)
var countryAlpha3 = map[string]bool{
	"ABW": true, "AFG": true, "AGO": true, "AIA": true, "ALA": true, "ALB": true, "AND": true, "ARE": true,
	"ARG": true, "ARM": true, "ASM": true, "ATA": true, "ATF": true, "ATG": true, "AUS": true, "AUT": true,
	"AZE": true, "BDI": true, "BEL": true, "BEN": true, "BES": true, "BFA": true, "BGD": true, "BGR": true,
	"BHR": true, "BHS": true, "BIH": true, "BLM": true, "BLR": true, "BLZ": true, "BMU": true, "BOL": true,
	"BRA": true, "BRB": true, "BRN": true, "BTN": true, "BVT": true, "BWA": true, "CAF": true, "CAN": true,
	"CCK": true, "CHE": true, "CHL": true, "CHN": true, "CIV": true, "CMR": true, "COD": true, "COG": true,
	"COK": true, "COL": true, "COM": true, "CPV": true, "CRI": true, "CUB": true, "CUW": true, "CXR": true,
	"CYM": true, "CYP": true, "CZE": true, "DEU": true, "DJI": true, "DMA": true, "DNK": true, "DOM": true,
	"DZA": true, "ECU": true, "EGY": true, "ERI": true, "ESH": true, "ESP": true, "EST": true, "ETH": true,
	"FIN": true, "FJI": true, "FLK": true, "FRA": true, "FRO": true, "FSM": true, "GAB": true, "GBR": true,
	"GEO": true, "GGY": true, "GHA": true, "GIB": true, "GIN": true, "GLP": true, "GMB": true, "GNB": true,
	"GNQ": true, "GRC": true, "GRD": true, "GRL": true, "GTM": true, "GUF": true, "GUM": true, "GUY": true,
	"HKG": true, "HMD": true, "HND": true, "HRV": true, "HTI": true, "HUN": true, "IDN": true, "IMN": true,
	"IND": true, "IOT": true, "IRL": true, "IRN": true, "IRQ": true, "ISL": true, "ISR": true, "ITA": true,
	"JAM": true, "JEY": true, "JOR": true, "JPN": true, "KAZ": true, "KEN": true, "KGZ": true, "KHM": true,
	"KIR": true, "KNA": true, "KOR": true, "KWT": true, "LAO": true, "LBN": true, "LBR": true, "LBY": true,
	"LCA": true, "LIE": true, "LKA": true, "LSO": true, "LTU": true, "LUX": true, "LVA": true, "MAC": true,
	"MAF": true, "MAR": true, "MCO": true, "MDA": true, "MDG": true, "MDV": true, "MEX": true, "MHL": true,
	"MKD": true, "MLI": true, "MLT": true, "MMR": true, "MNE": true, "MNG": true, "MNP": true, "MOZ": true,
	"MRT": true, "MSR": true, "MTQ": true, "MUS": true, "MWI": true, "MYS": true, "MYT": true, "NAM": true,
	"NCL": true, "NER": true, "NFK": true, "NGA": true, "NIC": true, "NIU": true, "NLD": true, "NOR": true,
	"NPL": true, "NRU": true, "NZL": true, "OMN": true, "PAK": true, "PAN": true, "PCN": true, "PER": true,
	"PHL": true, "PLW": true, "PNG": true, "POL": true, "PRI": true, "PRK": true, "PRT": true, "PRY": true,
	"PSE": true, "PYF": true, "QAT": true, "REU": true, "ROU": true, "RUS": true, "RWA": true, "SAU": true,
	"SDN": true, "SEN": true, "SGP": true, "SGS": true, "SHN": true, "SJM": true, "SLB": true, "SLE": true,
	"SLV": true, "SMR": true, "SOM": true, "SPM": true, "SRB": true, "SSD": true, "STP": true, "SUR": true,
	"SVK": true, "SVN": true, "SWE": true, "SWZ": true, "SXM": true, "SYC": true, "SYR": true, "TCA": true,
	"TCD": true, "TGO": true, "THA": true, "TJK": true, "TKL": true, "TKM": true, "TLS": true, "TON": true,
	"TTO": true, "TUN": true, "TUR": true, "TUV": true, "TWN": true, "TZA": true, "UGA": true, "UKR": true,
	"UMI": true, "URY": true, "USA": true, "UZB": true, "VAT": true, "VCT": true, "VEN": true, "VGB": true,
	"VIR": true, "VNM": true, "VUT": true, "WLF": true, "WSM": true, "YEM": true, "ZAF": true, "ZMB": true,
	"ZWE": true,
}

// ISO 3166-1 numeric country codes
var countryNumeric = map[string]bool{
	"004": true, "008": true, "010": true, "012": true, "016": true, "020": true, "024": true, "028": true,
	"031": true, "032": true, "036": true, "040": true, "044": true, "048": true, "050": true, "051": true,
	"052": true, "056": true, "060": true, "064": true, "068": true, "070": true, "072": true, "074": true,
	"076": true, "084": true, "086": true, "090": true, "092": true, "096": true, "100": true, "104": true,
	"108": true, "112": true, "116": true, "120": true, "124": true, "132": true, "136": true, "140": true,
	"144": true, "148": true, "152": true, "156": true, "158": true, "162": true, "166": true, "170": true,
	"174": true, "175": true, "178": true, "180": true, "184": true, "188": true, "191": true, "192": true,
	"196": true, "203": true, "204": true, "208": true, "212": true, "214": true, "218": true, "222": true,
	"226": true, "231": true, "232": true, "233": true, "234": true, "238": true, "242": true, "246": true,
	"248": true, "250": true, "254": true, "258": true, "260": true, "262": true, "266": true, "268": true,
	"270": true, "275": true, "276": true, "288": true, "292": true, "296": true, "300": true, "304": true,
	"308": true, "312": true, "316": true, "320": true, "324": true, "328": true, "332": true, "334": true,
	"336": true, "340": true, "344": true, "348": true, "352": true, "356": true, "360": true, "364": true,
	"368": true, "372": true, "376": true, "380": true, "384": true, "388": true, "392": true, "398": true,
	"400": true, "404": true, "408": true, "410": true, "414": true, "417": true, "418": true, "422": true,
	"426": true, "428": true, "430": true, "434": true, "438": true, "440": true, "442": true, "446": true,
	"450": true, "454": true, "458": true, "462": true, "466": true, "470": true, "474": true, "478": true,
	"480": true, "484": true, "492": true, "496": true, "498": true, "499": true, "500": true, "504": true,
	"508": true, "512": true, "516": true, "520": true, "524": true, "528": true, "531": true, "533": true,
	"534": true, "535": true, "540": true, "548": true, "554": true, "558": true, "562": true, "566": true,
	"570": true, "574": true, "578": true, "580": true, "581": true, "583": true, "584": true, "585": true,
	"586": true, "591": true, "598": true, "600": true, "604": true, "608": true, "612": true, "616": true,
	"620": true, "624": true, "626": true, "630": true, "634": true, "638": true, "642": true, "643": true,
	"646": true, "652": true, "654": true, "659": true, "660": true, "662": true, "663": true, "666": true,
	"670": true, "674": true, "678": true, "682": true, "686": true, "688": true, "690": true, "694": true,
	"702": true, "703": true, "704": true, "705": true, "706": true, "710": true, "716": true, "724": true,
	"728": true, "729": true, "732": true, "740": true, "744": true, "748": true, "752": true, "756": true,
	"760": true, "762": true, "764": true, "768": true, "772": true, "776": true, "780": true, "784": true,
	"788": true, "792": true, "795": true, "796": true, "798": true, "800": true, "804": true, "807": true,
	"818": true, "826": true, "831": true, "832": true, "833": true, "834": true, "840": true, "850": true,
	"854": true, "858": true, "860": true, "862": true, "876": true, "882": true, "887": true, "894": true,
}

// IsCountryAlpha2 validates ISO 3166-1 alpha-2 country code (2-letter)
func IsCountryAlpha2(str string) bool {
	return countryAlpha2[strings.ToUpper(str)]
}

// IsCountryAlpha3 validates ISO 3166-1 alpha-3 country code (3-letter)
func IsCountryAlpha3(str string) bool {
	return countryAlpha3[strings.ToUpper(str)]
}

// IsCountryNumeric validates ISO 3166-1 numeric country code
func IsCountryNumeric(str string) bool {
	// Pad with leading zeros if needed
	for len(str) < 3 {
		str = "0" + str
	}
	return countryNumeric[str]
}

// IsCountryCode validates any ISO 3166-1 country code format
func IsCountryCode(str string) bool {
	return IsCountryAlpha2(str) || IsCountryAlpha3(str) || IsCountryNumeric(str)
}
