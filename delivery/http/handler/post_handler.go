package handler

import (
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amthesonofGod/Notice-Board/form"

	// "encoding/json"

	"github.com/amthesonofGod/Notice-Board/company"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
)

// CompanyPostHandler handles post handler admin requests
type CompanyPostHandler struct {
	tmpl           *template.Template
	postSrv        post.PostService
	companySrv     company.CompanyService
	sessionService company.SessionServiceCamp
	campSess       *entity.CompanySession
}

// NewCompanyPostHandler initializes and returns new CompanyPostHandler
func NewCompanyPostHandler(T *template.Template, PS post.PostService, CP company.CompanyService) *CompanyPostHandler {
	return &CompanyPostHandler{tmpl: T, postSrv: PS, companySrv: CP}
}

// CompanyPosts handle requests on route /admin/posts
func (cph *CompanyPostHandler) CompanyPosts(w http.ResponseWriter, r *http.Request) {

	handler := CompanyHandler{loggedInUserCamp: currentCompUser}

	authorizedPost := []entity.Post{}

	posts, errs := cph.postSrv.Posts()

	if len(errs) > 0 {
		panic(errs)
	}
	for _, post := range posts {
		if handler.loggedInUserCamp.ID == post.CompanyID {
			authorizedPost = append(authorizedPost, post)
		}
	}

	fmt.Println("Current Post")
	fmt.Println(authorizedPost)
	cph.tmpl.ExecuteTemplate(w, "cmp_post.layout", authorizedPost)
}

// CompanyPostsNew hanlde requests on route /admin/posts/new
func (cph *CompanyPostHandler) CompanyPostsNew(w http.ResponseWriter, r *http.Request) {

	fmt.Println("companypostsnew function invoked! ")

	if r.Method == http.MethodPost {

		handler := CompanyHandler{loggedInUserCamp: currentCompUser}

		compID := handler.loggedInUserCamp.ID

		cmp, cerr := cph.companySrv.Company(compID)

		if len(cerr) > 0 {
			fmt.Println("i am the error")
			panic(cerr)
		}

		fmt.Println(cmp.Name)

		post := &entity.Post{}
		post.CompanyID = compID
		post.Owner = cmp.Name
		post.Title = r.FormValue("title")
		post.Description = r.FormValue("description")
		post.Category = r.Form.Get("category")

		fmt.Println(post.Category)

		mf, fh, err := r.FormFile("postimg")
		if err != nil {
			panic(err)
		}
		defer mf.Close()

		post.Image = fh.Filename

		writeFile(&mf, fh.Filename)

		// Validate the form contents
		singnUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		singnUpForm.Required("title", "description", "category")
		// If there are any errors, redisplay the signup form.
		if !singnUpForm.Valid() {
			cph.tmpl.ExecuteTemplate(w, "post-job.layout", singnUpForm)
			return
		}

		_, errs := cph.postSrv.StorePost(post)
		// cph.postSrv.StorePost(post)

		if len(errs) > 0 {
			panic(errs)
		}
		// fmt.Println(entity.Company.ID)
		fmt.Println(post)
		fmt.Println("post added to db")

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)

	} else {

		cph.tmpl.ExecuteTemplate(w, "post-job.layout", nil)

	}
}

// CompanyPostUpdate handle requests on /cmp/posts/update
func (cph *CompanyPostHandler) CompanyPostUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		post, errs := cph.postSrv.Post(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

		cph.tmpl.ExecuteTemplate(w, "post_update.layout", post)

	}
	if r.Method == http.MethodPost {

		pst := &entity.Post{}

		postid, _ := strconv.Atoi(r.FormValue("id"))

		companyid, _ := strconv.Atoi(r.FormValue("companyid"))

		pst.ID = uint(postid)

		pst.CompanyID = uint(companyid)

		pst.Owner = r.FormValue("owner")

		pst.Title = r.FormValue("title")
		pst.Description = r.FormValue("description")
		pst.Image = r.FormValue("image")
		pst.Category = r.Form.Get("category")

		mf, fh, err := r.FormFile("postimg")

		if err != nil {
			panic(err)
		}

		defer mf.Close()

		pst.Image = fh.Filename

		writeFile(&mf, pst.Image)

		_, errs := cph.postSrv.UpdatePost(pst)

		if len(errs) > 0 {
			panic(errs)
		}

		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)

	}

}

// CompanyPostDelete handle requests on route /cmp/posts/delete
func (cph *CompanyPostHandler) CompanyPostDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		idRaw := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		_, errs := cph.postSrv.DeletePost(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}

	}

	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func writeFile(mf *multipart.File, fname string) {

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	path := filepath.Join(wd, "../../", "ui", "assets", "img", fname)
	image, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	io.Copy(image, *mf)
}
