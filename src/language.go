package main

var (
	languageMap map[string]string
)

func languageInit() {
	languageMap = map[string]string{
		"en": "English",
		"fr": "French",
		"es": "Spanish",
		"ru": "Russian",
		"cn": "Chinese",
	}
}
