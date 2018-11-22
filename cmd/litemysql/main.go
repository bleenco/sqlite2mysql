package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	inputFile  = flag.String("inputFile", "input.sql", "input SQLite3 dump file")
	outputFile = flag.String("outputFile", "output.sql", "output MySQL dump file")
)

var (
	re1  *regexp.Regexp
	re2  *regexp.Regexp
	re3  *regexp.Regexp
	re4  *regexp.Regexp
	re5  *regexp.Regexp
	re6  *regexp.Regexp
	re7  *regexp.Regexp
	re8  *regexp.Regexp
	re9  *regexp.Regexp
	re10 *regexp.Regexp
	re11 *regexp.Regexp
	re12 *regexp.Regexp
	re13 *regexp.Regexp
)

func init() {
	re1 = regexp.MustCompile(`"`)
	re2 = regexp.MustCompile(`PRAGMA.*;`)
	re3 = regexp.MustCompile(`BEGIN TRANSACTION.*`)
	re4 = regexp.MustCompile(`.*sqlite_sequence.*;`)
	re5 = regexp.MustCompile("CREATE TABLE IF NOT EXISTS \\`(.*?)\\`")
	re6 = regexp.MustCompile(`(CREATE TABLE.*)(primary key) (autoincrement)(.*)\);`)
	re7 = regexp.MustCompile(`'t'`)
	re8 = regexp.MustCompile(`'f'`)
	re9 = regexp.MustCompile(`\b[0-9]{13}\b`)
	re10 = regexp.MustCompile(`\b[0-9]{10}\b`)
	re11 = regexp.MustCompile("([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})(.*?)[Z|\\`]")
	re12 = regexp.MustCompile(`[\+|\-][0-9]{2}:[0-9]{2}`)
	re13 = regexp.MustCompile(`[\x{2712}\x{2714}\x{2716}\x{271d}\x{2721}\x{2728}\x{2733}\x{2734}\x{2744}\x{2747}\x{274c}\x{274e}\x{2753}-\x{2755}\x{2757}\x{2763}\x{2764}\x{2795}-\x{2797}\x{27a1}\x{27b0}\x{27bf}\x{2934}\x{2935}\x{2b05}-\x{2b07}\x{2b1b}\x{2b1c}\x{2b50}\x{2b55}\x{3030}\x{303d}\x{1f004}\x{1f0cf}\x{1f170}\x{1f171}\x{1f17e}\x{1f17f}\x{1f18e}\x{1f191}-\x{1f19a}\x{1f201}\x{1f202}\x{1f21a}\x{1f22f}\x{1f232}-\x{1f23a}\x{1f250}\x{1f251}\x{1f300}-\x{1f321}\x{1f324}-\x{1f393}\x{1f396}\x{1f397}\x{1f399}-\x{1f39b}\x{1f39e}-\x{1f3f0}\x{1f3f3}-\x{1f3f5}\x{1f3f7}-\x{1f4fd}\x{1f4ff}-\x{1f53d}\x{1f549}-\x{1f54e}\x{1f550}-\x{1f567}\x{1f56f}\x{1f570}\x{1f573}-\x{1f579}\x{1f587}\x{1f58a}-\x{1f58d}\x{1f590}\x{1f595}\x{1f596}\x{1f5a5}\x{1f5a8}\x{1f5b1}\x{1f5b2}\x{1f5bc}\x{1f5c2}-\x{1f5c4}\x{1f5d1}-\x{1f5d3}\x{1f5dc}-\x{1f5de}\x{1f5e1}\x{1f5e3}\x{1f5ef}\x{1f5f3}\x{1f5fa}-\x{1f64f}\x{1f680}-\x{1f6c5}\x{1f6cb}-\x{1f6d0}\x{1f6e0}-\x{1f6e5}\x{1f6e9}\x{1f6eb}\x{1f6ec}\x{1f6f0}\x{1f6f3}\x{1f910}-\x{1f918}\x{1f980}-\x{1f984}\x{1f9c0}\x{3297}\x{3299}\x{a9}\x{ae}\x{203c}\x{2049}\x{2122}\x{2139}\x{2194}-\x{2199}\x{21a9}\x{21aa}\x{231a}\x{231b}\x{2328}\x{2388}\x{23cf}\x{23e9}-\x{23f3}\x{23f8}-\x{23fa}\x{24c2}\x{25aa}\x{25ab}\x{25b6}\x{25c0}\x{25fb}-\x{25fe}\x{2600}-\x{2604}\x{260e}\x{2611}\x{2614}\x{2615}\x{2618}\x{261d}\x{2620}\x{2622}\x{2623}\x{2626}\x{262a}\x{262e}\x{262f}\x{2638}-\x{263a}\x{2648}-\x{2653}\x{2660}\x{2663}\x{2665}\x{2666}\x{2668}\x{267b}\x{267f}\x{2692}-\x{2694}\x{2696}\x{2697}\x{2699}\x{269b}\x{269c}\x{26a0}\x{26a1}\x{26aa}\x{26ab}\x{26b0}\x{26b1}\x{26bd}\x{26be}\x{26c4}\x{26c5}\x{26c8}\x{26ce}\x{26cf}\x{26d1}\x{26d3}\x{26d4}\x{26e9}\x{26ea}\x{26f0}-\x{26f5}\x{26f7}-\x{26fa}\x{26fd}\x{2702}\x{2705}\x{2708}-\x{270d}\x{270f}]|\x{23}\x{20e3}|\x{2a}\x{20e3}|\x{30}\x{20e3}|\x{31}\x{20e3}|\x{32}\x{20e3}|\x{33}\x{20e3}|\x{34}\x{20e3}|\x{35}\x{20e3}|\x{36}\x{20e3}|\x{37}\x{20e3}|\x{38}\x{20e3}|\x{39}\x{20e3}|\x{1f1e6}[\x{1f1e8}-\x{1f1ec}\x{1f1ee}\x{1f1f1}\x{1f1f2}\x{1f1f4}\x{1f1f6}-\x{1f1fa}\x{1f1fc}\x{1f1fd}\x{1f1ff}]|\x{1f1e7}[\x{1f1e6}\x{1f1e7}\x{1f1e9}-\x{1f1ef}\x{1f1f1}-\x{1f1f4}\x{1f1f6}-\x{1f1f9}\x{1f1fb}\x{1f1fc}\x{1f1fe}\x{1f1ff}]|\x{1f1e8}[\x{1f1e6}\x{1f1e8}\x{1f1e9}\x{1f1eb}-\x{1f1ee}\x{1f1f0}-\x{1f1f5}\x{1f1f7}\x{1f1fa}-\x{1f1ff}]|\x{1f1e9}[\x{1f1ea}\x{1f1ec}\x{1f1ef}\x{1f1f0}\x{1f1f2}\x{1f1f4}\x{1f1ff}]|\x{1f1ea}[\x{1f1e6}\x{1f1e8}\x{1f1ea}\x{1f1ec}\x{1f1ed}\x{1f1f7}-\x{1f1fa}]|\x{1f1eb}[\x{1f1ee}-\x{1f1f0}\x{1f1f2}\x{1f1f4}\x{1f1f7}]|\x{1f1ec}[\x{1f1e6}\x{1f1e7}\x{1f1e9}-\x{1f1ee}\x{1f1f1}-\x{1f1f3}\x{1f1f5}-\x{1f1fa}\x{1f1fc}\x{1f1fe}]|\x{1f1ed}[\x{1f1f0}\x{1f1f2}\x{1f1f3}\x{1f1f7}\x{1f1f9}\x{1f1fa}]|\x{1f1ee}[\x{1f1e8}-\x{1f1ea}\x{1f1f1}-\x{1f1f4}\x{1f1f6}-\x{1f1f9}]|\x{1f1ef}[\x{1f1ea}\x{1f1f2}\x{1f1f4}\x{1f1f5}]|\x{1f1f0}[\x{1f1ea}\x{1f1ec}-\x{1f1ee}\x{1f1f2}\x{1f1f3}\x{1f1f5}\x{1f1f7}\x{1f1fc}\x{1f1fe}\x{1f1ff}]|\x{1f1f1}[\x{1f1e6}-\x{1f1e8}\x{1f1ee}\x{1f1f0}\x{1f1f7}-\x{1f1fb}\x{1f1fe}]|\x{1f1f2}[\x{1f1e6}\x{1f1e8}-\x{1f1ed}\x{1f1f0}-\x{1f1ff}]|\x{1f1f3}[\x{1f1e6}\x{1f1e8}\x{1f1ea}-\x{1f1ec}\x{1f1ee}\x{1f1f1}\x{1f1f4}\x{1f1f5}\x{1f1f7}\x{1f1fa}\x{1f1ff}]|\x{1f1f4}\x{1f1f2}|\x{1f1f5}[\x{1f1e6}\x{1f1ea}-\x{1f1ed}\x{1f1f0}-\x{1f1f3}\x{1f1f7}-\x{1f1f9}\x{1f1fc}\x{1f1fe}]|\x{1f1f6}\x{1f1e6}|\x{1f1f7}[\x{1f1ea}\x{1f1f4}\x{1f1f8}\x{1f1fa}\x{1f1fc}]|\x{1f1f8}[\x{1f1e6}-\x{1f1ea}\x{1f1ec}-\x{1f1f4}\x{1f1f7}-\x{1f1f9}\x{1f1fb}\x{1f1fd}-\x{1f1ff}]|\x{1f1f9}[\x{1f1e6}\x{1f1e8}\x{1f1e9}\x{1f1eb}-\x{1f1ed}\x{1f1ef}-\x{1f1f4}\x{1f1f7}\x{1f1f9}\x{1f1fb}\x{1f1fc}\x{1f1ff}]|\x{1f1fa}[\x{1f1e6}\x{1f1ec}\x{1f1f2}\x{1f1f8}\x{1f1fe}\x{1f1ff}]|\x{1f1fb}[\x{1f1e6}\x{1f1e8}\x{1f1ea}\x{1f1ec}\x{1f1ee}\x{1f1f3}\x{1f1fa}]|\x{1f1fc}[\x{1f1eb}\x{1f1f8}]|\x{1f1fd}\x{1f1f0}|\x{1f1fe}[\x{1f1ea}\x{1f1f9}]|\x{1f1ff}[\x{1f1e6}\x{1f1f2}\x{1f1fc}]`)
}

