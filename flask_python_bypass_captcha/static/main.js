window.addEventListener("load", () => {
    const myForm = document.getElementById("myForm");
    myForm.addEventListener('submit',(event)=>{
        event.preventDefault();
        sendData();
    });

    function goToStepTwo(){
        console.log("Step Two");
        let newElement = document.createElement('div');
        newElement.textContent = "This is step Two";
        while(myForm.firstChild){
            myForm.removeChild(myForm.lastChild);
        }
        var flagElement = document.createElement('h3');
        flagElement.textContent = "Step Two";
        newElement.appendChild(flagElement);
        myForm.appendChild(newElement);
    }


    function sendData(){
        const formData = new FormData(myForm);
        if ((formData.get('email') === undefined || formData.get('email') === '') &&
         (formData.get('captcha') === undefined || formData.get('captcha') === '')){
            console.log("Degerler bos.");
            return;
        }
        fetch('/check_token',{
            method: 'post',
            body: formData
        }).then(response =>response.json())
            .then(ress=> {
                console.log("validation => ",ress.validation);
                if(ress.validation==true){
                    goToStepTwo(ress.sessValid);
                }else{
                    window.location.reload();
                }
            })
        .catch((err)=>{});
    }
});




