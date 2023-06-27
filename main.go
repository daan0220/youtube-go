package main

import (
	"log"
	"net/http"
	"os/user"
	"strconv"
	"strings"
	"text/template"

	"github.com/daan0220/youtube-go/my"
	"github.com/gorilla/sessions"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//db variable
var dbDriver = "sqlite3"
var dbName = "data.sqlite"

//session variable
var sesName = "ytboard-session"
var cs = sessions.NewCookieStore([]byte("secret-key-1234"))

//login check
func checkLogin(w http.ResponseWriter, rq *http.Request) *my.User {
	ses, _ := cs.Get(rq, sesName)
	if ses.Values["login"] == nil || !ses.Values["login"].(bool) {
		http.Redirect(w, rq, "/login", 302)
	}
	ac :=""
	if ses.Values["account"] != nil {
		ac = ses.Values["account"].(string)
	}

	var user my.User
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	db.Where("account = ?", ac).First(&user)
	return &user
}

//Template for no-template.
func notemp() *template.Template {
	tmp, _ := template.New("index").Parse("NO PAGE.")
		return tmp
}

//get target Temlate
func page(fname string) *template.Template {
	tmps, _ := template.ParseFiles("templates/"+fname+".html",
			"templates/head.html", "templates/foot.html")
			return tmps
}
func createTemplate() *template.Template {
	tmp := template.New("index")
	tmp, _ = tmp.Parse("NO PAGE.")
	return tmp
}

//top page handler.
func index(w http.ResponseWriter, rq * http.Request) {
	user := checkLogin(w, rq)

	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	var pl []my.Post
	db.Where("group_id > 0").Order("created_at desc").Limit(10).Find(&pl)
	var gl []my.Group
	db.Order("created_at desc").Limit(10).Find(&gl)

	item := struct {
		Title string
		Message string
		Name string
		Account string
		Plist []my.Post
		Glist []my.Group
	}{
		Title: "Index",
		Message: "This is Top page",
		Name: user.Name,
		Account: user.Account,
		Plist: pl,
		Glist: gl,
	}
	er := page("index").Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

func Post(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)
    
	pid := rq.FormValue("pid")
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if rq.Method == "POST" {
		msg := rq.PostFormValue("message")
		pId, _ := strconv.Atoi(pid)
		cmt := my.Comment{
			UserId: int(user.Model.ID),
			PostId: pId,
			Message: msg,
		}
		db.Create(&cmt)
	}
	var pst my.Post
	var cmts []my.CommentJoin

	db.Where("id = ?", pid).First(&pst)
	db.Table("comments").Select("comments.*, users.id, users.name")
	.Joins("join users on users.id =comments.user_id")
	.Where("comments.post_id = ?", pid).Order("created_at desc").Find(&cmts)

	item := struct {
		Title string
		Message string
		Name string
		Account string
		Post my.Post
		Clist []my.CommentJoin
	}{
		Title: "Post",
		Message: "Post id=" + pid,
		Name: user.Name,
		Account: user.Account,
		Post: pst,
		Clist: cmts,
	}
	er := page("post").Execute(w, item)
	if er != nil {
		log.Fatal(er)
	}
}

//home handler
func home(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)
    
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if rq.Method == "POST" {
		switch rq.PostFormValue("form") {
		case "post":
			ad := rq.PostFormValue("address")
			ad := strings.TrimSpace(ad)
			if strings.HasPrefix(ad, "https://youtu.be/") {
				ad = strings.TrimPrefix(ad, "https://youtu.be/")
			}

			pt := my.Post{
				UserId: int(user.Model.ID),
				Address: ad,
				Message: rq.PostFormValue("message"),
			}
			db.Create(&pt)
		case "group":
			gp := my.Group{
				UserId: int(user.Model.ID),
				Name: rq.PostFormValue("name"),
				Message: rq.PostFormValue("message"),
			}
			db.Create(&gp)
		}
		
	}
	var pts []my.Post
	var gps []my.Group

	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&pts)
	db.Where("user_id=?", user.ID).Order("created_at desc").Limit(10).Find(&gps)
		

	itm := struct {
		Title string
		Message string
		Name string
		Account string
		Plist my.Post
		Glist []my.Group
	}{
		Title: "Home",
		Message: "User account=\"" + user.Account + "\".",
		Name: user.Name,
		Account: user.Account,
		Plist: pts,
		Glist: gps,
	}
	er := page("home").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}


//group handler
func group(w http.ResponseWriter, rq *http.Request) {
	user := checkLogin(w, rq)
    
	gid := rq.FormValue("gid")
	db, _ := gorm.Open(dbDriver, dbName)
	defer db.Close()

	if rq.Method == "POST" {
		switch rq.PostFormValue("form") {
		case "post":
			ad := rq.PostFormValue("address")
			ad := strings.TrimSpace(ad)
			if strings.HasPrefix(ad, "https://youtu.be/") {
				ad = strings.TrimPrefix(ad, "https://youtu.be/")
			}

			gid, _ := strconv.Atoi(gid)
			pt := my.Post{
				UserId: int(user.Model.ID),
				Address: ad,
				Message: rq.PostFormValue("message"),
			}
			db.Create(&pt)
		}
		
	}
	var grp my.Group
	var pts []my.Post

	db.Where("id=?", gid).First(&grp)
	db.Order("created_at desc").Model(&grp).Related(&pts)
		

	itm := struct {
		Title string
		Message string
		Name string
		Account string
		Group my.Group
		Plist []my.Post
	}{
		Title: "Group",
		Message: "Group id=" + gid,
		Name: user.Name,
		Account: user.Account,
		Group: grp,
		Plist: pts,
	}
	er := page("group").Execute(w, itm)
	if er != nil {
		log.Fatal(er)
	}
}

//login handler


//logout handler





// main program
func main() {
	my.Migrate()
}
