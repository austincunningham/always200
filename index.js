const express = require("express");
const app = express();

const PORT = 8080;
const HOST = '0.0.0.0';

app.get('/get',(req,res) => {
    console.log("GET response");
    res.sendStatus(200);
});

app.post('/post',(req,res) => {
    console.log("POST response");
    res.sendStatus(200);
});

app.listen(PORT, HOST, () => {
    console.log(`Started on http://${HOST}:${PORT}`);
  })