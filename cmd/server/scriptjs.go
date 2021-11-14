package main

var scriptjs = `
const contentData = document.getElementById("content-list");

let orderUpward = true;

let hackerList = JSON.parse(data);

/*****************init js*****************/
function init() {
    const dataInit = mapInfo(hackerList)
    contentData.innerHTML = dataInit.join(" ");
}

init();

/*********************Update Array***************/
function updateArray(data){
  hackerList.push(data);
  init();
}

/******************Filter Array*****************/
const selectElement = document.querySelector('#typeData');

selectElement.addEventListener('change', (event) => {
    const filterData = event.target.value === 'all'
      ? hackerList
      : hackerList.filter(alert => alert.type === event.target.value);

    const mapData = mapInfo(filterData)

  contentData.innerHTML = mapData.join(" ");
});

/*******************Order Array******************/
function functionOrderAndClass() {
  const contentButton = document.getElementById("contentOrder");

  const select = document.getElementById('typeData');
  const value = select.options[select.selectedIndex].value;
`
