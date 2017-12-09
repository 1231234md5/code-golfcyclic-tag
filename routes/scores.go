package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var holes = [][]string{
	{"all", "All Holes"},
	{"12-days-of-christmas", "12 Days of Christmas"},
	{"99-bottles-of-beer", "99 Bottles of Beer"},
	{"arabic-to-roman", "Arabic to Roman"},
	{"christmas-trees", "Christmas Trees"},
	{"emirp-numbers", "Emirp Numbers"},
	{"evil-numbers", "Evil Numbers"},
	{"fibonacci", "Fibonacci"},
	{"fizz-buzz", "Fizz Buzz"},
	{"happy-numbers", "Happy Numbers"},
	{"odious-numbers", "Odious Numbers"},
	{"pangram-grep", "Pangram Grep"},
	{"pascals-triangle", "Pascal's Triangle"},
	{"pernicious-numbers", "Pernicious Numbers"},
	{"prime-numbers", "Prime Numbers"},
	{"quine", "Quine"},
	{"roman-to-arabic", "Roman to Arabic"},
	{"sierpiński-triangle", "Sierpiński Triangle"},
	{"seven-segment", "Seven Segment"},
	{"spelling-numbers", "Spelling Numbers"},
	{"π", "π"},
	{"φ", "φ"},
	{"𝑒", "𝑒"},
	{"τ", "τ"},
}

var langs = [][]string{
	{"all", "All Langs"},
	{"bash", "Bash"},
	{"javascript", "JavaScript"},
	{"perl", "Perl"},
	{"perl6", "Perl 6"},
	{"php", "PHP"},
	{"python", "Python"},
	{"ruby", "Ruby"},
}

func scores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	parts := strings.Split(r.URL.Path, "/")

	hole := parts[2]
	lang := parts[3]

	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<script async src=" + scoresJsPath +
			"></script><main id=scores><select id=hole>",
	))

	for _, v := range holes {
		w.Write([]byte("<option "))
		if hole == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><select id=lang>"))

	for _, v := range langs {
		w.Write([]byte("<option "))
		if lang == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><table class=scores>"))

	var concat, where string

	if hole != "all" {
		where += " AND hole = '" + hole + "'"
		concat = "' class=', lang, '>', TO_CHAR(strokes, 'FM99,999'), '<td>(', TO_CHAR(score, 'FM9,999'), ' point', CASE WHEN score"
	} else {
		concat = "'>', TO_CHAR(score, 'FM9,999'), '<td>(', count, ' hole', CASE WHEN count"
	}

	if lang != "all" {
		where += " AND lang = '" + lang + "'"
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id,
		         lang
		    FROM solutions
		   WHERE true`+where+`
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 COUNT(*) OVER (PARTITION BY hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 COUNT(*) OVER (PARTITION BY hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id,
		         lang
		    FROM leaderboard
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(strokes) strokes,
		         SUM(score)   score,
		         COUNT(*),
		         MAX(submitted),
		         STRING_AGG(CONCAT(lang), '') lang
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             '<tr',
		             CASE WHEN user_id = $1 THEN ' class=me' END,
		             '><td>',
		             TO_CHAR(
		                 RANK() OVER (ORDER BY score DESC),
		                 'FM999"<sup>"th"</sup>"'
		             ),
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="/users/',
		             login,
		             '">',
		             login,
		             '</a><td',
		             `+concat+` > 1 THEN 's' END,
		             ')<td><time datetime=',
		             TO_CHAR(max, 'YYYY-MM-DD"T"HH24:MI:SS"Z>"FMDD Mon'),
		             '</time>'
		         )
		    FROM summed_leaderboard
		    JOIN users on user_id = id
		ORDER BY score DESC, max`,
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var row sql.RawBytes

		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		w.Write(row)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
