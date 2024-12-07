package locations

import (
	"github.com/grokify/mogo/math/mathutil"
	"google.golang.org/genproto/googleapis/type/latlng"
)

var (
	// https://aws.amazon.com/about-aws/global-infrastructure/regions_az/
	// https://github.com/jsonmaur/aws-regions/issues/11

	// https://stackoverflow.com/questions/13785466/default-center-on-united-states
	// USContiguousCenter is defined by U.S. National Geodetic Survey (NGS): https://en.wikipedia.org/wiki/Geographic_center_of_the_United_States
	USCenterContiguous = latlng.LatLng{
		Latitude:  mathutil.DegreesMinutesSecondsToDecimal(39, 50, 0),
		Longitude: mathutil.DegreesMinutesSecondsToDecimal(-98, 35, 0)}
	// USHIAKCenter is defined by U.S. National Geodetic Survey (NGS): https://en.wikipedia.org/wiki/Geographic_center_of_the_United_States
	USCenterAKHI = latlng.LatLng{
		Latitude:  mathutil.DegreesMinutesSecondsToDecimal(44, 58, 2.07622),
		Longitude: mathutil.DegreesMinutesSecondsToDecimal(-103, 46, 17.60283)}
	USHIEUCenter = latlng.LatLng{Latitude: 37.09024, Longitude: -70}
	// EU center is defined by French Institut Géographique National (IGN): https://en.wikipedia.org/wiki/Geographical_midpoint_of_Europe
	EUCenter2020 = latlng.LatLng{
		Latitude:  mathutil.DegreesMinutesSecondsToDecimal(49, 40, 34.8),
		Longitude: mathutil.DegreesMinutesSecondsToDecimal(9, 54, 7.4)}
	// Centers of Europe
	EuropeCenterHungary = latlng.LatLng{Latitude: 48.2352, Longitude: 21.22617} // HUTAY Tállya, Hungary (48°14' 21°14'), Europe Center 1992
	EuropeCenterBelarus = latlng.LatLng{
		Latitude:  mathutil.DegreesMinutesSecondsToDecimal(55, 30, 0),
		Longitude: mathutil.DegreesMinutesSecondsToDecimal(28, 48, 0)} // Polotsk, Belarus

	CACAL = latlng.LatLng{Latitude: 51.0447, Longitude: -114.0719} // Calgary, Calgary
	CAMTR = latlng.LatLng{Latitude: 45.5019, Longitude: -73.5674}  // Montreal, Quebec
	CATOR = latlng.LatLng{Latitude: 43.6532, Longitude: -79.3832}  // Toronto, Ontario
	CAVAN = latlng.LatLng{Latitude: 49.2827, Longitude: -123.1207} // Vancouver, BC

	USQAS = latlng.LatLng{Latitude: 39.0403, Longitude: -77.4852}       // Ashuburn, Virginia
	USATL = latlng.LatLng{Latitude: 33.7488, Longitude: -84.3877}       // Atlanta, Georgia
	USBOS = latlng.LatLng{Latitude: 42.3601, Longitude: -71.0589}       // Boston, MA
	USCHI = latlng.LatLng{Latitude: 41.8781, Longitude: -87.6298}       // Chicago, IL; CHCHIpop = 2714856
	USCMH = latlng.LatLng{Latitude: 39.9612, Longitude: -82.9988}       // Columbus, OH
	USDAL = latlng.LatLng{Latitude: 32.7767, Longitude: -96.7970}       // Dallas, TX
	USDEN = latlng.LatLng{Latitude: 39.7392, Longitude: -104.9903}      // Denver, CO
	USDET = latlng.LatLng{Latitude: 42.3314, Longitude: -83.0458}       // Detroit
	USDUO = latlng.LatLng{Latitude: 40.1055, Longitude: -83.1375}       // Dublin, Ohio
	USHNL = latlng.LatLng{Latitude: 21.3099, Longitude: -157.8581}      // Honolulu
	USHOU = latlng.LatLng{Latitude: 29.7604, Longitude: -95.3698}       // Houston, TX
	USJAX = latlng.LatLng{Latitude: 30.3322, Longitude: -81.6557}       // Jacksonville, FL
	USMKC = latlng.LatLng{Latitude: 39.11924, Longitude: -94.57237}     // Kansas City, MO - https://www.marinetraffic.com/en/ais/details/ports/21810
	US4NB = latlng.LatLng{Latitude: 39.8097343, Longitude: -98.5556199} // Lebanon, Kansas (US Center, Continental)
	USLAX = latlng.LatLng{Latitude: 34.0522, Longitude: -118.2437}      // Los Angeles, CA
	USMES = latlng.LatLng{Latitude: 44.9778, Longitude: -93.2650}       // Minneapolis, MN
	USNYC = latlng.LatLng{Latitude: 40.7128, Longitude: -74.0060}       // New York, NY; USNYCpop = 8405837
	USMIA = latlng.LatLng{Latitude: 25.761681, Longitude: -80.191788}   // Miami, Florida
	USPHX = latlng.LatLng{Latitude: 33.4484, Longitude: -112.0740}      // Phoenix, AZ
	USPDX = latlng.LatLng{Latitude: 45.5152, Longitude: -122.6784}      // Portland, Oregon
	USSTL = latlng.LatLng{Latitude: 38.6270, Longitude: -90.1994}       // St. Louis, MO
	USSFO = latlng.LatLng{Latitude: 37.77395, Longitude: -122.3963}     // San Francisco, California
	USSJC = latlng.LatLng{Latitude: 37.3387, Longitude: -121.8853}      // San Jose, California
	USSEA = latlng.LatLng{Latitude: 47.6062, Longitude: -122.3321}      // Sesattle, WA

	COBOG = latlng.LatLng{Latitude: 4.7110, Longitude: -74.0721}   // Bogota
	ARBUE = latlng.LatLng{Latitude: -34.6037, Longitude: -58.3816} // Buenos Aires
	BRCWB = latlng.LatLng{Latitude: -25.4372, Longitude: -49.2700} // Curitiba
	BRFOR = latlng.LatLng{Latitude: -3.7327, Longitude: -38.5270}  // Fortaleza
	PELIM = latlng.LatLng{Latitude: -12.0464, Longitude: -77.0428} // Lima
	BRSAO = latlng.LatLng{Latitude: -23.5558, Longitude: -46.6396} // São Paulo
	CLSCL = latlng.LatLng{Latitude: -33.4489, Longitude: -70.6693} // Santiago
	BRRIO = latlng.LatLng{Latitude: -22.9068, Longitude: -43.1729} // Rio de Janeiro

	NLAMS = latlng.LatLng{Latitude: 52.3676, Longitude: 4.9041}  // Amsterdam
	BEBRU = latlng.LatLng{Latitude: 50.8476, Longitude: 4.3572}  // Brussels
	DKCPH = latlng.LatLng{Latitude: 55.6761, Longitude: 12.5683} // Copenhagen
	IEDUB = latlng.LatLng{Latitude: 53.3498, Longitude: -6.2603} // Dublin
	DEFRA = latlng.LatLng{Latitude: 50.1109, Longitude: 8.6821}  // Frankfurt
	FIHEL = latlng.LatLng{Latitude: 60.1699, Longitude: 24.9384} // Helsinki
	PTLIS = latlng.LatLng{Latitude: 38.7223, Longitude: -9.1393} // Lisbon
	GBLON = latlng.LatLng{Latitude: 51.5072, Longitude: -0.1276} // London
	ESMAD = latlng.LatLng{Latitude: 40.4168, Longitude: -3.7038} // Madrid
	GBMAN = latlng.LatLng{Latitude: 53.4808, Longitude: -2.2426} // Manchester
	FRMRS = latlng.LatLng{Latitude: 43.2965, Longitude: 5.3698}  // Marseille
	ITMIL = latlng.LatLng{Latitude: 45.4642, Longitude: 9.1900}  // Milan
	DEMUC = latlng.LatLng{Latitude: 48.1351, Longitude: 11.5820} // Munich
	NOOSL = latlng.LatLng{Latitude: 59.9139, Longitude: 10.7522} // Oslo
	ITPMO = latlng.LatLng{Latitude: 38.1157, Longitude: 13.3615} // Palermo
	FRPAR = latlng.LatLng{Latitude: 48.8566, Longitude: 2.3522}  // Paris
	CZPRG = latlng.LatLng{Latitude: 50.0755, Longitude: 14.4378} // Prague - EU Center
	ITROM = latlng.LatLng{Latitude: 41.9028, Longitude: 12.4964} // Rome
	SESTO = latlng.LatLng{Latitude: 59.3293, Longitude: 18.0686} // Stockholm
	BGSOF = latlng.LatLng{Latitude: 42.6977, Longitude: 23.3219} // Sofia
	ATVIE = latlng.LatLng{Latitude: 48.2082, Longitude: 16.3738} // Vienna
	CHZRH = latlng.LatLng{Latitude: 47.3769, Longitude: 8.5417}  // Zurich

	THBKK = latlng.LatLng{Latitude: 13.7563, Longitude: 100.5018} // Bangkok
	CNBJS = latlng.LatLng{Latitude: 39.9042, Longitude: 116.4074} // Beijing
	INMAA = latlng.LatLng{Latitude: 13.0827, Longitude: 80.2707}  // Chennai
	AEDXB = latlng.LatLng{Latitude: 25.2048, Longitude: 55.2708}  // Dubai
	AEFJR = latlng.LatLng{Latitude: 25.1288, Longitude: 56.3265}  // Fujairah
	HKHKG = latlng.LatLng{Latitude: 22.3193, Longitude: 114.1694} // Hong Kong
	IDJKT = latlng.LatLng{Latitude: -6.2088, Longitude: 106.8456} // Jakarta
	INHYD = latlng.LatLng{Latitude: 17.3850, Longitude: 78.4867}  // Hyderabad
	INCCU = latlng.LatLng{Latitude: 22.5726, Longitude: 88.3639}  // Kolkata
	MYKUL = latlng.LatLng{Latitude: 3.1569, Longitude: 101.7123}  // Kuala Lumpur
	BHAMH = latlng.LatLng{Latitude: 26.2235, Longitude: 50.5876}  // Manama, Bahrain
	PHMNL = latlng.LatLng{Latitude: 14.5995, Longitude: 120.9842} // Manila
	INBOM = latlng.LatLng{Latitude: 19.0760, Longitude: 72.8777}  // Mumbai
	INICD = latlng.LatLng{Latitude: 28.6139, Longitude: 77.2090}  // New Delhi
	JPOSA = latlng.LatLng{Latitude: 34.6937, Longitude: 135.5023} // Osaka
	KRSEL = latlng.LatLng{Latitude: 37.5665, Longitude: 126.9780} // Seoul
	SGSIN = latlng.LatLng{Latitude: 1.3521, Longitude: 103.8198}  // Singapore
	JPTYO = latlng.LatLng{Latitude: 35.6762, Longitude: 139.6503} // Tokyo
	CNYCH = latlng.LatLng{Latitude: 38.4863, Longitude: 106.2324} // Yinchuan, Ningxia

	GHACC = latlng.LatLng{Latitude: 5.6037, Longitude: -0.1870}   // Accra
	ZACPT = latlng.LatLng{Latitude: -33.9249, Longitude: 18.4241} // Cape Town
	ZAJNB = latlng.LatLng{Latitude: -26.2041, Longitude: 28.0473} // Johannesburg

	AUADL = latlng.LatLng{Latitude: -34.9285, Longitude: 138.6007} // Adelaide
	NZAKL = latlng.LatLng{Latitude: -36.8509, Longitude: 174.7645} // Auckland
	AUBNE = latlng.LatLng{Latitude: -27.4705, Longitude: 153.0260} // Brisbane
	NZCHC = latlng.LatLng{Latitude: -43.5320, Longitude: 172.6306} // Christchurch
	AUMEL = latlng.LatLng{Latitude: -37.8136, Longitude: 144.9631} // Melbourne
	AUPER = latlng.LatLng{Latitude: -31.9523, Longitude: 115.8613} // Perth
	AUSYD = latlng.LatLng{Latitude: -33.8688, Longitude: 151.2093} // Sydney
	NZWLG = latlng.LatLng{Latitude: -41.2924, Longitude: 174.7787} // Wellington
)
