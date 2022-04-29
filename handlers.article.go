// handlers.article.go

package main

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
  articles := getAllArticles()

  // Call the render function with the name of the template to render
  render(c, gin.H{
    "title":   "Home Page",
    "payload": articles}, "index.html")

}

func getArticle(c *gin.Context) {
  // Check if the article ID is valid
  if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
    // Check if the article exists
    if article, err := getArticleByID(articleID); err == nil {
      // Call the HTML method of the Context to render a template
      c.HTML(
        // Set the HTTP status to 200 (OK)
        http.StatusOK,
        // Use the article.html template
        "article.html",
        // Pass the data that the page uses
        gin.H{
          "title":   article.Title,
          "payload": article,
        },
      )

    } else {
      // If the article is not found, abort with an error
      c.AbortWithError(http.StatusNotFound, err)
    }

  } else {
    // If an invalid article ID is specified in the URL, abort with an error
    c.AbortWithStatus(http.StatusNotFound)
  }
  
}
func showArticleCreationPage(c *gin.Context) {
    render(c, gin.H{
        "title": "Create New Article"}, "create-article.html")
}

func createArticle(c *gin.Context) {
    title := c.PostForm("title")
    content := c.PostForm("content")

    if a, err := createNewArticle(title, content); err == nil {
        render(c, gin.H{
            "title":   "Submission Successful",
            "payload": a}, "submission-successful.html")
    } else {
        c.AbortWithStatus(http.StatusBadRequest)
    }
}