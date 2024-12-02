package country

import "github.com/grokify/mogo/type/maputil"

type CountriesMap map[string]string

func (cm CountriesMap) Add(c CountriesMap) CountriesMap {
	for k, v := range c {
		cm[k] = v
	}
	return cm
}

func (cm CountriesMap) ISO3166P1Alpha2s() []string {
	return maputil.Keys(cm)
}

func (cm CountriesMap) Sub(c CountriesMap) CountriesMap {
	out := CountriesMap{}
	for k, v := range cm {
		if _, ok := c[k]; !ok {
			out[k] = v
		}
	}
	return out
}

func CountriesMapANZ() CountriesMap {
	return map[string]string{
		"AU": "Australia",
		"NZ": "New Zealand",
	}
}

func CountriesMapBalkan() CountriesMap {
	return map[string]string{
		"AL": "Albania",
		"BA": "Bosnia and Herzegovina",
		"LT": "Bulgaria",
		"HR": "Croatia",
		"GR": "Greece", // Kosovo
		"ME": "Montenegro",
		"MK": "North Macedonia",
		"RO": "Romania",
		"RS": "Serbia",
		"SI": "Slovenia",
	}
}

func CountriesMapBaltic() CountriesMap {
	return map[string]string{
		"EE": "Estonia",
		"LV": "Latvia",
		"LT": "Lithuania",
	}
}

func CountriesMapBenelux() CountriesMap {
	return map[string]string{
		"BE": "Belgium",
		"LU": "Luxembourg",
		"NL": "Netherlands",
	}
}

func CountriesMapCaucasus() CountriesMap {
	return map[string]string{
		"AM": "Armenia",
		"AZ": "Azerbaijan",
		"GE": "Georgia",
	}
}

func CountriesMapDACH() CountriesMap {
	return map[string]string{
		"AT": "Austria",
		"DE": "Germany",
		"CH": "Switzerland",
	}
}

func CountriesMapEurope() CountriesMap {
	return map[string]string{
		"AL": "Albania",
		"AD": "Andorra",
		"AM": "Armenia",
		"AT": "Austria",
		"BY": "Belarus",
		"BE": "Belgium",
		"BA": "Bosnia and Herzegovina",
		"BG": "Bulgaria",
		"HR": "Croatia",
		"CY": "Cyprus",
		"CZ": "Czechia",
		"DK": "Denmark",
		"EE": "Estonia",
		"FI": "Finland",
		"FR": "France",
		"GE": "Georgia",
		"DE": "Germany",
		"GR": "Greece",
		"HU": "Hungary",
		"IS": "Iceland",
		"IE": "Ireland",
		"IT": "Italy",
		"XK": "Kosovo",
		"LV": "Latvia",
		"LI": "Liechtenstein",
		"LT": "Lithuania",
		"LU": "Luxembourg",
		"MT": "Malta",
		"MD": "Moldova",
		"MC": "Monaco",
		"ME": "Montenegro",
		"NL": "Netherlands",
		"MK": "North Macedonia",
		"NO": "Norway",
		"PL": "Poland",
		"PT": "Portugal",
		"RO": "Romania",
		"RU": "Russia",
		"SM": "San Marino",
		"RS": "Serbia",
		"SK": "Slovakia",
		"SI": "Slovenia",
		"ES": "Spain",
		"SE": "Sweden",
		"CH": "Switzerland",
		"TR": "Turkey",
		"UA": "Ukraine",
		"GB": "United Kingdom",
		"VA": "Vatican City",
	}
}
func CountriesMapEasternEurope() CountriesMap {
	return map[string]string{
		"BY": "Belarus",
		"MD": "Moldova",
		"RO": "Romania",
		"RU": "Russia",
		"UA": "Ukraine",
	}
}

func CountriesMapEasternEuropeFull() CountriesMap {
	out := CountriesMap{}
	funcs := []func() CountriesMap{
		CountriesMapBalkan,
		CountriesMapBaltic,
		CountriesMapCaucasus,
		CountriesMapEasternEurope,
		CountriesMapVisegrad,
	}
	for _, fi := range funcs {
		out = out.Add(fi())
	}
	return out
}

func CountriesMapNearAndMiddleEast() CountriesMap {
	return map[string]string{
		"AF": "Afghanistan",
		"AM": "Armenia",
		"AZ": "Azerbaijan",
		"BH": "Bahrain",
		"CY": "Cyprus",
		"EG": "Egypt",
		"GE": "Georgia",
		"IR": "Iran",
		"IQ": "Iraq",
		"IL": "Israel",
		"JO": "Jordan",
		"KW": "Kuwait",
		"LB": "Lebanon",
		"OM": "Oman",
		"PS": "Palestine",
		"QA": "Qatar",
		"SA": "Saudi Arabia",
		"SY": "Syria",
		"TR": "Turkey",
		"AE": "United Arab Emirates",
		"YE": "Yemen",
	}
}

