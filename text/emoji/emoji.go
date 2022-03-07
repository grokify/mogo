package emoji

/*
https://emojipedia.org/github/
http://emojicodes.com/
https://emojipedia.org/emoji/%F0%9F%86%98/
https://unicode.org/emoji/charts/full-emoji-list.html
*/

import (
	"regexp"
	"strings"
)

const gomojiRaw string = `
:angry:	😠	:@
:anguished:	😧
:astonished:	😲
:confused:	😕	>:\
:cry:	😢	:'(
:disappointed:	😞	:(
:dizzy_face:	😵	#)
:expressionless:	😑	-_-
:fearful:	😨	D:
:flushed:	😳	:$
:frowning:	🙁	:(
:innocent:	😇	O:)
:joy:	😂	:')
:kissing_heart:	:^*
:laughing:	😆	>:)
:no_mouth:	:X
:ok_woman:	*\0/*
:open_mouth:	😮	>:O
:persevere:	😣	>.<
:relaxed:	🙂
:relieved:	😌
:scream:	😱
:slight_smile:	🙂	:)
:smile:	😀	:)
:smiley:	😃	:D
:snowflake:	❄
:snowman:	☃
:stuck_out_tongue:	😛	:P
:stuck_out_tongue_winking_eye:	😜	>:P
:sunglasses:	😎	B)
:sweat:	😰	':(
:sweat_smile:	😅	':)
:victory_hand:	✌
:wink:	😉	;)

:+1:	👍	+1
:-1:	👎	-1
:love_you_gesture:	🤟
:ok_hand:	👌

:alien:	👽
:ghost:	👻
:goblin:	👺
:ogre:	👹
:robot:	🤖
:skull:	💀

:beer:	🍺
:beers:	🍻
:boom:	💥
:chopsticks:	🥢
:droplet:	💧
:exclamation:	❗	!
:fire:	🔥
:minus_sign:	➖
:no_entry:	⛔
:sos:	🆘	SOS
:spoon:	🥄
:sun:,:sunny:	☀️
:umbrella:	☂️
:white_check_mark:	✅

:black_heart:	🖤
:blue_heart;	💙
:broken_heart:	💔	</3
:green_heart:	💚
:heart:	🧡	<3
:heart_declaration:	💟
:heart_exclamation:	❣
:orange_heart:	🧡
:purple_heart:	💜
:red_heart:	❤
:revolving_heart:	💞
:two_hearts:	💕
:yellow_heart:	💛

:kiss_mark:	💋
:love_letter:	💌

:crying_cat:	😿
:grinning_cat:	😺
:pouting_cat:	😾
:weary_cat:	🙀

:black_flag:	🏴
:checkered_flag:	🏁
:crossed_flags:	🎌
:rainbow_flag:	🏳️‍🌈
:triangular_flag:	🚩
:white_flag:	🏳
`

/*
:relaxed: Y
:smiley:	Y
:relieved: Y
:green_heart:	Y
:+1: Y
:ok_hand:	Y
:sunny:	Y
:beers:	Y
:white_check_mark:	Y

:frowning: Y
:anguished: Y
:open_mouth:	Y
:confused:	Y
:scream:	Y
:broken_heart:	Y
:boom: Y
:exclamation: Y
:fire:	Y
:-1: Y
:umbrella:	Y
:sos:	Y

> DOWN_EMOJI = %w(:frowning: :anguished: :open_mouth: :confused: :scream: :broken_heart: :boom: :exclamation: :fire: :-1: :umbrella: :sos:)
> UP_EMOJI = %w(:relaxed: :smiley: :relieved: :green_heart: :+1: :ok_hand: :sunny: :beers: :white_check_mark:)
*/

var rxEmojiShortcode *regexp.Regexp = regexp.MustCompile(`:[\+\-]?[0-9a-z_]+:`)

func GetEmojiDataShortcodeMap() map[string]Emoji {
	data := map[string]Emoji{}
	rx := regexp.MustCompile(`\s+`)
	lines := strings.Split(gomojiRaw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Index(line, ":") != 0 {
			continue
		}
		line = rx.ReplaceAllString(line, " ")
		parts := strings.Split(line, " ")
		if len(parts) == 2 || len(parts) == 3 {
			emo := Emoji{
				Shortcode:   strings.TrimSpace(parts[0]),
				Unicode:     strings.TrimSpace(parts[1]),
				ShortcodeRx: regexp.MustCompile(regexp.QuoteMeta(parts[0]))}
			if len(parts) == 3 {
				emo.ASCII = strings.TrimSpace(parts[2])
			}
			data[parts[0]] = emo
		} else if len(parts) > 3 {
			panic("E_BAD_FORMATTING")
		}
	}
	return data
}

type Emoji struct {
	ASCII       string
	Shortcode   string
	Unicode     string
	ShortcodeRx *regexp.Regexp
}

type EmojiType int

const (
	Shortcode EmojiType = iota
	ASCII
	Unicode
)

type Converter struct {
	data map[string]Emoji
}

func NewConverter() Converter { return Converter{data: GetEmojiDataShortcodeMap()} }

func (conv *Converter) ConvertShortcodesString(input string, emoType EmojiType) string {
	if emoType == Shortcode {
		return input
	}
	matches := rxEmojiShortcode.FindAllString(input, -1)
	output := input
	for _, emoShortcode := range matches {
		if einfo, ok := conv.data[emoShortcode]; ok {
			if emoType == Unicode {
				if len(einfo.Unicode) > 0 {
					output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.Unicode)
				} else {
					output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.ASCII)
				}
			} else {
				if len(einfo.ASCII) > 0 {
					output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.ASCII)
				} else {
					output = einfo.ShortcodeRx.ReplaceAllString(output, einfo.Unicode)
				}
			}
		}
	}
	return output
}
