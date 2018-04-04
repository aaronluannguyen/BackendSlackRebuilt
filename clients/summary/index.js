// @ts-check
"use strict";

let baseURL = "http://localhost:4000";
let queryCall = "/v1/summary?url=";
let website = "";
let webSearch = document.querySelector("input");
let result = document.querySelector("#result");
let alertDiv = document.querySelector("#alert-box");

document.querySelector("form")
    .addEventListener("submit", (evt) => {
       evt.preventDefault();
       website = webSearch.value;

       fetch(baseURL + queryCall + website)
           .then((response) => {
               alertDiv.innerHTML = "";
               result.innerHTML = "";
               if (response.ok) {
                   return response.json()
               }
           })
           .then((data) => {
               let link = document.createElement("a");
               result.appendChild(link);

               let card = document.createElement("div");
               card.classList.add("card");
               link.appendChild(card);

               let imageDiv = document.createElement("div");
               imageDiv.classList.add("imgDiv");
               card.appendChild(imageDiv);

               for (let i = 0; i < data.images.length; i++) {
                   let img = document.createElement("img");
                   img.classList.add("card-img-top");
                   img.src = data.images[i].url;
                   imageDiv.appendChild(img);
               }

               let cardBody = document.createElement("div");
               cardBody.classList.add("card-body");
               card.appendChild(cardBody);

               let title = document.createElement("h2");
               title.classList.add("card-title");
               title.textContent = data.title;
               cardBody.appendChild(title);

               let description = document.createElement("p");
               description.classList.add("card-text");
               description.textContent = data.description;
               cardBody.appendChild(description)
           })
           .catch(err => {
               let alert = document.createElement("div");
               alert.classList.add("alert");
               alert.classList.add("alert-danger");
               alertDiv.appendChild(alert);

               let alertTitle = document.createElement("h");
               alertTitle.textContent = err;

               alert.appendChild(alertTitle);
           })
    });