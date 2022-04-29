// models.article_test.go

package main

import (
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  "strconv"
  "net/url"
)

// Test the function that fetches all articles
func TestGetAllArticles(t *testing.T) {
  alist := getAllArticles()

  // Check that the length of the list of articles returned is the
  // same as the length of the global variable holding the list
  if len(alist) != len(articleList) {
    t.Fail()
  }

  // Check that each member is identical
  for i, v := range alist {
    if v.Content != articleList[i].Content ||
      v.ID != articleList[i].ID ||
      v.Title != articleList[i].Title {

      t.Fail()
      break
    }
  }
}

func TestShowIndexPageUnauthenticated(t *testing.T) {
  r := getRouter(true)

  r.GET("/", showIndexPage)

  // Create a request to send to the above route
  req, _ := http.NewRequest("GET", "/", nil)

  testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
    // Test that the http status code is 200
    statusOK := w.Code == http.StatusOK

    // Test that the page title is "Home Page"
    // You can carry out a lot more detailed tests using libraries that can
    // parse and process HTML pages
    p, err := ioutil.ReadAll(w.Body)
    pageOK := err == nil && strings.Index(string(p), "<title>Home Page</title>") > 0

    return statusOK && pageOK
  })
}

func TestCreateNewArticle(t *testing.T) {
    originalLength := len(getAllArticles())

    a, err := createNewArticle("New test title", "New test content")

    allArticles := getAllArticles()
    newLength := len(allArticles)

    if err != nil || newLength != originalLength+1 ||
        a.Title != "New test title" || a.Content != "New test content" {

        t.Fail()
    }
}

func TestArticleCreationAuthenticated(t *testing.T) {
    saveLists()
    w := httptest.NewRecorder()

    r := getRouter(true)

    http.SetCookie(w, &http.Cookie{Name: "token", Value: "123"})

    r.POST("/article/create", createArticle)

    articlePayload := getArticlePOSTPayload()
    req, _ := http.NewRequest("POST", "/article/create", strings.NewReader(articlePayload))
    req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("Content-Length", strconv.Itoa(len(articlePayload)))

    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fail()
    }

    p, err := ioutil.ReadAll(w.Body)
    if err != nil || strings.Index(string(p), "<title>Submission Successful</title>") < 0 {
        t.Fail()
    }
    restoreLists()
}

func getArticlePOSTPayload() string {
    params := url.Values{}
    params.Add("title", "Test Article Title")
    params.Add("content", "Test Article Content")

    return params.Encode()
}