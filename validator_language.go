package validator

import "strings"

// ISO 639-1 language codes (2-letter)
var languageAlpha2 = map[string]bool{
	"aa": true, "ab": true, "ae": true, "af": true, "ak": true, "am": true, "an": true, "ar": true,
	"as": true, "av": true, "ay": true, "az": true, "ba": true, "be": true, "bg": true, "bh": true,
	"bi": true, "bm": true, "bn": true, "bo": true, "br": true, "bs": true, "ca": true, "ce": true,
	"ch": true, "co": true, "cr": true, "cs": true, "cu": true, "cv": true, "cy": true, "da": true,
	"de": true, "dv": true, "dz": true, "ee": true, "el": true, "en": true, "eo": true, "es": true,
	"et": true, "eu": true, "fa": true, "ff": true, "fi": true, "fj": true, "fo": true, "fr": true,
	"fy": true, "ga": true, "gd": true, "gl": true, "gn": true, "gu": true, "gv": true, "ha": true,
	"he": true, "hi": true, "ho": true, "hr": true, "ht": true, "hu": true, "hy": true, "hz": true,
	"ia": true, "id": true, "ie": true, "ig": true, "ii": true, "ik": true, "io": true, "is": true,
	"it": true, "iu": true, "ja": true, "jv": true, "ka": true, "kg": true, "ki": true, "kj": true,
	"kk": true, "kl": true, "km": true, "kn": true, "ko": true, "kr": true, "ks": true, "ku": true,
	"kv": true, "kw": true, "ky": true, "la": true, "lb": true, "lg": true, "li": true, "ln": true,
	"lo": true, "lt": true, "lu": true, "lv": true, "mg": true, "mh": true, "mi": true, "mk": true,
	"ml": true, "mn": true, "mr": true, "ms": true, "mt": true, "my": true, "na": true, "nb": true,
	"nd": true, "ne": true, "ng": true, "nl": true, "nn": true, "no": true, "nr": true, "nv": true,
	"ny": true, "oc": true, "oj": true, "om": true, "or": true, "os": true, "pa": true, "pi": true,
	"pl": true, "ps": true, "pt": true, "qu": true, "rm": true, "rn": true, "ro": true, "ru": true,
	"rw": true, "sa": true, "sc": true, "sd": true, "se": true, "sg": true, "si": true, "sk": true,
	"sl": true, "sm": true, "sn": true, "so": true, "sq": true, "sr": true, "ss": true, "st": true,
	"su": true, "sv": true, "sw": true, "ta": true, "te": true, "tg": true, "th": true, "ti": true,
	"tk": true, "tl": true, "tn": true, "to": true, "tr": true, "ts": true, "tt": true, "tw": true,
	"ty": true, "ug": true, "uk": true, "ur": true, "uz": true, "ve": true, "vi": true, "vo": true,
	"wa": true, "wo": true, "xh": true, "yi": true, "yo": true, "za": true, "zh": true, "zu": true,
}

