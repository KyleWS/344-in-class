# Web Client for your Zip Server

Create a new web client in this directory that makes AJAX requests to the server you wrote in the `../zipsvr/` directory. Follow these steps:

- Start your zip server if it's not already running.
- Create a web page that has at least one `<input>` into which a user can type a city name. 
- You may style this page in any way you want. You can use your own style rules or import a CSS framework like [Bootstrap](https://getbootstrap.com/).
- When the user hits Enter or clicks/taps a submit button, make an AJAX request to your zip server to get the states and zips for that city. If you wrap the `<input>` in a `<form>` element, remember to prevent the default behavior so that the browser doesn't try to submit the form and refresh the page.
- Display the results on the web page by creating new elements, populating them with the data, and adding them to the page. Remember to clear the previous results before showing new ones.
- If you get an error status code (e.g., 404), show the error response text to the user by populating some sort of error alert element on the page.
- Extra kudos if you make this city name `<input>` part of a larger `<form>` with `<select>` elements for the possible states and zips. Do the AJAX request during the `"change"` event on the city `<input>` (which fires when the input looses focus), and add `<option>` elements to the state and zip `<select>` elements for the distinct states and zips that are returned.

Bragging rights go to the student(s) who finish this first!
