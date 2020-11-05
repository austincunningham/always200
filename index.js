const express = require("express");
const app = express();

app.get('/',(req,res) => {
    res.sendStatus(200);
});

app.post('/',(req,res) => {
    console.log("do i get here")
    res.sendStatus(200);
});

app.listen(8080,() => {
    console.log("Started on PORT 8080");
  })