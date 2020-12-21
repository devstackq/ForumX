function showAnswerForm(idAnswer, btn) {
    console.log(idAnswer);
    let f = document.getElementById("answerFormId");
    document
        .getElementById("btn-reply-answer")
        .addEventListener("click", function() {
            f.classList.toggle("visible-answer-form");
        });
}

const test = (data) => {
    console.log(data, "D")
};
//window.onload = fetch request - > backend, backen response json { authentficate value }
// const showButton = () => {

//     let btn = document.getElementById("create_post");
//     if
//     btn.sty
// }

window.onload = async function() {
    let response = await fetch('http://localhost:6969/', {
        mode: 'cors',
        method: 'get',
    })
    let d = await response.json()
    console.log("windwos onload1")

    .then(d => test(d));
    console.log("windwos onload2")
        //let data = await response.json()
}