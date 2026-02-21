package filters

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/PriestFaria/lingo/internal/analyzer/log"
)

// EmojiStrictFilter reports log messages that contain emoji characters or
// repeated punctuation sequences (e.g. !!, ???, ...).
type EmojiStrictFilter struct{}

// repeatedPunct matches two or more consecutive ! or ? and two or more dots.
var repeatedPunct = regexp.MustCompile(`[!?]{2,}|\.{2,}`)

// emojiRanges is a unicode.RangeTable that covers all emoji code points
// defined in the Unicode Emoji specification v15.1.
// R16 entries fit in uint16; R32 entries require uint32.
var emojiRanges = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x00A9, 0x00A9, 1}, // © Copyright
		{0x00AE, 0x00AE, 1}, // ® Registered
		{0x203C, 0x203C, 1}, // ‼ Double exclamation mark
		{0x2049, 0x2049, 1}, // ⁉ Exclamation question mark
		{0x2122, 0x2122, 1}, // ™ Trade mark sign
		{0x2139, 0x2139, 1}, // ℹ Information source
		{0x2194, 0x2199, 1}, // ↔–↙ Arrows
		{0x21A9, 0x21AA, 1}, // ↩↪ Arrows with hook
		{0x231A, 0x231B, 1}, // ⌚⌛ Watch, hourglass
		{0x2328, 0x2328, 1}, // ⌨ Keyboard
		{0x23CF, 0x23CF, 1}, // ⏏ Eject symbol
		{0x23E9, 0x23F3, 1}, // ⏩–⏳ Fast-forward, clocks
		{0x23F8, 0x23FA, 1}, // ⏸–⏺ Pause, stop, record
		{0x24C2, 0x24C2, 1}, // Ⓜ Circled M
		{0x25AA, 0x25AB, 1}, // ▪▫ Small squares
		{0x25B6, 0x25B6, 1}, // ▶ Black right-pointing triangle
		{0x25C0, 0x25C0, 1}, // ◀ Black left-pointing triangle
		{0x25FB, 0x25FE, 1}, // ◻–◾ Medium squares
		{0x2600, 0x2604, 1}, // ☀–🌄 Misc symbols block start
		{0x260E, 0x260E, 1}, // ☎ Black telephone
		{0x2611, 0x2611, 1}, // ☑ Ballot box with check
		{0x2614, 0x2615, 1}, // ☔☕ Umbrella, hot beverage
		{0x2618, 0x2618, 1}, // ☘ Shamrock
		{0x261D, 0x261D, 1}, // ☝ White up-pointing index
		{0x2620, 0x2620, 1}, // ☠ Skull and crossbones
		{0x2622, 0x2623, 1}, // ☢☣ Radioactive, biohazard
		{0x2626, 0x2626, 1}, // ☦ Orthodox cross
		{0x262A, 0x262A, 1}, // ☪ Star and crescent
		{0x262E, 0x262F, 1}, // ☮☯ Peace, yin yang
		{0x2638, 0x263A, 1}, // ☸–☺ Wheel, frowning, smiling face
		{0x2640, 0x2640, 1}, // ♀ Female sign
		{0x2642, 0x2642, 1}, // ♂ Male sign
		{0x2648, 0x2653, 1}, // ♈–♓ Zodiac signs
		{0x265F, 0x2660, 1}, // ♟♠ Chess pawn, spade suit
		{0x2663, 0x2663, 1}, // ♣ Club suit
		{0x2665, 0x2666, 1}, // ♥♦ Heart and diamond suits
		{0x2668, 0x2668, 1}, // ♨ Hot springs
		{0x267B, 0x267B, 1}, // ♻ Black universal recycling symbol
		{0x267E, 0x267F, 1}, // ♾♿ Infinity, wheelchair
		{0x2692, 0x2697, 1}, // ⚒–⚗ Hammer&pick through alembic
		{0x2699, 0x2699, 1}, // ⚙ Gear
		{0x269B, 0x269C, 1}, // ⚛⚜ Atom, fleur-de-lis
		{0x26A0, 0x26A1, 1}, // ⚠⚡ Warning, lightning
		{0x26A7, 0x26A7, 1}, // ⚧ Male with stroke and male and female sign
		{0x26AA, 0x26AB, 1}, // ⚪⚫ Medium circles
		{0x26B0, 0x26B1, 1}, // ⚰⚱ Coffin, funeral urn
		{0x26BD, 0x26BE, 1}, // ⚽⚾ Soccer ball, baseball
		{0x26C4, 0x26C5, 1}, // ⛄⛅ Snowman, sun behind cloud
		{0x26CE, 0x26CF, 1}, // ⛎⛏ Ophiuchus, pick
		{0x26D1, 0x26D1, 1}, // ⛑ Helmet with white cross
		{0x26D3, 0x26D4, 1}, // ⛓⛔ Chains, no entry
		{0x26E9, 0x26EA, 1}, // ⛩⛪ Shinto shrine, church
		{0x26F0, 0x26F5, 1}, // ⛰–⛵ Mountain through sailboat
		{0x26F7, 0x26FA, 1}, // ⛷–⛺ Skier through tent
		{0x26FD, 0x26FD, 1}, // ⛽ Fuel pump
		{0x2702, 0x2702, 1}, // ✂ Black scissors
		{0x2705, 0x2705, 1}, // ✅ White heavy check mark
		{0x2708, 0x270D, 1}, // ✈–✍ Airplane through writing hand
		{0x270F, 0x270F, 1}, // ✏ Pencil
		{0x2712, 0x2712, 1}, // ✒ Black nib
		{0x2714, 0x2714, 1}, // ✔ Heavy check mark
		{0x2716, 0x2716, 1}, // ✖ Heavy multiplication X
		{0x271D, 0x271D, 1}, // ✝ Latin cross
		{0x2721, 0x2721, 1}, // ✡ Star of David
		{0x2728, 0x2728, 1}, // ✨ Sparkles
		{0x2733, 0x2734, 1}, // ✳✴ Eight-spoked/pointed asterisk
		{0x2744, 0x2744, 1}, // ❄ Snowflake
		{0x2747, 0x2747, 1}, // ❇ Sparkle
		{0x274C, 0x274C, 1}, // ❌ Cross mark
		{0x274E, 0x274E, 1}, // ❎ Cross mark button
		{0x2753, 0x2755, 1}, // ❓❔❕ Question marks
		{0x2757, 0x2757, 1}, // ❗ Heavy exclamation mark ornament
		{0x2763, 0x2764, 1}, // ❣❤ Heart exclamation, heart
		{0x2795, 0x2797, 1}, // ➕➖➗ Plus, minus, division
		{0x27A1, 0x27A1, 1}, // ➡ Black rightwards arrow
		{0x27B0, 0x27B0, 1}, // ➰ Curly loop
		{0x27BF, 0x27BF, 1}, // ➿ Double curly loop
		{0x2934, 0x2935, 1}, // ⤴⤵ Arrows
		{0x2B05, 0x2B07, 1}, // ⬅–⬇ Arrows
		{0x2B1B, 0x2B1C, 1}, // ⬛⬜ Large squares
		{0x2B50, 0x2B50, 1}, // ⭐ White medium star
		{0x2B55, 0x2B55, 1}, // ⭕ Heavy large circle
		{0x3030, 0x3030, 1}, // 〰 Wavy dash
		{0x303D, 0x303D, 1}, // 〽 Part alternation mark
		{0x3297, 0x3297, 1}, // ㊗ Circled ideograph congratulation
		{0x3299, 0x3299, 1}, // ㊙ Circled ideograph secret
	},
	R32: []unicode.Range32{
		{0x1F004, 0x1F004, 1}, // 🀄 Mahjong Red Dragon
		{0x1F0CF, 0x1F0CF, 1}, // 🃏 Playing card black joker
		{0x1F170, 0x1F171, 1}, // 🅰🅱 Blood type buttons
		{0x1F17E, 0x1F17F, 1}, // 🅾🅿 Blood type / parking buttons
		{0x1F18E, 0x1F18E, 1}, // 🆎 AB button
		{0x1F191, 0x1F19A, 1}, // 🆑–🆚 Squared Latin buttons
		{0x1F1E0, 0x1F1FF, 1}, // Regional indicator symbols (flag sequences)
		{0x1F201, 0x1F202, 1}, // 🈁🈂 Japanese buttons
		{0x1F21A, 0x1F21A, 1}, // 🈚 Japanese "free of charge"
		{0x1F22F, 0x1F22F, 1}, // 🈯 Japanese "reserved"
		{0x1F232, 0x1F23A, 1}, // 🈲–🈺 Japanese CJK buttons
		{0x1F250, 0x1F251, 1}, // 🉐🉑 Japanese "bargain"/"acceptable"
		{0x1F300, 0x1F321, 1}, // 🌀–🌡 Misc symbols & pictographs
		{0x1F324, 0x1F393, 1}, // 🌤–🎓 Weather, activities
		{0x1F396, 0x1F397, 1}, // 🎖🎗 Military medal, reminder ribbon
		{0x1F399, 0x1F39B, 1}, // 🎙–🎛 Studio microphone, knob
		{0x1F39E, 0x1F3F0, 1}, // 🎞–🏰 Film frames through castle
		{0x1F3F3, 0x1F3F5, 1}, // 🏳–🏵 White flag through rosette
		{0x1F3F7, 0x1F4FD, 1}, // 🏷–📽 Label through film projector
		{0x1F4FF, 0x1F53D, 1}, // 📿–🔽 Prayer beads through downward button
		{0x1F549, 0x1F54E, 1}, // 🕉–🕎 Om through menorah
		{0x1F550, 0x1F567, 1}, // 🕐–🕧 Clock faces
		{0x1F56F, 0x1F570, 1}, // 🕯🕰 Candle, mantelpiece clock
		{0x1F573, 0x1F57A, 1}, // 🕳–🕺 Hole through man dancing
		{0x1F587, 0x1F587, 1}, // 🖇 Linked paperclips
		{0x1F58A, 0x1F58D, 1}, // 🖊–🖍 Pens and crayon
		{0x1F590, 0x1F590, 1}, // 🖐 Raised hand with fingers splayed
		{0x1F595, 0x1F596, 1}, // 🖕🖖 Middle finger, vulcan salute
		{0x1F5A4, 0x1F5A5, 1}, // 🖤🖥 Black heart, desktop computer
		{0x1F5A8, 0x1F5A8, 1}, // 🖨 Printer
		{0x1F5B1, 0x1F5B2, 1}, // 🖱🖲 Computer mouse, trackball
		{0x1F5BC, 0x1F5BC, 1}, // 🖼 Frame with picture
		{0x1F5C2, 0x1F5C4, 1}, // 🗂–🗄 Card index dividers, cabinet
		{0x1F5D1, 0x1F5D3, 1}, // 🗑–🗓 Wastebasket, spiral calendars
		{0x1F5DC, 0x1F5DE, 1}, // 🗜–🗞 Compression, rolled-up newspaper
		{0x1F5E1, 0x1F5E1, 1}, // 🗡 Dagger knife
		{0x1F5E3, 0x1F5E3, 1}, // 🗣 Speaking head in silhouette
		{0x1F5E8, 0x1F5E8, 1}, // 🗨 Left speech bubble
		{0x1F5EF, 0x1F5EF, 1}, // 🗯 Right anger bubble
		{0x1F5F3, 0x1F5F3, 1}, // 🗳 Ballot box with ballot
		{0x1F5FA, 0x1F64F, 1}, // 🗺–🙏 World map through folded hands
		{0x1F680, 0x1F6C5, 1}, // 🚀–🛅 Transport & map symbols
		{0x1F6CB, 0x1F6D2, 1}, // 🛋–🛒 Couch through shopping trolley
		{0x1F6D5, 0x1F6D7, 1}, // 🛕–🛗 Hindu temple, elevator
		{0x1F6DC, 0x1F6E5, 1}, // 🛜–🛥 Wireless, motor boat
		{0x1F6E9, 0x1F6E9, 1}, // 🛩 Small airplane
		{0x1F6EB, 0x1F6EC, 1}, // 🛫🛬 Airplane departure/arrival
		{0x1F6F0, 0x1F6F0, 1}, // 🛰 Satellite
		{0x1F6F3, 0x1F6FC, 1}, // 🛳–🛼 Passenger ship through roller skate
		{0x1F7E0, 0x1F7EB, 1}, // 🟠–🟫 Colored circles and squares
		{0x1F7F0, 0x1F7F0, 1}, // 🟰 Heavy equals sign
		{0x1F90C, 0x1F9FF, 1}, // 🤌–🧿 Supplemental symbols & pictographs
		{0x1FA00, 0x1FA53, 1}, // 🨀–🩓 Chess symbols
		{0x1FA60, 0x1FA6D, 1}, // 🩠–🩭 Game pieces
		{0x1FA70, 0x1FA7C, 1}, // 🩰–🩼 Medical symbols
		{0x1FA80, 0x1FA88, 1}, // 🪀–🪈 Yo-yo through flute
		{0x1FA90, 0x1FABD, 1}, // 🪐–🪽 Ringed planet through wing
		{0x1FABF, 0x1FAC5, 1}, // 🪿–🫅 Goose through person with crown
		{0x1FACE, 0x1FADB, 1}, // 🫎–🫛 Moose through pea pod
		{0x1FAE0, 0x1FAE8, 1}, // 🫠–🫨 Melting face through shaking face
		{0x1FAF0, 0x1FAF8, 1}, // 🫰–🫸 Hand with index finger and thumb
	},
	LatinOffset: 0,
}

// isEmoji reports whether the rune is an emoji character per Unicode v15.1.
func isEmoji(r rune) bool {
	return unicode.Is(emojiRanges, r)
}

func (f *EmojiStrictFilter) Apply(context *log.LogContext) []FilterIssue {
	var issues []FilterIssue
	for _, part := range context.Parts {
		if !part.IsLiteral {
			continue
		}

		for _, r := range part.Value {
			if isEmoji(r) {
				issues = append(issues, FilterIssue{
					Message: fmt.Sprintf("log message must not contain emoji: %q", r),
					Pos:     part.Pos,
				})
				break
			}
		}

		if loc := repeatedPunct.FindStringIndex(part.Value); loc != nil {
			issues = append(issues, FilterIssue{
				Message: fmt.Sprintf("log message must not contain repeated punctuation: %q", part.Value[loc[0]:loc[1]]),
				Pos:     part.Pos,
			})
		}
	}
	return issues
}