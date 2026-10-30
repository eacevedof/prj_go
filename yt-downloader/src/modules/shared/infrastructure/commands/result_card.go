package commands

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"yt-downloader-go/src/core/boot"
)

// printCard el resultado (map anidado) como tarjeta para humanos: "clave: valor" por linea, mapas anidados con
// sangria, listas de escalares en una linea y listas de mapas como "N items". Las rutas (file_path, output_path) van
// como hiperenlace OSC 8 (clic en Windows Terminal) y en amarillo.
func (c *AbstractConsoleCommand) printCard(result any) {
	fmt.Println()
	c.printCardValue(reflect.ValueOf(result), 0)
}

func (c *AbstractConsoleCommand) printCardValue(value reflect.Value, depth int) {
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Map {
		fmt.Printf("%s%v\n", getIndent(depth), value.Interface())
		return
	}
	keys := make([]string, 0, value.Len())
	for _, key := range value.MapKeys() {
		keys = append(keys, fmt.Sprint(key.Interface()))
	}
	sort.Strings(keys)
	width := 0
	for _, key := range keys {
		if len(key) > width {
			width = len(key)
		}
	}
	for _, key := range keys {
		item := value.MapIndex(reflect.ValueOf(key))
		for item.Kind() == reflect.Interface && !item.IsNil() {
			item = item.Elem()
		}
		label := getIndent(depth) + c.GetGreen(fmt.Sprintf("%-*s", width, key)) + "  "
		switch item.Kind() {
		case reflect.Map:
			fmt.Println(label)
			c.printCardValue(item, depth+1)
		case reflect.Slice, reflect.Array:
			fmt.Println(label + c.getCardList(item))
		default:
			fmt.Println(label + c.getCardScalar(key, fmt.Sprint(item.Interface())))
		}
	}
}

func (c *AbstractConsoleCommand) getCardList(list reflect.Value) string {
	if list.Len() == 0 {
		return "-"
	}
	first := list.Index(0)
	for first.Kind() == reflect.Interface && !first.IsNil() {
		first = first.Elem()
	}
	if first.Kind() == reflect.Map || first.Kind() == reflect.Slice {
		return fmt.Sprintf("%d items", list.Len())
	}
	parts := make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		parts = append(parts, fmt.Sprint(list.Index(i).Interface()))
	}
	return strings.Join(parts, ", ")
}

// getCardScalar las rutas como enlace file:// clicable (OSC 8) y en amarillo; el resto tal cual.
func (c *AbstractConsoleCommand) getCardScalar(key string, text string) string {
	if !strings.HasSuffix(key, "_path") || text == "" {
		return text
	}
	return c.GetYellow(getHyperlink("file:///"+filepath.ToSlash(text), text))
}

// getHyperlink OSC 8: ESC ] 8 ; ; url ESC \ texto ESC ] 8 ; ; ESC \ . Solo si la terminal interpreta ANSI.
func getHyperlink(url string, text string) string {
	if !boot.AreColorsEnabled() {
		return text
	}
	return "\033]8;;" + url + "\033\\" + text + "\033]8;;\033\\"
}

func getIndent(depth int) string {
	return strings.Repeat("    ", depth)
}
