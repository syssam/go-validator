package validator

import "strings"

// ISO 4217 currency codes
var currencyCodes = map[string]bool{
	"AED": true, "AFN": true, "ALL": true, "AMD": true, "ANG": true, "AOA": true, "ARS": true, "AUD": true,
	"AWG": true, "AZN": true, "BAM": true, "BBD": true, "BDT": true, "BGN": true, "BHD": true, "BIF": true,
	"BMD": true, "BND": true, "BOB": true, "BOV": true, "BRL": true, "BSD": true, "BTN": true, "BWP": true,
	"BYN": true, "BZD": true, "CAD": true, "CDF": true, "CHE": true, "CHF": true, "CHW": true, "CLF": true,
	"CLP": true, "CNY": true, "COP": true, "COU": true, "CRC": true, "CUC": true, "CUP": true, "CVE": true,
	"CZK": true, "DJF": true, "DKK": true, "DOP": true, "DZD": true, "EGP": true, "ERN": true, "ETB": true,
	"EUR": true, "FJD": true, "FKP": true, "GBP": true, "GEL": true, "GHS": true, "GIP": true, "GMD": true,
	"GNF": true, "GTQ": true, "GYD": true, "HKD": true, "HNL": true, "HRK": true, "HTG": true, "HUF": true,
	"IDR": true, "ILS": true, "INR": true, "IQD": true, "IRR": true, "ISK": true, "JMD": true, "JOD": true,
	"JPY": true, "KES": true, "KGS": true, "KHR": true, "KMF": true, "KPW": true, "KRW": true, "KWD": true,
	"KYD": true, "KZT": true, "LAK": true, "LBP": true, "LKR": true, "LRD": true, "LSL": true, "LYD": true,
	"MAD": true, "MDL": true, "MGA": true, "MKD": true, "MMK": true, "MNT": true, "MOP": true, "MRU": true,
	"MUR": true, "MVR": true, "MWK": true, "MXN": true, "MXV": true, "MYR": true, "MZN": true, "NAD": true,
	"NGN": true, "NIO": true, "NOK": true, "NPR": true, "NZD": true, "OMR": true, "PAB": true, "PEN": true,
	"PGK": true, "PHP": true, "PKR": true, "PLN": true, "PYG": true, "QAR": true, "RON": true, "RSD": true,
	"RUB": true, "RWF": true, "SAR": true, "SBD": true, "SCR": true, "SDG": true, "SEK": true, "SGD": true,
	"SHP": true, "SLE": true, "SLL": true, "SOS": true, "SRD": true, "SSP": true, "STN": true, "SVC": true,
	"SYP": true, "SZL": true, "THB": true, "TJS": true, "TMT": true, "TND": true, "TOP": true, "TRY": true,
	"TTD": true, "TWD": true, "TZS": true, "UAH": true, "UGX": true, "USD": true, "USN": true, "UYI": true,
	"UYU": true, "UYW": true, "UZS": true, "VED": true, "VES": true, "VND": true, "VUV": true, "WST": true,
	"XAF": true, "XAG": true, "XAU": true, "XBA": true, "XBB": true, "XBC": true, "XBD": true, "XCD": true,
	"XDR": true, "XOF": true, "XPD": true, "XPF": true, "XPT": true, "XSU": true, "XTS": true, "XUA": true,
	"XXX": true, "YER": true, "ZAR": true, "ZMW": true, "ZWL": true,
}

