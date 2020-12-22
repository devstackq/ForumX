function showAnswerForm(idAnswer, btn) {
  console.log(idAnswer);
  let f = document.getElementById("answerFormId");
  document
    .getElementById("btn-reply-answer")
    .addEventListener("click", function () {
      f.classList.toggle("visible-answer-form");
    });
}

const test = (data) => {
  console.log(data, "D");
};
//window.onload = fetch request - > backend, backen response json { authentficate value }
// const showButton = () => {

//     let btn = document.getElementById("create_post");
//     if
//     btn.sty
// }
