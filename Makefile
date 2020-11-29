
all:

extract:
	goi18n extract -sourceLanguage en-us -outdir l10n -format toml

translate:
ifeq ($(lang),)
	echo "use: make translate lang=de-DE"
	exit
else
	touch translate.$(lang).toml
	goi18n merge -sourceLanguage en-us l10n/active.en-US.toml translate.$(lang).toml
endif
