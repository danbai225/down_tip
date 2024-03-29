package keylog

import (
	"runtime"
)

var keyMap = map[uint16]string{
	14:    "BackSpace",
	15:    "Tab",
	29:    "Left Control",
	0:     "Menu",
	28:    "Enter",
	42:    "Left Shift",
	54:    "Right Shift",
	3613:  "Right Control",
	3675:  "Win",
	56:    "Left Alt",
	3640:  "Right Alt",
	58:    "Cape Lock",
	1:     "Esc",
	57:    "Space Bar",
	3657:  "Page Up",
	3665:  "Page Down",
	3663:  "End",
	3655:  "Home",
	3666:  "Insert",
	3667:  "Delete",
	57419: "Left Arrow",
	57416: "Up Arrow",
	57421: "Right Arrow",
	57424: "Dw Arrow",
	11:    "0",
	2:     "1",
	3:     "2",
	4:     "3",
	5:     "4",
	6:     "5",
	7:     "6",
	8:     "7",
	9:     "8",
	10:    "9",
	30:    "A",
	48:    "B",
	46:    "C",
	32:    "D",
	18:    "E",
	33:    "F",
	34:    "G",
	35:    "H",
	23:    "I",
	36:    "J",
	37:    "K",
	38:    "L",
	50:    "M",
	49:    "N",
	24:    "O",
	25:    "P",
	16:    "Q",
	19:    "R",
	31:    "S",
	20:    "T",
	22:    "U",
	47:    "V",
	17:    "W",
	45:    "X",
	21:    "Y",
	44:    "Z",
	82:    "Num 0",
	79:    "Num 1",
	80:    "Num 2",
	81:    "Num 3",
	75:    "Num 4",
	76:    "Num 5",
	77:    "Num 6",
	71:    "Num 7",
	72:    "Num 8",
	73:    "Num 9",
	55:    "Num *",
	78:    "Num +",
	3612:  "Num Enter",
	74:    "Num -",
	83:    "Num .",
	3637:  "Num /",
	59:    "F1",
	60:    "F2",
	61:    "F3",
	62:    "F4",
	63:    "F5",
	64:    "F6",
	65:    "F7",
	66:    "F8",
	67:    "F9",
	68:    "F10",
	87:    "F11",
	88:    "F12",
	91:    "PrtSc",
	92:    "ScrLk",
	93:    "Pause",
	69:    "Num Lock",
	39:    ";:",
	13:    "\u003d+",
	51:    ",\u003c\u003e",
	12:    "-_",
	52:    ".\u003e",
	53:    "/?",
	41:    "`~",
	26:    "[{",
	43:    "\\|",
	27:    "]}",
	40:    "\u0027\"",
}

func init() {
	var winMap = map[uint16]string{
		61003: "Left Arrow",
		61000: "Up Arrow",
		61005: "Right Arrow",
		61008: "Dw Arrow",
		61001: "Page Up",
		61009: "Page Down",
		61007: "End",
		60999: "Home",
		61010: "Insert",
		61011: "Delete",
		3639:  "PrtSc",
		70:    "ScrLk",
		3653:  "Pause",
		3677:  "Menu",
		3613:  "Right Control",
	}
	winKey2 := make(map[string]uint16)

	if runtime.GOOS == "windows" {
		for k, v := range winMap {
			winKey2[v] = k
		}
		for k1, v1 := range keyMap {
			if _, has := winKey2[v1]; !has {
				winMap[k1] = v1
			}
		}
		keyMap = winMap
	}
}
func getKeyName(code uint16) string {
	if name, has := keyMap[code]; has {
		return name
	}
	return "无"
}