func main() {
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	reader := bufio.NewReader(file)
	var i int32

	if _, err := io.WriteString(output, `SET FOREIGN_KEY_CHECKS=0;\n`); err != nil {
		log.Fatal(err)
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		result := re1.ReplaceAllString(line, "`")
		result = re2.ReplaceAllString(result, ``)
		result = re3.ReplaceAllString(result, ``)
		result = re4.ReplaceAllString(result, ``)
		result = re5.ReplaceAllString(result, "DROP TABLE IF EXISTS `$1`;\nCREATE TABLE `$1`")
		result = re6.ReplaceAllString(result, "$1 AUTO_INCREMENT $4, PRIMARY KEY(id));")
		result = re7.ReplaceAllString(result, `1`)
		result = re8.ReplaceAllString(result, `0`)
		result = re9.ReplaceAllStringFunc(result, func(m string) string {
			return `'` + formatDate(m) + `'`
		})
		result = re10.ReplaceAllStringFunc(result, func(m string) string {
			return `'` + formatDate(m) + `'`
		})
		result = re11.ReplaceAllStringFunc(result, func(m string) string {
			return `'` + parseDate(m) + `'`
		})
		result = re13.ReplaceAllString(result, ``)

		if result != "\n" {
			if _, err := io.WriteString(output, result); err != nil {
				log.Fatal(err)
			}
		}

		i++
	}

	if _, err := io.WriteString(output, `SET FOREIGN_KEY_CHECKS=1;\n`); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func formatDate(unixTimestamp string) string {
	i, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		panic(err)
	}

	if len(unixTimestamp) == 13 {
		i = i / 1000
	}

	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func parseDate(d string) string {
	layout := "2006-01-02T15:04:05.9Z"

	d = re12.ReplaceAllString(d, "")

	last := d[len(d)-1:]
	if last == "`" {
		d = strings.TrimSuffix(d, "`")
	}
	last = d[len(d)-1:]
	if last != "Z" {
		d = d + `Z`
	}
	d = strings.Replace(d, "\\", "", -1)

	t, err := time.Parse(layout, d)
	if err != nil {
		log.Fatal(err)
	}

	return t.Format("2006-01-02 15:04:05")
}
