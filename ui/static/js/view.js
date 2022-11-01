const buttonLeft = document.getElementById('btn-left');
const buttonRight = document.getElementById('btn-right');

let elems = document.getElementsByTagName('input');
let len = elems.length;

//Show modal and add to btn eventlistener
const showmodal = () => {
  console.log('object');
  modal.style.display = 'block';
};

buttonRight.addEventListener('click', showmodal);

const toggleInput = () => {
  for (var i = 0; i < elems.length; i++) {
    console.log('first');
    elems[i].disabled = false;
  }
  buttonLeft.innerText = 'Save';
  buttonRight.innerText = 'Cancel';
  buttonLeft.addEventListener('click', submitForm);
  buttonRight.removeEventListener('click', showmodal);
  buttonRight.addEventListener('click', cancel);
};

const cancel = () => {
  for (var i = 0; i < elems.length; i++) {
    console.log('second');
    elems[i].disabled = true;
  }
  buttonLeft.innerText = 'Edit';
  buttonRight.innerText = 'Delete';
  buttonLeft.removeEventListener('click', submitForm);
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

function checkboxControl(j) {
  var total = 0;
  var elem = document.getElementsByClassName('checkbox');
  var error = document.getElementById('error');
  for (var i = 0; i < elem.length; i++) {
    if (elem[i].checked == true) {
      total = total + 1;
    }
    if (total > 3) {
      error.textContent = 'You Must Select at Most 3 index';
      elem[j].checked = false;
      return false;
    }
  }
}

const submitForm = () => {
  document.getElementById('put-form').submit();
};