func CountriesMapNordic() CountriesMap {
	return map[string]string{
		"DK": "Denmark",
		"FI": "Finland",
		"IS": "Iceland",
		"NO": "Norway",
		"SE": "Sweden",
	}
}

func CountriesMapVisegrad() CountriesMap {
	return map[string]string{
		"CZ": "Czechia",
		"HU": "Hungary",
		"PL": "Poland",
		"SK": "Slovakia",
	}
}

func CountriesMapAll() CountriesMap {
	return map[string]string{
		"AF": "Afghanistan",
		"AL": "Albania",
		"DZ": "Algeria",
		"AS": "American Samoa",
		"AD": "Andorra",
		"AO": "Angola",
		"AI": "Anguilla",
		"AQ": "Antarctica",
		"AG": "Antigua and Barbuda",
		"AR": "Argentina",
		"AM": "Armenia",
		"AW": "Aruba",
		"AU": "Australia",
		"AT": "Austria",
		"AZ": "Azerbaijan",
		"BS": "Bahamas",
		"BH": "Bahrain",
		"BD": "Bangladesh",
		"BB": "Barbados",
		"BY": "Belarus",
		"BE": "Belgium",
		"BZ": "Belize",
		"BJ": "Benin",
		"BM": "Bermuda",
		"BT": "Bhutan",
		"BO": "Bolivia",
		"BA": "Bosnia and Herzegovina",
		"BW": "Botswana",
		"BR": "Brazil",
		"IO": "British Indian Ocean Territory",
		"BN": "Brunei",
		"BG": "Bulgaria",
		"BF": "Burkina Faso",
		"BI": "Burundi",
		"CV": "Cabo Verde",
		"KH": "Cambodia",
		"CM": "Cameroon",
		"CA": "Canada",
		"KY": "Cayman Islands",
		"CF": "Central African Republic",
		"TD": "Chad",
		"CL": "Chile",
		"CN": "China",
		"CO": "Colombia",
		"KM": "Comoros",
		"CG": "Congo (Brazzaville)",
		"CD": "Congo (Kinshasa)",
		"CR": "Costa Rica",
		"HR": "Croatia",
		"CU": "Cuba",
		"CY": "Cyprus",
		"CZ": "Czechia",
		"DK": "Denmark",
		"DJ": "Djibouti",
		"DM": "Dominica",
		"DO": "Dominican Republic",
		"EC": "Ecuador",
		"EG": "Egypt",
		"SV": "El Salvador",
		"GQ": "Equatorial Guinea",
		"ER": "Eritrea",
		"EE": "Estonia",
		"SZ": "Eswatini",
		"ET": "Ethiopia",
		"FJ": "Fiji",
		"FI": "Finland",
		"FR": "France",
		"GA": "Gabon",
		"GM": "Gambia",
		"GE": "Georgia",
		"DE": "Germany",
		"GH": "Ghana",
		"GR": "Greece",
		"GD": "Grenada",
		"GU": "Guam",
		"GT": "Guatemala",
		"GN": "Guinea",
		"GW": "Guinea-Bissau",
		"GY": "Guyana",
		"HT": "Haiti",
		"HN": "Honduras",
		"HU": "Hungary",
		"IS": "Iceland",
		"IN": "India",
		"ID": "Indonesia",
		"IR": "Iran",
		"IQ": "Iraq",
		"IE": "Ireland",
		"IL": "Israel",
		"IT": "Italy",
		"JM": "Jamaica",
		"JP": "Japan",
		"JO": "Jordan",
		"KZ": "Kazakhstan",
		"KE": "Kenya",
		"KI": "Kiribati",
		"KP": "Korea (North)",
		"KR": "Korea (South)",
		"KW": "Kuwait",
		"KG": "Kyrgyzstan",
		"LA": "Laos",
		"LV": "Latvia",
		"LB": "Lebanon",
		"LS": "Lesotho",
		"LR": "Liberia",
		"LY": "Libya",
		"LI": "Liechtenstein",
		"LT": "Lithuania",
		"LU": "Luxembourg",
		"MG": "Madagascar",
		"MW": "Malawi",
		"MY": "Malaysia",
		"MV": "Maldives",
		"ML": "Mali",
		"MT": "Malta",
		"MH": "Marshall Islands",
		"MR": "Mauritania",
		"MU": "Mauritius",
		"MX": "Mexico",
		"FM": "Micronesia",
		"MD": "Moldova",
		"MC": "Monaco",
		"MN": "Mongolia",
		"ME": "Montenegro",
		"MA": "Morocco",
		"MZ": "Mozambique",
		"MM": "Myanmar",
		"NA": "Namibia",
		"NR": "Nauru",
		"NP": "Nepal",
		"NL": "Netherlands",
		"NZ": "New Zealand",
		"NI": "Nicaragua",
		"NE": "Niger",
		"NG": "Nigeria",
		"MK": "North Macedonia",
		"NO": "Norway",
		"OM": "Oman",
		"PK": "Pakistan",
		"PW": "Palau",
		"PS": "Palestine",
		"PA": "Panama",
		"PG": "Papua New Guinea",
		"PY": "Paraguay",
		"PE": "Peru",
		"PH": "Philippines",
		"PL": "Poland",
		"PT": "Portugal",
		"QA": "Qatar",
		"RO": "Romania",
		"RU": "Russia",
		"RW": "Rwanda",
		"WS": "Samoa",
		"SM": "San Marino",
		"ST": "Sao Tome and Principe",
		"SA": "Saudi Arabia",
		"SN": "Senegal",
		"RS": "Serbia",
		"SC": "Seychelles",
		"SL": "Sierra Leone",
		"SG": "Singapore",
		"SK": "Slovakia",
		"SI": "Slovenia",
		"SB": "Solomon Islands",
		"SO": "Somalia",
		"ZA": "South Africa",
		"SS": "South Sudan",
		"ES": "Spain",
		"LK": "Sri Lanka",
		"SD": "Sudan",
		"SR": "Suriname",
		"SE": "Sweden",
		"CH": "Switzerland",
		"SY": "Syria",
		"TW": "Taiwan",
		"TJ": "Tajikistan",
		"TZ": "Tanzania",
		"TH": "Thailand",
		"TL": "Timor-Leste",
		"TG": "Togo",
		"TO": "Tonga",
		"TT": "Trinidad and Tobago",
		"TN": "Tunisia",
		"TR": "Turkey",
		"TM": "Turkmenistan",
		"TV": "Tuvalu",
		"UG": "Uganda",
		"UA": "Ukraine",
		"AE": "United Arab Emirates",
		"GB": "United Kingdom",
		"US": "United States",
		"UY": "Uruguay",
		"UZ": "Uzbekistan",
		"VU": "Vanuatu",
		"VA": "Vatican City",
		"VE": "Venezuela",
		"VN": "Vietnam",
		"YE": "Yemen",
		"ZM": "Zambia",
		"ZW": "Zimbabwe",
	}
}

