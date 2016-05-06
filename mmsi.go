// Copyright (c) 2015, Marios Andreopoulos.
//
// This file is part of aislib.
//
//  Aislib is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
//  Aislib is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
// along with aislib.  If not, see <http://www.gnu.org/licenses/>.

package aislib

// Contains MMSI owners' descriptions. Currently not used anywhere
var MmsiCodes = [...]string{
	"Ship", "Coastal Station", "Group of ships", "SAR —Search and Rescue Aircraft",
	"Diver's radio", "Aids to navigation", "Auxiliary craft associated with parent ship",
	"AIS SART —Search and Rescue Transmitter", "MOB —Man Overboard Device",
	"EPIRB —Emergency Position Indicating Radio Beacon", "Invalid MMSI",
}

// DecodeMMSI returns a string with the type of the owner of the MMSI and its country
// Some MMSIs aren't valid. There is some more information in some MMSIs (the satellite
// equipment of the ship). We may add them in the future.
// Have a look at http://en.wikipedia.org/wiki/Maritime_Mobile_Service_Identity
func DecodeMMSI(m uint32) string {
	owner := ""
	country := ""
	mid := uint32(1000)

	// Current intervals:
	// [0 00999999][010000000 099999999][100000000 199999999][200000000 799999999]
	// [800000000 899999999]...[970000000 970999999]...[972000000 972999999]...
	// [974000000 974999999]...[98000000 98999999][99000000 999999999]
	switch {
	case m >= 200000000 && m < 800000000:
		mid = m / 1000000
		owner = "Ship"
	case m <= 9999999:
		mid = m / 10000
		owner = "Coastal Station"
	case m <= 99999999:
		mid = m / 100000
		owner = "Group of ships"
	case m <= 199999999:
		mid = m/1000 - 111000
		owner = "SAR —Search and Rescue Aircraft"
	case m < 900000000:
		mid = m/100000 - 8000
		owner = "Diver's radio"
	case m >= 990000000 && m < 1000000000:
		mid = m/10000 - 99000
		owner = "Aids to navigation"
	case m >= 980000000 && m < 990000000:
		mid = m/10000 - 98000
		owner = "Auxiliary craft associated with parent ship"
	case m >= 970000000 && m < 970999999:
		mid = m/1000 - 970000
		owner = "AIS SART —Search and Rescue Transmitter"
	case m >= 972000000 && m < 972999999:
		owner = "MOB —Man Overboard Device"
	case m >= 974000000 && m < 974999999:
		owner = "EPIRB —Emergency Position Indicating Radio Beacon"
	default:
		owner = "Invalid MMSI"
	}

	if mid < 1000 {
		country = Mid[int(mid)]
		if country == "" {
			country = "Unknown Country ID"
		}
		return owner + ", " + country
	}
	return owner
}

