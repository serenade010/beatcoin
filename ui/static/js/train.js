const selectStep = (s) => {
  //Select Tab
  lis = document.querySelectorAll('li');
  lis.forEach((li) => {
    li.classList.remove('is-active');
  });
  lisSec = document.querySelectorAll('section');
  // Select Section

  lisSec.forEach((li) => {
    li.classList.add('is-hide');
  });
  document.querySelector(`.step-${s}`).classList.remove('is-hide');
};