// Cryptocurrency codes (common cryptocurrencies)
var currencyCrypto = map[string]bool{
	// Major cryptocurrencies
	"BTC":  true, // Bitcoin
	"ETH":  true, // Ethereum
	"USDT": true, // Tether
	"USDC": true, // USD Coin
	"BNB":  true, // Binance Coin
	"XRP":  true, // Ripple
	"ADA":  true, // Cardano
	"DOGE": true, // Dogecoin
	"SOL":  true, // Solana
	"DOT":  true, // Polkadot
	"MATIC": true, // Polygon
	"LTC":  true, // Litecoin
	"SHIB": true, // Shiba Inu
	"TRX":  true, // Tron
	"AVAX": true, // Avalanche
	"LINK": true, // Chainlink
	"ATOM": true, // Cosmos
	"XMR":  true, // Monero
	"ETC":  true, // Ethereum Classic
	"BCH":  true, // Bitcoin Cash
	"XLM":  true, // Stellar
	"ALGO": true, // Algorand
	"VET":  true, // VeChain
	"FIL":  true, // Filecoin
	"ICP":  true, // Internet Computer
	"NEAR": true, // NEAR Protocol
	"APT":  true, // Aptos
	"ARB":  true, // Arbitrum
	"OP":   true, // Optimism
	"AAVE": true, // Aave
	"UNI":  true, // Uniswap
	"MKR":  true, // Maker
	"CRO":  true, // Cronos
	"QNT":  true, // Quant
	"GRT":  true, // The Graph
	"FTM":  true, // Fantom
	"SAND": true, // The Sandbox
	"MANA": true, // Decentraland
	"AXS":  true, // Axie Infinity
	"THETA": true, // Theta Network
	"EGLD": true, // MultiversX
	"EOS":  true, // EOS
	"XTZ":  true, // Tezos
	"FLOW": true, // Flow
	"CHZ":  true, // Chiliz
	"CAKE": true, // PancakeSwap
	"ZEC":  true, // Zcash
	"NEO":  true, // Neo
	"KAVA": true, // Kava
	"DASH": true, // Dash
	"WAVES": true, // Waves
	"MINA": true, // Mina Protocol
	"ZIL":  true, // Zilliqa
	"ENJ":  true, // Enjin Coin
	"BAT":  true, // Basic Attention Token
	"1INCH": true, // 1inch
	"COMP": true, // Compound
	"SNX":  true, // Synthetix
	"YFI":  true, // yearn.finance
	"SUSHI": true, // SushiSwap
	"CRV":  true, // Curve DAO Token
	"LDO":  true, // Lido DAO
	"RPL":  true, // Rocket Pool
	"RUNE": true, // THORChain
	"INJ":  true, // Injective
	"SUI":  true, // Sui
	"SEI":  true, // Sei
	"TIA":  true, // Celestia
	"JUP":  true, // Jupiter
	"PYTH": true, // Pyth Network
	"WLD":  true, // Worldcoin
	"BLUR": true, // Blur
	"PEPE": true, // Pepe
	"BONK": true, // Bonk
	"WIF":  true, // dogwifhat
	"FLOKI": true, // Floki
	"FET":  true, // Fetch.ai
	"RNDR": true, // Render
	"AGIX": true, // SingularityNET
	"OCEAN": true, // Ocean Protocol
	"TAO":  true, // Bittensor
	"AR":   true, // Arweave
	"STX":  true, // Stacks
	"IMX":  true, // Immutable X
	"GMX":  true, // GMX
	"DYDX": true, // dYdX
	"OSMO": true, // Osmosis
	"CFX":  true, // Conflux
	"ROSE": true, // Oasis Network
	"CELO": true, // Celo
	"KDA":  true, // Kadena
	"HBAR": true, // Hedera
	"IOTA": true, // IOTA
	"XDC":  true, // XDC Network
	"KLAY": true, // Klaytn
	"ONE":  true, // Harmony
	"ICX":  true, // ICON
	"ZRX":  true, // 0x Protocol
	"LRC":  true, // Loopring
	"GALA": true, // Gala
	"APE":  true, // ApeCoin
	"MASK": true, // Mask Network
	"ENS":  true, // Ethereum Name Service
	"SSV":  true, // SSV Network
	"PENDLE": true, // Pendle
	"STRK": true, // Starknet
}

// IsCurrencyFiat validates ISO 4217 fiat currency code (3-letter)
func IsCurrencyFiat(str string) bool {
	return currencyCodes[strings.ToUpper(str)]
}

// IsCurrencyCrypto validates cryptocurrency code
func IsCurrencyCrypto(str string) bool {
	return currencyCrypto[strings.ToUpper(str)]
}

// IsCurrencyAll validates any currency (fiat + crypto)
func IsCurrencyAll(str string) bool {
	upper := strings.ToUpper(str)
	return currencyCodes[upper] || currencyCrypto[upper]
}
