package forms

import (
	"gozdman/data"
	"strconv"
)

// html form validation : cf. https://developer.mozilla.org/fr/docs/Learn/Forms/Form_validation

func Ziform() string {
	return `
	  <form hx-post="/listzi" hx-target="#content" hx-swap="innerHTML" >
		    <label for="carac">Character:</label>
		    <input id="carac" name="carac" type="text" autofocus required minlength="1" maxlength="1">
		    <button class="menubouton" type="submit">Click to submit </button>
			<button class="menubouton" hx-get="/remove" hx-target="#content", hx-swap="innerHTML">Cancel</button>
	  </form>
	`
}

func Pyform() string {
	return `
	  <form hx-post="/listpy" hx-target="#content" hx-swap="innerHTML" >
		    <label for="pinyin">Pinyin+tone (using pattern ^[a-z,ü]+[0-4]?) :</label>
		    <input id="pinyin" name="pinyin" type="text" pattern="^[a-z,ü]+[0-4]?" autofocus>
		    <button class="menubouton" type="submit">Click to submit </button>
			<button class="menubouton" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	`
}

func Selupdate() string {
	return `
	  <form hx-get="/updatezi" hx-target="#content" hx-swap="innerHTML" >
		    <label for="id">Id:</label>
		    <input id="id" name="id" required type="number" autofocus>
		    <button class="menubouton" type="submit">Click to submit </button>
			<button class="menubouton" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	  `
}

func Addziform() string {
	// Cancel button : cf. https://alchemists.io/projects/htmx-remove
	return `
	  <form id="addziform" class="multi" hx-post="/addzi" hx-target="#content" hx-swap="innerHTML">
	      <p id="formhead">Add character to dictionary :</p>
		  <label for="pinyin_ton">Pinyin+tone (using pattern ^[a-z,ü]+[0-4]?) :</label>
		  <input id="pinyin_ton" name="pinyin_ton" type="text" pattern="^[a-z,ü]+[0-4]?" autofocus><br />
		  <label for="unicode">Unicode:</label>
		  <input id="unicode" name="unicode" type="text"><span id="viewcar"> </span><br />
		  <label for="sens">Meaning:</label>
		  <input id="sens" name="sens" type="text"><br />
		
		<button class="formbut" type="submit">Submit</button>
		<button class="formbut" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	  <script>
	    function displayChar(){
			var s= document.getElementById("unicode");
			document.getElementById("viewcar").innerHTML = String.fromCharCode(parseInt(s.value,16))	;
		}
	    document.getElementById("unicode").addEventListener("keyup", displayChar);
		document.getElementById("unicode").addEventListener("change", displayChar);
	  </script>
	`
}

func Updateziform(zi data.Zi) string {

	retour := `
	  <form class="multi" hx-put="/doupdate/` + strconv.Itoa(zi.Id) + `" hx-target="#content" hx-swap="innerHTML">
	    <p id="formhead">Update dictionary entry :</p>
	    <label for="Id">Id:</label>
	    <input id="Id" name="Id" type="text" readonly value="` + strconv.Itoa(zi.Id) + `"><br />
	    <label for="pinyin_ton">Pinyin+tone (using pattern ^[a-z,ü]+[0-4]?) :</label>
		<input id="pinyin_ton" name="pinyin_ton" type="text" pattern="^[a-z,ü]+[0-4]?" value="` + zi.Pinyin_ton + `"><br />
		<label for="unicode">Unicode:</label>
		<input id="unicode" name="unicode" type="text" value="` + zi.Unicode + `"><span id="viewcar"> ` + zi.Hanzi + ` </span><br />
		<label for="sens">Meaning:</label>
		<input id="sens" name="sens" type="text" value="` + zi.Sens + `"><br />
		<button class="formbut" type="submit">Submit</button>
		<button class="formbut" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	  <script>
	    function displayChar(){
			var s= document.getElementById("unicode");
			document.getElementById("viewcar").innerHTML = String.fromCharCode(parseInt(s.value,16))	;
		}
	    document.getElementById("unicode").addEventListener("keyup", displayChar);
		document.getElementById("unicode").addEventListener("change", displayChar);
	  </script>
	`
	return retour
}

func Seldelete() string {
	return `
	  <form hx-post="/seldelete" hx-target="#content" hx-swap="innerHTML" >
		    <p>Enter Id of entry to delete :</p>
			<label for="id">Id:</label>
		    <input id="id" name="id" required type="number" autofocus>
		    <button class="menubouton" type="submit">Submit</button>
			<button class="menubouton" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	  `
}

func Confdelete(zi data.Zi) string {
	id := strconv.Itoa(zi.Id)
	return `
	  <form class="multi" hx-delete="/delete/` + id + `" hx-target="#content" hx-swap="innerHTML">
		    <p>Confirm delete of entry:</p>
			<table>` + data.Printzi(zi) + `</table>
		    
			<button class="formbut" type="submit">Yes</button>
		    <button class="formbut" hx-get="/remove" hx-target="#content" hx-swap="innerHTML">Cancel</button>
	  </form>
	  `
}