// Maritime Identification Digits, have a look at http://www.itu.int/online/mms/glad/cga_mids.sh?lang=en
var Mid = map[int]string{
	201: "Albania (Republic of)",
	202: "Andorra (Principality of)",
	203: "Austria",
	204: "Azores - Portugal",
	205: "Belgium",
	206: "Belarus (Republic of)",
	207: "Bulgaria (Republic of)",
	208: "Vatican City State",
	209: "Cyprus (Republic of)",
	210: "Cyprus (Republic of)",
	211: "Germany (Federal Republic of)",
	212: "Cyprus (Republic of)",
	213: "Georgia",
	214: "Moldova (Republic of)",
	215: "Malta",
	216: "Armenia (Republic of)",
	218: "Germany (Federal Republic of)",
	219: "Denmark",
	220: "Denmark",
	224: "Spain",
	225: "Spain",
	226: "France",
	227: "France",
	228: "France",
	229: "Malta",
	230: "Finland",
	231: "Faroe Islands - Denmark",
	232: "United Kingdom of Great Britain and Northern Ireland",
	233: "United Kingdom of Great Britain and Northern Ireland",
	234: "United Kingdom of Great Britain and Northern Ireland",
	235: "United Kingdom of Great Britain and Northern Ireland",
	236: "Gibraltar - United Kingdom of Great Britain and Northern Ireland",
	237: "Greece",
	238: "Croatia (Republic of)",
	239: "Greece",
	240: "Greece",
	241: "Greece",
	242: "Morocco (Kingdom of)",
	243: "Hungary",
	244: "Netherlands (Kingdom of the)",
	245: "Netherlands (Kingdom of the)",
	246: "Netherlands (Kingdom of the)",
	247: "Italy",
	248: "Malta",
	249: "Malta",
	250: "Ireland",
	251: "Iceland",
	252: "Liechtenstein (Principality of)",
	253: "Luxembourg",
	254: "Monaco (Principality of)",
	255: "Madeira - Portugal",
	256: "Malta",
	257: "Norway",
	258: "Norway",
	259: "Norway",
	261: "Poland (Republic of)",
	262: "Montenegro",
	263: "Portugal",
	264: "Romania",
	265: "Sweden",
	266: "Sweden",
	267: "Slovak Republic",
	268: "San Marino (Republic of)",
	269: "Switzerland (Confederation of)",
	270: "Czech Republic",
	271: "Turkey",
	272: "Ukraine",
	273: "Russian Federation",
	274: "The Former Yugoslav Republic of Macedonia",
	275: "Latvia (Republic of)",
	276: "Estonia (Republic of)",
	277: "Lithuania (Republic of)",
	278: "Slovenia (Republic of)",
	279: "Serbia (Republic of)",
	301: "Anguilla - United Kingdom of Great Britain and Northern Ireland",
	303: "Alaska (State of) - United States of America",
	304: "Antigua and Barbuda",
	305: "Antigua and Barbuda",
	306: "Curacao, Sint Maarten (Dutch part), Bonaire, Sint Eustatius and Saba - Netherlands (Kingdom of the)",
	307: "Aruba - Netherlands (Kingdom of the)",
	308: "Bahamas (Commonwealth of the)",
	309: "Bahamas (Commonwealth of the)",
	310: "Bermuda - United Kingdom of Great Britain and Northern Ireland",
	311: "Bahamas (Commonwealth of the)",
	312: "Belize",
	314: "Barbados",
	316: "Canada",
	319: "Cayman Islands - United Kingdom of Great Britain and Northern Ireland",
	321: "Costa Rica",
	323: "Cuba",
	325: "Dominica (Commonwealth of)",
	327: "Dominican Republic",
	329: "Guadeloupe (French Department of) - France",
	330: "Grenada",
	331: "Greenland - Denmark",
	332: "Guatemala (Republic of)",
	334: "Honduras (Republic of)",
	336: "Haiti (Republic of)",
	338: "United States of America",
	339: "Jamaica",
	341: "Saint Kitts and Nevis (Federation of)",
	343: "Saint Lucia",
	345: "Mexico",
	347: "Martinique (French Department of) - France",
	348: "Montserrat - United Kingdom of Great Britain and Northern Ireland",
	350: "Nicaragua",
	351: "Panama (Republic of)",
	352: "Panama (Republic of)",
	353: "Panama (Republic of)",
	354: "Panama (Republic of)",
	355: " - ",
	356: " - ",
	357: " - ",
	358: "Puerto Rico - United States of America",
	359: "El Salvador (Republic of)",
	361: "Saint Pierre and Miquelon (Territorial Collectivity of) - France",
	362: "Trinidad and Tobago",
	364: "Turks and Caicos Islands - United Kingdom of Great Britain and Northern Ireland",
	366: "United States of America",
	367: "United States of America",
	368: "United States of America",
	369: "United States of America",
	370: "Panama (Republic of)",
	371: "Panama (Republic of)",
	372: "Panama (Republic of)",
	373: "Panama (Republic of)",
	375: "Saint Vincent and the Grenadines",
	376: "Saint Vincent and the Grenadines",
	377: "Saint Vincent and the Grenadines",
	378: "British Virgin Islands - United Kingdom of Great Britain and Northern Ireland",
	379: "United States Virgin Islands - United States of America",
	401: "Afghanistan",
	403: "Saudi Arabia (Kingdom of)",
	405: "Bangladesh (People's Republic of)",
	408: "Bahrain (Kingdom of)",
	410: "Bhutan (Kingdom of)",
	412: "China (People's Republic of)",
	413: "China (People's Republic of)",
	414: "China (People's Republic of)",
	416: "Taiwan (Province of China) - China (People's Republic of)",
	417: "Sri Lanka (Democratic Socialist Republic of)",
	419: "India (Republic of)",
	422: "Iran (Islamic Republic of)",
	423: "Azerbaijan (Republic of)",
	425: "Iraq (Republic of)",
	428: "Israel (State of)",
	431: "Japan",
	432: "Japan",
	434: "Turkmenistan",
	436: "Kazakhstan (Republic of)",
	437: "Uzbekistan (Republic of)",
	438: "Jordan (Hashemite Kingdom of)",
	440: "Korea (Republic of)",
	441: "Korea (Republic of)",
	443: "State of Palestine (In accordance with Resolution 99 Rev. Guadalajara, 2010)",
	445: "Democratic People's Republic of Korea",
	447: "Kuwait (State of)",
	450: "Lebanon",
	451: "Kyrgyz Republic",
	453: "Macao (Special Administrative Region of China) - China (People's Republic of)",
	455: "Maldives (Republic of)",
	457: "Mongolia",
	459: "Nepal (Federal Democratic Republic of)",
	461: "Oman (Sultanate of)",
	463: "Pakistan (Islamic Republic of)",
	466: "Qatar (State of)",
	468: "Syrian Arab Republic",
	470: "United Arab Emirates",
	472: "Tajikistan (Republic of)",
	473: "Yemen (Republic of)",
	475: "Yemen (Republic of)",
	477: "Hong Kong (Special Administrative Region of China) - China (People's Republic of)",
	478: "Bosnia and Herzegovina",
	501: "Adelie Land - France",
	503: "Australia",
	506: "Myanmar (Union of)",
	508: "Brunei Darussalam",
	510: "Micronesia (Federated States of)",
	511: "Palau (Republic of)",
	512: "New Zealand",
	514: "Cambodia (Kingdom of)",
	515: "Cambodia (Kingdom of)",
	516: "Christmas Island (Indian Ocean) - Australia",
	518: "Cook Islands - New Zealand",
	520: "Fiji (Republic of)",
	523: "Cocos (Keeling) Islands - Australia",
	525: "Indonesia (Republic of)",
	529: "Kiribati (Republic of)",
	531: "Lao People's Democratic Republic",
	533: "Malaysia",
	536: "Northern Mariana Islands (Commonwealth of the) - United States of America",
	538: "Marshall Islands (Republic of the)",
	540: "New Caledonia - France",
	542: "Niue - New Zealand",
	544: "Nauru (Republic of)",
	546: "French Polynesia - France",
	548: "Philippines (Republic of the)",
	553: "Papua New Guinea",
	555: "Pitcairn Island - United Kingdom of Great Britain and Northern Ireland",
	557: "Solomon Islands",
	559: "American Samoa - United States of America",
	561: "Samoa (Independent State of)",
	563: "Singapore (Republic of)",
	564: "Singapore (Republic of)",
	565: "Singapore (Republic of)",
	566: "Singapore (Republic of)",
	567: "Thailand",
	570: "Tonga (Kingdom of)",
	572: "Tuvalu",
	574: "Viet Nam (Socialist Republic of)",
	576: "Vanuatu (Republic of)",
	577: "Vanuatu (Republic of)",
	578: "Wallis and Futuna Islands - France",
	601: "South Africa (Republic of)",
	603: "Angola (Republic of)",
	605: "Algeria (People's Democratic Republic of)",
	607: "Saint Paul and Amsterdam Islands - France",
	608: "Ascension Island - United Kingdom of Great Britain and Northern Ireland",
	609: "Burundi (Republic of)",
	610: "Benin (Republic of)",
	611: "Botswana (Republic of)",
	612: "Central African Republic",
	613: "Cameroon (Republic of)",
	615: "Congo (Republic of the)",
	616: "Comoros (Union of the)",
	617: "Cabo Verde (Republic of)",
	618: "Crozet Archipelago - France",
	619: "Cote d'Ivoire (Republic of)",
	620: "Comoros (Union of the)",
	621: "Djibouti (Republic of)",
	622: "Egypt (Arab Republic of)",
	624: "Ethiopia (Federal Democratic Republic of)",
	625: "Eritrea",
	626: "Gabonese Republic",
	627: "Ghana",
	629: "Gambia (Republic of the)",
	630: "Guinea-Bissau (Republic of)",
	631: "Equatorial Guinea (Republic of)",
	632: "Guinea (Republic of)",
	633: "Burkina Faso",
	634: "Kenya (Republic of)",
	635: "Kerguelen Islands - France",
	636: "Liberia (Republic of)",
	637: "Liberia (Republic of)",
	638: "South Sudan (Republic of)",
	642: "Libya",
	644: "Lesotho (Kingdom of)",
	645: "Mauritius (Republic of)",
	647: "Madagascar (Republic of)",
	649: "Mali (Republic of)",
	650: "Mozambique (Republic of)",
	654: "Mauritania (Islamic Republic of)",
	655: "Malawi",
	656: "Niger (Republic of the)",
	657: "Nigeria (Federal Republic of)",
	659: "Namibia (Republic of)",
	660: "Reunion (French Department of) - France",
	661: "Rwanda (Republic of)",
	662: "Sudan (Republic of the)",
	663: "Senegal (Republic of)",
	664: "Seychelles (Republic of)",
	665: "Saint Helena - United Kingdom of Great Britain and Northern Ireland",
	666: "Somalia (Federal Republic of)",
	667: "Sierra Leone",
	668: "Sao Tome and Principe (Democratic Republic of)",
	669: "Swaziland (Kingdom of)",
	670: "Chad (Republic of)",
	671: "Togolese Republic",
	672: "Tunisia",
	674: "Tanzania (United Republic of)",
	675: "Uganda (Republic of)",
	676: "Democratic Republic of the Congo",
	677: "Tanzania (United Republic of)",
	678: "Zambia (Republic of)",
	679: "Zimbabwe (Republic of)",
	701: "Argentine Republic",
	710: "Brazil (Federative Republic of)",
	720: "Bolivia (Plurinational State of)",
	725: "Chile",
	730: "Colombia (Republic of)",
	735: "Ecuador",
	740: "Falkland Islands (Malvinas) - United Kingdom of Great Britain and Northern Ireland",
	745: "Guiana (French Department of) - France",
	750: "Guyana",
	755: "Paraguay (Republic of)",
	760: "Peru",
	765: "Suriname (Republic of)",
	770: "Uruguay (Eastern Republic of)",
}
