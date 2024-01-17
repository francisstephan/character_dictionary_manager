# Chinese character dictionary management
While my `Chinese character writer` is focussed on using and learning chinese characters, this dictionary manager aims at managing the dictionary (a sqlite3 database).

It is a CRUD program written in GO with GIN, HTMX and HTML5 validation.

### Watch live at https://gozdman.fly.dev/

A word about HTMX, which is an enhanced HTML enabling AJAX calls through simple tags included in HTML elements.

## Overview of the program flow:

For example, the `Size` button, when clicked, will start a `GET` AJAX request over the `/size` route, the result of which will overwrite the content (innerHTML) of the `div` element with CSS id `"content"` (further referred to as the `#content div`):

```html
<button class = "menubouton" hx-get="/size" hx-target="#content" hx-swap="innerHTML" >Size</button>
```

The `div` source code :

```html
<div id="content">{{ .content }}</div>
```

This div's content may be overwritten in two ways:
- with the GO templating system, through the {{.content}} tag
- with the htmx ajax requests, such as here above, which will overwrite the text with the dictionary's size. The /size route is handled with the `dicsize` handler :

```go
func dicsize(c *gin.Context) {
	len, time := data.Dicsize()
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("The dictionary presently contains "+len+" entries ; last updated on "+time))
}
```

We do not need to bother where this text will be displayed: this was already specified in the `hx-target = "#content"` tag of the button element hereabove.


Identifying a chinese character (`Zi => Pinyin`) is done in two steps (two AJAX requests):
### - one `GET` request to display the adequate form:

```html
<button class = "menubouton" hx-get="/getziform" hx-target="#content" hx-swap="innerHTML" >Zi => Pinyin</button>
````

The handler displays the form within the #content div:

```go
func getziform(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Ziform()))
}
```

Since the program includes not less than 7 forms, we chose to have a `forms` module, containing all forms. The Ziform function contains the following form:

```go
func Ziform() string {
	return `
	  <form hx-post="/listzi" hx-target="#content" hx-swap="innerHTML" >
		    <label for="carac">Character:</label>
		    <input id="carac" name="carac" type="text" autofocus required minlength="1" maxlength="1">
		    <button class="menubouton" type="submit">Click to submit </button>
			<button class="menubouton" hx-get="/remove">Cancel</button>
	  </form>
	`
}
```

### - one `POST` request:

When the form is submitted, it will be processed as a `POST` request by the adequate handler and the result will, once again, overwrite the #content div.

## Some problems we had :

HTMX makes for an elegant and simple program flow, using all HTTP verbs. Here we use 4 verbs: GET, POST, PUT and DELETE.

We faced 3 issues:

### 1. having a `Cancel` button in our forms: 

We tried the `htmx-remove` extension, but this did not satisfy our requirements.

We solved this by associating a specific GET request with these buttons, requesting for removal of the `#content` div's content:

```html
<button class="menubouton" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
```

and the corresponding handler :

```go
func remove(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Form canceled."))
}
```

## 2. Displaying server errors:

With HTMX (version 1.9.9) server errors are only shown in the console, which is not acceptable in production code. We found good information in https://xvello.net/blog/htmx-error-handling/ , and we used an EventListener in javascript, based on the `htmx:afterRequest` event:


```javascript
window.onload = function() {
  elem = document.body;
  document.body.addEventListener('htmx:afterRequest', function (evt) {
      contenu = document.getElementById("content");
      if (evt.detail.failed && evt.detail.xhr) {  // display error message within the document (and not only the console)
        // Server error with response contents, equivalent to htmx:responseError
        console.warn("Server error", evt.detail)
        const xhr = evt.detail.xhr;
        contenu.innerHTML = `Unexpected server error: ${xhr.status} - ${xhr.statusText}`;
      } 
    });
};```
```

## 3. Keyboard shortcuts:
We use 3 keyboard shortcuts:
- z to load the `Zi => Pinyin` form,
- p to load the `Pinyin => Zi` form, and
- Esc to abort a form.

z and p eventlisteners should only be active when `no` form is displayed, while Esc should only be active in the opposite case, i.e. when a form is displayed.
This is easily done using HTMX ajax function. For instance for the Esc key the event listener is :

```javascript
    function esckey(e) {
		if(e.keyCode==27) htmx.ajax('GET', "/remove", {target: "#content", swap: "innerHTML"}); // key esc : cancel form
	}
```

The syntax of this function call is very similar to that of the Cancel button here above.