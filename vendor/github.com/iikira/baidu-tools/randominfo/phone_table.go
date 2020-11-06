package randominfo

var (
	// PhoneModelDataBase 手机型号库
	PhoneModelDataBase = []string{
		"1501_M02",        //360 F4
		"1503-M02",        //360 N4
		"1505-A01",        //360 N4S
		"303SH",           //夏普 Aquos Crystal Xx Mini 303SH
		"304SH",           //夏普 Aquos Crystal Xx SoftBank
		"305SH",           //夏普 Aquos Crystal Y
		"306SH",           //夏普 Aquos Crystal 306SH
		"360 Q5 Plus",     //360 Q5 Plus
		"360 Q5",          //360 Q5
		"402SH",           //夏普 Aquos Crystal X
		"502SH",           //夏普 Aquos Crystal Xx2
		"6607",            //OPPO U3
		"A1001",           //一加手机1
		"ASUS_A001",       //华硕 ZenFone 3 Ultra
		"ASUS_A001",       //华硕 ZenFone 3 Ultra
		"ASUS_Z00ADB",     //华硕 ZenFone 2
		"ASUS_Z00UDB",     //华硕 Zenfone Selfie
		"ASUS_Z00XSB",     //华硕 ZenFone Zoom
		"ASUS_Z012DE",     //华硕 ZenFone 3
		"ASUS_Z012DE",     //华硕 ZenFone 3
		"ASUS_Z016D",      //华硕 ZenFone 3 尊爵
		"ATH-TL00H",       //华为 荣耀 7i
		"Aster T",         //Vertu Aster T
		"BLN-AL10",        //华为 荣耀 畅玩6X
		"BND-AL10",        //荣耀7X
		"BTV-W09",         //华为 M3
		"CAM-UL00",        //华为 荣耀 畅玩5A
		"Constellation V", //Vertu Constellation V
		"D6683",           //索尼 Xperia Z3 Dual TD
		"DIG-AL00",        //华为 畅享 6S
		"E2312",           //索尼 Xperia M4 Aqua
		"E2363 ",          //索尼 Xperia M4 Aqua Dual
		"E5363",           //索尼 Xperia C4
		"E5563",           //索尼 Xperia C5
		"E5663",           //索尼 Xperia M5
		"E5823",           //索尼 Xperia Z5 Compact
		"E6533",           //索尼 Xperia Z3+
		"E6683",           //索尼 Xperia Z5
		"E6883",           //索尼 Xperia Z5 Premium
		"EBEN M2",         //8848 M2
		"EDI-AL10",        //华为 荣耀 Note 8
		"EVA-AL00",        //华为 P9
		"F100A",           //金立 F100
		"F103B",           //金立 F103B
		"F3116",           //索尼 Xperia XA
		"F3216",           //索尼 Xperia XA Ultra
		"F5121 / F5122",   //索尼 Xperia X
		"F5321",           //索尼 Xperia X Compact
		"F8132",           //索尼 Xperia X Performance
		"F8332",           //索尼 Xperia XZ
		"FRD-AL00",        //华为 荣耀 8
		"FS8001",          //夏普 C1
		"FS8002",          //夏普 A1
		"G0111",           //格力手机 1
		"G0215",           //格力手机 2
		"G8142",           //索尼Xperia XZ Premium G8142
		"G8342",           //索尼Xperia XZ1
		"GIONEE S9",       //金立 S9
		"GN5001S",         //金立 金钢
		"GN5003",          //金立 大金钢
		"GN8002S",         //金立 M6 Plus
		"GN8003",          //金立 M6
		"GN9011",          //金立 S8
		"GN9012",          //金立 S6 Pro
		"GRA-A0",          //Coolpad Cool Play 6C
		"H60-L11",         //华为 荣耀 6
		"HN3-U01",         //华为 荣耀 3
		"HTC D10w",        //HTC Desire 10 Pro
		"HTC E9pw",        //HTC One E9+
		"HTC M10u",        //HTC 10
		"HTC M8St",        //HTC One M8
		"HTC M9PT",        //HTC One M9+
		"HTC M9e",         //HTC One M9
		"HTC One A9",      //HTC One A9
		"HTC U-1w",        //HTC U Ultra
		"HTC X9u",         //HTC One X9
		"HTC_M10h",        //HTC 10 国际版
		"HUAWEI CAZ-AL00", //华为 Nova
		"HUAWEI CRR-UL00", //华为 Mate S
		"HUAWEI GRA-UL10", //华为 P8
		"HUAWEI MLA-AL10", //华为 麦芒 5
		"HUAWEI MT7-AL00", //华为 mate 7
		"HUAWEI MT7-TL00", //华为 Mate 7
		"HUAWEI NXT-AL10", //华为 Mate 8
		"HUAWEI P7-L00",   //华为 P7
		"HUAWEI RIO-AL00", //华为 麦芒 4
		"HUAWEI TAG-AL00", //华为 畅享 5S
		"HUAWEI VNS-AL00", //华为 G9
		"IUNI N1",         //艾优尼 N1
		"IUNI i1",         //艾优尼 i1
		"KFAPWI",          //Amazon Kindle Fire HDX 8.9
		"KFSOWI",          //Amazon Kindle Fire HDX 7
		"KFTHWI",          //Amazon Kindle Fire HD
		"KIW-TL00H",       //华为 荣耀 畅玩5X
		"KNT-AL10",        //华为 荣耀 V8
		"L55t",            //索尼 Xperia Z3
		"L55u",            //索尼 Xperia Z3
		"LEX626",          //乐视 乐S3
		"LEX720",          //乐视 乐Pro3
		"LG-D858",         //LG G3
		"LG-H818",         //LG G4
		"LG-H848",         //LG G5 SE
		"LG-H868",         //LG G5
		"LG-H968",         //LG V10
		"LON-AL00",        //华为 Mate 9 Pro
		"LON-AL00-PD",     //华为 Mate 9 Porsche Design
		"LT18i",           //Sony Ericsson Xperia Arc S
		"LT22i",           //Sony Ericsson Xperia P
		"LT26i",           //Sony Ericsson Xperia S
		"LT26ii",          //Sony Ericsson Xperia SL
		"LT26w",           //Sony Ericsson Xperia Acro S
		"Le X520",         //乐视 乐2
		"Le X620",         //乐视 乐2Pro
		"Le X820",         //乐视 乐Max2
		"Lenovo A3580",    //联想 黄金斗士 A8 畅玩
		"Lenovo A7600-m",  //联想 黄金斗士 S8
		"Lenovo A938t",    //联想 黄金斗士 Note8
		"Lenovo K10e70",   //联想 乐檬K10
		"Lenovo K30-T",    //联想 乐檬 K3
		"Lenovo K32C36",   //联想 乐檬3
		"Lenovo K50-t3s",  //联想 乐檬 K3 Note
		"Lenovo K52-T38",  //联想 乐檬 K5 Note
		"Lenovo K52e78",   //Lenovo K5 Note
		"Lenovo P2c72",    //联想 P2
		"Lenovo X3c50",    //联想 乐檬 X3
		"Lenovo Z90-3",    //联想 VIBE Shot大拍
		"M040",            //魅族 MX 2
		"M1 E",            //魅蓝 E
		"M2-801w",         //华为 M2
		"M2017",           //金立 M2017
		"M3",              //EBEN M3
		"M355",            //魅族 MX 3
		"MHA-AL00",        //华为 Mate 9
		"MI 4LTE",         //小米手机4
		"MI 4S",           //小米手机4S
		"MI 5",            //小米手机5
		"MI 5s Plus",      //小米手机5s Plus
		"MI 5s",           //小米手机5s
		"MI MAX",          //小米Max
		"MI Note Pro",     //小米Note顶配版
		"MI PAD 2",        //小米平板 2
		"MIX",             //小米MIX
		"MLA-UL00",        //华为 G9 Plus
		"MP1503",          //美图 M6
		"MP1512",          //美图 M6s
		"MT27i",           //Sony Ericsson Xperia Sola
		"MX4 Pro",         //魅族 MX 4 Pro
		"MX4",             //魅族 MX 4
		"MX5",             //魅族 MX 5
		"MX6",             //魅族 MX 6
		"Meitu V4s",       //美图 V4s
		"Meizu M3 Max",    //魅蓝max
		"Meizu U20",       //魅蓝U20
		"Mi 5",
		"Mi 6",
		"Mi A1",      //MI androidone
		"Mi Note 2",  //小米Note2
		"MiTV2S-48",  //小米电视2s
		"Moto G (4)", //摩托罗拉 G⁴ Plus
		"N1",         //Nokia N1
		"NCE-AL00",   //华为 畅享 6
		"NTS-AL00",   //华为 荣耀 Magic
		"NWI-AL10",   //nova2s
		"NX508J",     //努比亚 Z9
		"NX511J",     //努比亚 小牛4 Z9 Mini
		"NX512J",     //努比亚 大牛 Z9 Max
		"NX513J",     //努比亚 My 布拉格
		"NX513J",     //努比亚 布拉格S
		"NX523J",     //努比亚 Z11 Max
		"NX529J",     //努比亚 小牛5 Z11 Mini
		"NX531J",     //努比亚 Z11
		"NX549J",     //努比亚 小牛6 Z11 MiniS
		"NX563J",     //努比亚Z17
		"Nexus 4",
		"Nexus 5X",
		"Nexus 6",
		"Nexus 6P",
		"Nexus 7",
		"Nexus 9",
		"Nokia_X",       //Nokia X
		"Nokia_XL_4G",   //Nokia XL
		"ONE A2001",     //一加手机2
		"ONE E1001",     //一加手机X
		"ONEPLUS A5010", //一加5T
		"OPPO A53",      //OPPO A53
		"OPPO A59M",     //OPPO A59
		"OPPO A59s",     //OPPO A59s
		"OPPO R11",
		"OPPO R7",          //OPPO R7
		"OPPO R7Plus",      //OPPO R7Plus
		"OPPO R7S",         //OPPO R7S
		"OPPO R7sPlus",     //OPPO R7sPlus
		"OPPO R9 Plustm A", //OPPO R9Plus
		"OPPO R9s Plus",    //OPPO R9s Plus
		"OPPO R9s",
		"OPPO R9s",        //OPPO R9s
		"OPPO R9tm",       //OPPO R9
		"PE-TL10",         //华为 荣耀 6 Plus
		"PLK-TL01H",       //华为 荣耀 7
		"Pro 5",           //魅族 Pro 5
		"Pro 6",           //魅族 Pro 6
		"Pro 6s",          //魅族 Pro 6s
		"RM-1010",         //Nokia Lumia 638
		"RM-1018",         //Nokia Lumia 530
		"RM-1087",         //Nokia Lumia 930
		"RM-1090",         //Nokia Lumia 535
		"RM-867",          //Nokia Lumia 920
		"RM-875",          //Nokia Lumia 1020
		"RM-887",          //Nokia Lumia 720
		"RM-892",          //Nokia Lumia 925
		"RM-927",          //Nokia Lumia 929
		"RM-937",          //Nokia Lumia 1520
		"RM-975",          //Nokia Lumia 635
		"RM-977",          //Nokia Lumia 630
		"RM-984",          //Nokia Lumia 830
		"RM-996",          //Nokia Lumia 1320
		"Redmi 3S",        //红米3s
		"Redmi 4",         //小米 红米4
		"Redmi 4A",        //小米 红米4A
		"Redmi Note 2",    //小米 红米Note2
		"Redmi Note 3",    //小米 红米Note3
		"Redmi Note 4",    //小米 红米Note4
		"Redmi Pro",       //小米 红米Pro
		"S3",              //佳域S3
		"SCL-TL00H",       //华为 荣耀 4A
		"SD4930UR",        //Amazon Fire Phone
		"SH-03G",          //夏普 Aquos Zeta SH-03G
		"SH-04F",          //夏普 Aquos Zeta SH-04F
		"SHV31",           //夏普 Aquos Serie Mini SHV31
		"SM-A5100",        //Samsung Galaxy A5
		"SM-A7100",        //Samsung Galaxy A7
		"SM-A8000",        //Samsung Galaxy A8
		"SM-A9000",        //Samsung Galaxy A9
		"SM-A9100",        //Samsung Galaxy A9 高配版
		"SM-C5000",        //Samsung Galaxy C5
		"SM-C5010",        //Samsung Galaxy C5 Pro
		"SM-C7000",        //Samsung Galaxy C7
		"SM-C7010",        //Samsung Galaxy C7 Pro
		"SM-C9000",        //Samsung Galaxy C9 Pro
		"SM-G1600",        //Samsung Galaxy Folder
		"SM-G5500",        //Samsung Galaxy On5
		"SM-G6000",        //Samsung Galaxy On7
		"SM-G7100",        //Samsung Galaxy On7(2016)
		"SM-G7200",        //Samsung Galasy Grand Max
		"SM-G9198",        //Samsung 领世旗舰Ⅲ
		"SM-G9208",        //Samsung Galaxy S6
		"SM-G9250",        //Samsung Galasy S7 Edge
		"SM-G9280",        //Samsung Galaxy S6 Edge+
		"SM-G9300",        //Samsung Galaxy S7
		"SM-G9350",        //Samsung Galaxy S7 Edge
		"SM-G9500",        //Samsung Galaxy S8
		"SM-G9550",        //Samsung Galaxy S8+
		"SM-G9600",        //Samsung Galaxy S9
		"SM-G960F",        //Galaxy S9 Dual SIM
		"SM-G9650",        //Samsung Galaxy S9+
		"SM-G965F",        //Galaxy S9+ Dual SIM
		"SM-J3109",        //Samsung Galaxy J3
		"SM-J3110",        //Samsung Galaxy J3 Pro
		"SM-J327A",        //Samsung Galaxy J3 Emerge
		"SM-J5008",        //Samsung Galaxy J5
		"SM-J7008",        //Samsung Galaxy J7
		"SM-N9108V",       //Samsung Galasy Note4
		"SM-N9200",        //Samsung Galaxy Note5
		"SM-N9300",        //Samsung Galaxy Note 7
		"SM-N935S",        //Samsung Galaxy Note Fan Edition
		"SM-N9500",        //Samsung Galasy Note8
		"SM-W2015",        //Samsung W2015
		"SM-W2016",        //Samsung W2016
		"SM-W2017",        //Samsung W2017
		"SM705",           //锤子 T1
		"SM801",           //锤子 T2
		"SM901",           //锤子 M1
		"SM919",           //锤子 M1L
		"ST18i",           //Sony Ericsson Xperia Ray
		"ST25i",           //Sony Ericsson Xperia U
		"STV100-1",        //黑莓Priv
		"Signature Touch", //Vertu Signature Touch
		"TA-1000",         //Nokia 6
		"TA-1000",         //HMD Nokia 6
		"TA-1041",         //Nokia 7
		"VERTU Ti",        //Vertu Ti
		"VIE-AL10",        //华为 P9 Plus
		"VIVO X20",
		"VIVO X20A",
		"W909",           //金立 天鉴 W909
		"X500",           //乐视 乐1S
		"X608",           //乐视 乐1
		"X800",           //乐视 乐1Pro
		"X900",           //乐视 乐Max
		"XT1085",         //摩托罗拉 X
		"XT1570",         //摩托罗拉 X Style
		"XT1581",         //摩托罗拉 X 极
		"XT1585",         //摩托罗拉 Droid Turbo 2
		"XT1635",         //摩托罗拉 Z Play
		"XT1635-02",      //摩托罗拉 Z Play
		"XT1650",         //摩托罗拉 Z
		"XT1650-05",      //摩托罗拉 Z
		"XT1706",         //摩托罗拉 E³ POWER
		"YD201",          //YotaPhone2
		"YD206",          //YotaPhone2
		"YQ60",           //锤子 坚果
		"ZTE A2015",      //中兴 AXON 天机
		"ZTE A2017",      //中兴 AXON 天机 7
		"ZTE B2015",      //中兴 AXON 天机 MINI
		"ZTE BV0720",     //中兴 Blade A2
		"ZTE BV0730",     //中兴 Blade A2 Plus
		"ZTE C2016",      //中兴 AXON 天机 MAX
		"ZTE C2017",      //中兴 AXON 天机 7 MAX
		"ZTE G720C",      //中兴 星星2号
		"ZUK Z2121",      //ZUK Z2 Pro
		"ZUK Z2131",      //ZUK Z2
		"ZUK Z2151",      //ZUK Edge
		"ZUK Z2155",      //ZUK Edge L
		"m030",           //魅族mx
		"m1 metal",       //魅蓝metal
		"m1 note",        //魅蓝 Note
		"m1",             //魅蓝
		"m2 note",        //魅蓝 Note 2
		"m2",             //魅蓝 2
		"m3 note",        //魅蓝 Note 3
		"m3",             //魅蓝 3
		"m3s",            //魅蓝 3S
		"m9",             //魅族m9
		"marlin",         //Google Pixel XL
		"sailfish",       //Google Pixel
		"vivo V3Max",     //vivo V3Max
		"vivo X6D",       //vivo X6
		"vivo X6PlusD",   //vivo X6Plus
		"vivo X6S",       //vivo X6S
		"vivo X6SPlus",   //vivo X6SPlus
		"vivo X7",        //vivo X7
		"vivo X7Plus",    //vivo X7Plus
		"vivo X9",        //vivo X9
		"vivo X9Plus",    //vivo X9Plus
		"vivo Xplay5A 金", //vivo Xplay5
		"vivo Xplay6",    //vivo Xplay6
		"vivo Y66",       //vivo Y66
		"vivo Y67",       //vivo Y67
		"z1221",          //ZUK Z1
	}
)
