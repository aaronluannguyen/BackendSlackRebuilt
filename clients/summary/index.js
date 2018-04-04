// @ts-check
"use strict";

let baseURL = "http://localhost:4000";
let queryCall = "/v1/summary?url=";
let website = "";
let webSearch = document.querySelector("input");

document.querySelector("form")
    .addEventListener("submit", (evt) => {
       evt.preventDefault();
       website = webSearch.value;
       console.log(website);
    });