func CountriesMapAPAC() CountriesMap {
	return map[string]string{
		"AF": "Afghanistan",
		"AM": "Armenia",
		"AZ": "Azerbaijan",
		"AU": "Australia",
		"BD": "Bangladesh",
		"BT": "Bhutan",
		"BN": "Brunei",
		"KH": "Cambodia",
		"CN": "China",
		"FJ": "Fiji",
		"GE": "Georgia",
		"IN": "India",
		"ID": "Indonesia",
		"JP": "Japan",
		"KZ": "Kazakhstan",
		"KR": "Korea (South)",
		"KP": "Korea (North)",
		"KG": "Kyrgyzstan",
		"LA": "Laos",
		"MY": "Malaysia",
		"MV": "Maldives",
		"MN": "Mongolia",
		"MM": "Myanmar",
		"NP": "Nepal",
		"NZ": "New Zealand",
		"PK": "Pakistan",
		"PW": "Palau",
		"PG": "Papua New Guinea",
		"PH": "Philippines",
		"SG": "Singapore",
		"SB": "Solomon Islands",
		"LK": "Sri Lanka",
		"TJ": "Tajikistan",
		"TH": "Thailand",
		"TL": "Timor-Leste",
		"TM": "Turkmenistan",
		"UZ": "Uzbekistan",
		"VN": "Vietnam",
		"VU": "Vanuatu",
	}
}

func CountriesMapAPACSales() CountriesMap {
	return map[string]string{
		"AU": "Australia",
		"BD": "Bangladesh",
		"BN": "Brunei",
		"KH": "Cambodia",
		"CN": "China",
		"FJ": "Fiji",
		"IN": "India",
		"ID": "Indonesia",
		"JP": "Japan",
		"KR": "Korea (South)",
		"MY": "Malaysia",
		"MV": "Maldives",
		"MN": "Mongolia",
		"MM": "Myanmar",
		"NP": "Nepal",
		"NZ": "New Zealand",
		"PK": "Pakistan",
		"PG": "Papua New Guinea",
		"PH": "Philippines",
		"SG": "Singapore",
		"LK": "Sri Lanka",
		"TH": "Thailand",
		"TL": "Timor-Leste",
		"VN": "Vietnam",
	}
}
