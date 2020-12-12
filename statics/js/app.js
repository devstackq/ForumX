function showAnswerForm(idAnswer, btn) {
  console.log(idAnswer);
  let f = document.getElementById("answerFormId");
  document
    .getElementById("btn-reply-answer")
    .addEventListener("click", function () {
      f.classList.toggle("visible-answer-form");
    });
}

const showAnswerComment = () => {};
