package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gummy789j/bookings/internal/config"
	"github.com/gummy789j/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	tc := make(map[string]*template.Template)

	if app.UseCache {
		tc, _ = CreateTemplateCache()
	} else {
		tc = app.TemplateCache
	}

	t, ok := tc[tmpl]
	//fmt.Println(tmpl)
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	//parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)

	//err := parsedTemplate.Execute(w, nil)

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
	}
}

// *template.Template是一個解析過後的html(...等)的file，也就是一些儲存text的fragment的在的記憶體位置
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		// 一個Tempalte重要的包含物 Name & content
		// Template是一個定義好的struct 裡面包含 Tree struct 跟 nameSpace
		// Tree 的 型別是 *parse.Tree 是定義在parse package中，
		// 我們New一個Template，傳進去的name就是存在Tree.Name中，作為這個parse file的名字
		// 而 nameSpace對應到的就是associate file 也就是我們最尾端呼叫的 ParseFile method所傳入的我們"原本的"html file
		// 要做的就是重新建立一個Template然後幫這個template加入一個function map(讓以後擴展性更高)，再把原內容加入進去
		// 這麼做的目的有3個
		// 1.為了讓他快速讀取修改後的html file(不然都要重新run程式開socket)
		// 2.也包含將layout的定義和內容加入page中
		// 3.可以自訂義新的tempalte的function(增加靈活性)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		//fmt.Println(ts)

		// ts, err := template.ParseFiles(page)
		// if err != nil {
		// 	return myCache, err
		// }

		// matches, err := filepath.Glob("./templates/*.layout.tmpl")
		// if err != nil {
		// 	return myCache, err
		// }

		// if len(matches) > 0 {
		// 	ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
		// 	if err != nil {
		// 		return myCache, err
		// 	}
		// }

		ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts

	}

	return myCache, nil

}
