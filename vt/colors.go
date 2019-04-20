package vt

import (
	"fmt"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	// Colors contains a list of terminal colors
	// https://jonasjacek.github.io/colors/
	Colors = map[uint8]colorful.Color{
		0:   colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		1:   colorful.Color{R: 128 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		2:   colorful.Color{R: 0 / 255.0, G: 128 / 255.0, B: 0 / 255.0},
		3:   colorful.Color{R: 128 / 255.0, G: 128 / 255.0, B: 0 / 255.0},
		4:   colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 128 / 255.0},
		5:   colorful.Color{R: 128 / 255.0, G: 0 / 255.0, B: 128 / 255.0},
		6:   colorful.Color{R: 0 / 255.0, G: 128 / 255.0, B: 128 / 255.0},
		7:   colorful.Color{R: 192 / 255.0, G: 192 / 255.0, B: 192 / 255.0},
		8:   colorful.Color{R: 128 / 255.0, G: 128 / 255.0, B: 128 / 255.0},
		9:   colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		10:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		11:  colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		12:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		13:  colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		14:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		15:  colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		16:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		17:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		18:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		19:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		20:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		21:  colorful.Color{R: 0 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		22:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		23:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		24:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		25:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		26:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		27:  colorful.Color{R: 0 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		28:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		29:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		30:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		31:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		32:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		33:  colorful.Color{R: 0 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		34:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		35:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		36:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		37:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		38:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		39:  colorful.Color{R: 0 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		40:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		41:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		42:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		43:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		44:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		45:  colorful.Color{R: 0 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		46:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		47:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		48:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		49:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		50:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		51:  colorful.Color{R: 0 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		52:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		53:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		54:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		55:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		56:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		57:  colorful.Color{R: 95 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		58:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		59:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		60:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		61:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		62:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		63:  colorful.Color{R: 95 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		64:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		65:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		66:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		67:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		68:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		69:  colorful.Color{R: 95 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		70:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		71:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		72:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		73:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		74:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		75:  colorful.Color{R: 95 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		76:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		77:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		78:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		79:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		80:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		81:  colorful.Color{R: 95 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		82:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		83:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		84:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		85:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		86:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		87:  colorful.Color{R: 95 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		88:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		89:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		90:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		91:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		92:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		93:  colorful.Color{R: 135 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		94:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		95:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		96:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		97:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		98:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		99:  colorful.Color{R: 135 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		100: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		101: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		102: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		103: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		104: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		105: colorful.Color{R: 135 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		106: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		107: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		108: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		109: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		110: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		111: colorful.Color{R: 135 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		112: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		113: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		114: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		115: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		116: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		117: colorful.Color{R: 135 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		118: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		119: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		120: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		121: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		122: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		123: colorful.Color{R: 135 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		124: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		125: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		126: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		127: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		128: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		129: colorful.Color{R: 175 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		130: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		131: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		132: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		133: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		134: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		135: colorful.Color{R: 175 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		136: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		137: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		138: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		139: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		140: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		141: colorful.Color{R: 175 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		142: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		143: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		144: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		145: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		146: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		147: colorful.Color{R: 175 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		148: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		149: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		150: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		151: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		152: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		153: colorful.Color{R: 175 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		154: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		155: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		156: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		157: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		158: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		159: colorful.Color{R: 175 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		160: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		161: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		162: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		163: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		164: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		165: colorful.Color{R: 215 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		166: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		167: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		168: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		169: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		170: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		171: colorful.Color{R: 215 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		172: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		173: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		174: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		175: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		176: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		177: colorful.Color{R: 215 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		178: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		179: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		180: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		181: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		182: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		183: colorful.Color{R: 215 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		184: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		185: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		186: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		187: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		188: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		189: colorful.Color{R: 215 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		190: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		191: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		192: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		193: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		194: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		195: colorful.Color{R: 215 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		196: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 0 / 255.0},
		197: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 95 / 255.0},
		198: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 135 / 255.0},
		199: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 175 / 255.0},
		200: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 215 / 255.0},
		201: colorful.Color{R: 255 / 255.0, G: 0 / 255.0, B: 255 / 255.0},
		202: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 0 / 255.0},
		203: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 95 / 255.0},
		204: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 135 / 255.0},
		205: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 175 / 255.0},
		206: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 215 / 255.0},
		207: colorful.Color{R: 255 / 255.0, G: 95 / 255.0, B: 255 / 255.0},
		208: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 0 / 255.0},
		209: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 95 / 255.0},
		210: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 135 / 255.0},
		211: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 175 / 255.0},
		212: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 215 / 255.0},
		213: colorful.Color{R: 255 / 255.0, G: 135 / 255.0, B: 255 / 255.0},
		214: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 0 / 255.0},
		215: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 95 / 255.0},
		216: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 135 / 255.0},
		217: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 175 / 255.0},
		218: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 215 / 255.0},
		219: colorful.Color{R: 255 / 255.0, G: 175 / 255.0, B: 255 / 255.0},
		220: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 0 / 255.0},
		221: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 95 / 255.0},
		222: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 135 / 255.0},
		223: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 175 / 255.0},
		224: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 215 / 255.0},
		225: colorful.Color{R: 255 / 255.0, G: 215 / 255.0, B: 255 / 255.0},
		226: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 0 / 255.0},
		227: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 95 / 255.0},
		228: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 135 / 255.0},
		229: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 175 / 255.0},
		230: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 215 / 255.0},
		231: colorful.Color{R: 255 / 255.0, G: 255 / 255.0, B: 255 / 255.0},
		232: colorful.Color{R: 8 / 255.0, G: 8 / 255.0, B: 8 / 255.0},
		233: colorful.Color{R: 18 / 255.0, G: 18 / 255.0, B: 18 / 255.0},
		234: colorful.Color{R: 28 / 255.0, G: 28 / 255.0, B: 28 / 255.0},
		235: colorful.Color{R: 38 / 255.0, G: 38 / 255.0, B: 38 / 255.0},
		236: colorful.Color{R: 48 / 255.0, G: 48 / 255.0, B: 48 / 255.0},
		237: colorful.Color{R: 58 / 255.0, G: 58 / 255.0, B: 58 / 255.0},
		238: colorful.Color{R: 68 / 255.0, G: 68 / 255.0, B: 68 / 255.0},
		239: colorful.Color{R: 78 / 255.0, G: 78 / 255.0, B: 78 / 255.0},
		240: colorful.Color{R: 88 / 255.0, G: 88 / 255.0, B: 88 / 255.0},
		241: colorful.Color{R: 98 / 255.0, G: 98 / 255.0, B: 98 / 255.0},
		242: colorful.Color{R: 108 / 255.0, G: 108 / 255.0, B: 108 / 255.0},
		243: colorful.Color{R: 118 / 255.0, G: 118 / 255.0, B: 118 / 255.0},
		244: colorful.Color{R: 128 / 255.0, G: 128 / 255.0, B: 128 / 255.0},
		245: colorful.Color{R: 138 / 255.0, G: 138 / 255.0, B: 138 / 255.0},
		246: colorful.Color{R: 148 / 255.0, G: 148 / 255.0, B: 148 / 255.0},
		247: colorful.Color{R: 158 / 255.0, G: 158 / 255.0, B: 158 / 255.0},
		248: colorful.Color{R: 168 / 255.0, G: 168 / 255.0, B: 168 / 255.0},
		249: colorful.Color{R: 178 / 255.0, G: 178 / 255.0, B: 178 / 255.0},
		250: colorful.Color{R: 188 / 255.0, G: 188 / 255.0, B: 188 / 255.0},
		251: colorful.Color{R: 198 / 255.0, G: 198 / 255.0, B: 198 / 255.0},
		252: colorful.Color{R: 208 / 255.0, G: 208 / 255.0, B: 208 / 255.0},
		253: colorful.Color{R: 218 / 255.0, G: 218 / 255.0, B: 218 / 255.0},
		254: colorful.Color{R: 228 / 255.0, G: 228 / 255.0, B: 228 / 255.0},
		255: colorful.Color{R: 238 / 255.0, G: 238 / 255.0, B: 238 / 255.0},
	}

	colStore = map[int64]uint8{}
	colMutex sync.Mutex
)

// GetRGBInt takes in an RGB int64 and return a terminal color ID
func GetRGBInt(rgb int64) uint8 {
	if u, ok := colStore[rgb]; ok {
		return u
	}

	var (
		b = rgb & 255
		g = (rgb >> 8) & 255
		r = (rgb >> 16) & 255
	)

	u := GetColorInt(
		float64(r)/255,
		float64(g)/255,
		float64(b)/255,
	)

	colMutex.Lock()
	defer colMutex.Unlock()

	colStore[rgb] = u
	return u
}

// GetColorInt takes in (r|g|b)/255
func GetColorInt(r, g, b float64) uint8 {
	var (
		c = colorful.Color{R: r, G: g, B: b}
		m float64
		j uint
	)

	for i := uint(0); i < 256; i++ {
		rgb := c.DistanceRgb(Colors[uint8(i)])
		if i == 0 || rgb < m {
			m = rgb
			j = i
		}
	}

	return uint8(j)
}

// FmtColorForeground turns a terminal color ID into string
func FmtColorForeground(c uint8) string {
	return fmt.Sprintf("\033[38;5;%dm", c)
}

func ColorString(c uint8, s string) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", c, s)
}
