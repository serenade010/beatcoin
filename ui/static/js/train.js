const selectStep = (s) => {
  //Select Tab
  lis = document.querySelectorAll('li');
  lis.forEach((li) => {
    li.classList.remove('is-active');
  });

  // let sselectedtwo = document.getElementById(`step-two`);
  // sselectedtwo.classList.remove('is-active');
  // let sselectedthree = document.getElementById(`step-three`);
  // sselectedthree.classList.remove('is-active');

  let sselected = document.getElementById(`step-${s}`);
  sselected.classList.add('is-active');

  //Select Section
  lisSec = document.querySelectorAll('section');

  lisSec.forEach((li) => {
    li.classList.add('is-hide');
  });
  document.querySelector(`.step-${s}`).classList.remove('is-hide');

  return false;
};
