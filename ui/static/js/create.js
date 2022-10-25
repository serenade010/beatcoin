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
