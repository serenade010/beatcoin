const buttonLeft = document.getElementById('btn-left');
const buttonRight = document.getElementById('btn-right');

let elems = document.getElementsByTagName('input');
let len = elems.length;

const showmodal = () => {
  console.log('object');
  modal.style.display = 'block';
};

buttonRight.addEventListener('click', showmodal);

const toggleInput = () => {
  for (var i = 0; i < len; i++) {
    console.log('first');
    elems[i].disabled = false;
  }

  buttonLeft.innerText = 'Submit';
  buttonRight.innerText = 'Cancel';
  buttonRight.removeEventListener('click', showmodal);
  buttonRight.addEventListener('click', cancel);
};

const cancel = () => {
  for (var i = 0; i < len; i++) {
    console.log('second');
    elems[i].disabled = true;
  }
  buttonLeft.innerText = 'Edit';
  buttonRight.innerText = 'Delete';
  buttonRight.removeEventListener('click', cancel);
  buttonRight.addEventListener('click', showmodal);
};

// Get the modal
var modal = document.getElementById('myModal');

// Get the button that opens the modal
var btn = document.getElementById('myBtn');

// Get the <span> element that closes the modal
var span = document.getElementsByClassName('close')[0];

// When the user clicks on the button, open the modal

// When the user clicks on <span> (x), close the modal
span.onclick = function () {
  modal.style.display = 'none';
};

// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
  if (event.target == modal) {
    modal.style.display = 'none';
  }
};