// ISO 639-2 language codes (3-letter) - bibliographic codes
var languageAlpha3 = map[string]bool{
	"aar": true, "abk": true, "ace": true, "ach": true, "ada": true, "ady": true, "afa": true, "afh": true,
	"afr": true, "ain": true, "aka": true, "akk": true, "alb": true, "ale": true, "alg": true, "alt": true,
	"amh": true, "ang": true, "anp": true, "apa": true, "ara": true, "arc": true, "arg": true, "arm": true,
	"arn": true, "arp": true, "art": true, "arw": true, "asm": true, "ast": true, "ath": true, "aus": true,
	"ava": true, "ave": true, "awa": true, "aym": true, "aze": true, "bad": true, "bai": true, "bak": true,
	"bal": true, "bam": true, "ban": true, "baq": true, "bas": true, "bat": true, "bej": true, "bel": true,
	"bem": true, "ben": true, "ber": true, "bho": true, "bih": true, "bik": true, "bin": true, "bis": true,
	"bla": true, "bnt": true, "bod": true, "bos": true, "bra": true, "bre": true, "btk": true, "bua": true,
	"bug": true, "bul": true, "bur": true, "byn": true, "cad": true, "cai": true, "car": true, "cat": true,
	"cau": true, "ceb": true, "cel": true, "ces": true, "cha": true, "chb": true, "che": true, "chg": true,
	"chi": true, "chk": true, "chm": true, "chn": true, "cho": true, "chp": true, "chr": true, "chu": true,
	"chv": true, "chy": true, "cmc": true, "cnr": true, "cop": true, "cor": true, "cos": true, "cpe": true,
	"cpf": true, "cpp": true, "cre": true, "crh": true, "crp": true, "csb": true, "cus": true, "cym": true,
	"cze": true, "dak": true, "dan": true, "dar": true, "day": true, "del": true, "den": true, "deu": true,
	"dgr": true, "din": true, "div": true, "doi": true, "dra": true, "dsb": true, "dua": true, "dum": true,
	"dut": true, "dyu": true, "dzo": true, "efi": true, "egy": true, "eka": true, "ell": true, "elx": true,
	"eng": true, "enm": true, "epo": true, "est": true, "eus": true, "ewe": true, "ewo": true, "fan": true,
	"fao": true, "fas": true, "fat": true, "fij": true, "fil": true, "fin": true, "fiu": true, "fon": true,
	"fra": true, "fre": true, "frm": true, "fro": true, "frr": true, "frs": true, "fry": true, "ful": true,
	"fur": true, "gaa": true, "gay": true, "gba": true, "gem": true, "geo": true, "ger": true, "gez": true,
	"gil": true, "gla": true, "gle": true, "glg": true, "glv": true, "gmh": true, "goh": true, "gon": true,
	"gor": true, "got": true, "grb": true, "grc": true, "gre": true, "grn": true, "gsw": true, "guj": true,
	"gwi": true, "hai": true, "hat": true, "hau": true, "haw": true, "heb": true, "her": true, "hil": true,
	"him": true, "hin": true, "hit": true, "hmn": true, "hmo": true, "hrv": true, "hsb": true, "hun": true,
	"hup": true, "hye": true, "iba": true, "ibo": true, "ice": true, "ido": true, "iii": true, "ijo": true,
	"iku": true, "ile": true, "ilo": true, "ina": true, "inc": true, "ind": true, "ine": true, "inh": true,
	"ipk": true, "ira": true, "iro": true, "isl": true, "ita": true, "jav": true, "jbo": true, "jpn": true,
	"jpr": true, "jrb": true, "kaa": true, "kab": true, "kac": true, "kal": true, "kam": true, "kan": true,
	"kar": true, "kas": true, "kat": true, "kau": true, "kaw": true, "kaz": true, "kbd": true, "kha": true,
	"khi": true, "khm": true, "kho": true, "kik": true, "kin": true, "kir": true, "kmb": true, "kok": true,
	"kom": true, "kon": true, "kor": true, "kos": true, "kpe": true, "krc": true, "krl": true, "kro": true,
	"kru": true, "kua": true, "kum": true, "kur": true, "kut": true, "lad": true, "lah": true, "lam": true,
	"lao": true, "lat": true, "lav": true, "lez": true, "lim": true, "lin": true, "lit": true, "lol": true,
	"loz": true, "ltz": true, "lua": true, "lub": true, "lug": true, "lui": true, "lun": true, "luo": true,
	"lus": true, "mac": true, "mad": true, "mag": true, "mah": true, "mai": true, "mak": true, "mal": true,
	"man": true, "mao": true, "map": true, "mar": true, "mas": true, "may": true, "mdf": true, "mdr": true,
	"men": true, "mga": true, "mic": true, "min": true, "mis": true, "mkd": true, "mkh": true, "mlg": true,
	"mlt": true, "mnc": true, "mni": true, "mno": true, "moh": true, "mon": true, "mos": true, "mri": true,
	"msa": true, "mul": true, "mun": true, "mus": true, "mwl": true, "mwr": true, "mya": true, "myn": true,
	"myv": true, "nah": true, "nai": true, "nap": true, "nau": true, "nav": true, "nbl": true, "nde": true,
	"ndo": true, "nds": true, "nep": true, "new": true, "nia": true, "nic": true, "niu": true, "nld": true,
	"nno": true, "nob": true, "nog": true, "non": true, "nor": true, "nqo": true, "nso": true, "nub": true,
	"nwc": true, "nya": true, "nym": true, "nyn": true, "nyo": true, "nzi": true, "oci": true, "oji": true,
	"ori": true, "orm": true, "osa": true, "oss": true, "ota": true, "oto": true, "paa": true, "pag": true,
	"pal": true, "pam": true, "pan": true, "pap": true, "pau": true, "peo": true, "per": true, "phi": true,
	"phn": true, "pli": true, "pol": true, "pon": true, "por": true, "pra": true, "pro": true, "pus": true,
	"que": true, "raj": true, "rap": true, "rar": true, "roa": true, "roh": true, "rom": true, "ron": true,
	"rum": true, "run": true, "rup": true, "rus": true, "sad": true, "sag": true, "sah": true, "sai": true,
	"sal": true, "sam": true, "san": true, "sas": true, "sat": true, "scn": true, "sco": true, "sel": true,
	"sem": true, "sga": true, "sgn": true, "shn": true, "sid": true, "sin": true, "sio": true, "sit": true,
	"sla": true, "slk": true, "slo": true, "slv": true, "sma": true, "sme": true, "smi": true, "smj": true,
	"smn": true, "smo": true, "sms": true, "sna": true, "snd": true, "snk": true, "sog": true, "som": true,
	"son": true, "sot": true, "spa": true, "sqi": true, "srd": true, "srn": true, "srp": true, "srr": true,
	"ssa": true, "ssw": true, "suk": true, "sun": true, "sus": true, "sux": true, "swa": true, "swe": true,
	"syc": true, "syr": true, "tah": true, "tai": true, "tam": true, "tat": true, "tel": true, "tem": true,
	"ter": true, "tet": true, "tgk": true, "tgl": true, "tha": true, "tib": true, "tig": true, "tir": true,
	"tiv": true, "tkl": true, "tlh": true, "tli": true, "tmh": true, "tog": true, "ton": true, "tpi": true,
	"tsi": true, "tsn": true, "tso": true, "tuk": true, "tum": true, "tup": true, "tur": true, "tut": true,
	"tvl": true, "twi": true, "tyv": true, "udm": true, "uga": true, "uig": true, "ukr": true, "umb": true,
	"und": true, "urd": true, "uzb": true, "vai": true, "ven": true, "vie": true, "vol": true, "vot": true,
	"wak": true, "wal": true, "war": true, "was": true, "wel": true, "wen": true, "wln": true, "wol": true,
	"xal": true, "xho": true, "yao": true, "yap": true, "yid": true, "yor": true, "ypk": true, "zap": true,
	"zbl": true, "zen": true, "zgh": true, "zha": true, "zho": true, "znd": true, "zul": true, "zun": true,
	"zxx": true, "zza": true,
}

// IsLanguageAlpha2 validates ISO 639-1 language code (2-letter)
func IsLanguageAlpha2(str string) bool {
	return languageAlpha2[strings.ToLower(str)]
}

// IsLanguageAlpha3 validates ISO 639-2 language code (3-letter)
func IsLanguageAlpha3(str string) bool {
	return languageAlpha3[strings.ToLower(str)]
}

// IsLanguageCode validates any ISO 639 language code format
func IsLanguageCode(str string) bool {
	return IsLanguageAlpha2(str) || IsLanguageAlpha3(str)
}
