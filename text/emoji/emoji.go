package emoji

/*
https://emojipedia.org/github/
http://emojicodes.com/
https://unicode.org/emoji/charts/full-emoji-list.html
*/

import (
	"regexp"
	"strings"
)

const gomojiRaw string = `
:+1: +1
:angry:	:@
:broken_heart:	</3
:confused:	>:\
:cry:	:'(
:disappointed:	:(
:dizzy_face:	#)
:expressionless:	-_-
:fearful:	D:
:flushed:	:$
:frowning: :( 🙁
:heart:	<3
:innocent:	O:)
:joy:	:')
:kissing_heart:	:^*
:laughing:	>:)
:no_mouth:	:X
:ok_woman:	*\0/*
:open_mouth:	>:O
:persevere:	>.<
:slight_smile:	:)
:smile: :) 😀
:smiley:	:D
:stuck_out_tongue:	:P
:stuck_out_tongue_winking_eye:	>:P
:sunglasses:	B)
:sweat:	':(
:sweat_smile:	':)
:wink: ;) 😉`

func GetEmojiToAsciiMap() map[string]Emoji {
	data := map[string]Emoji{}
	rx := regexp.MustCompile(`\s+`)
	lines := strings.Split(gomojiRaw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = rx.ReplaceAllString(line, " ")
		parts := strings.Split(line, " ")
		if len(parts) == 2 || len(parts) == 3 {
			emo := Emoji{
				Shortcode:   strings.TrimSpace(parts[0]),
				Ascii:       strings.TrimSpace(parts[1]),
				ShortcodeRx: regexp.MustCompile(regexp.QuoteMeta(parts[0]))}
			if len(parts) == 3 {
				emo.Unicode = strings.TrimSpace(parts[2])
			}
			data[parts[0]] = emo
		}
	}
	return data
}

type Emoji struct {
	Ascii       string
	Shortcode   string
	Unicode     string
	ShortcodeRx *regexp.Regexp
}

type Converter struct {
	data       map[string]Emoji
	UseUnicode bool
}

func NewConverter() Converter {
	return Converter{
		data:       GetEmojiToAsciiMap(),
		UseUnicode: true}
}

func (conv *Converter) EmojiToAscii(input string) string {
	rx := regexp.MustCompile(`:\+?[0-9a-z_]+:`)
	matches := rx.FindAllString(input, -1)
	output := input
	for _, emo := range matches {
		if einfo, ok := conv.data[emo]; ok {
			if conv.UseUnicode && len(einfo.Unicode) > 0 {
				output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.Unicode)
			} else {
				output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.Ascii)
			}
		}
	}
	return output
}
