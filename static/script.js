$('#form').find('input, textarea').on('keyup blur focus', function (e) {
  
  var $this = $(this),
      label = $this.prev('label');

	  if (e.type === 'keyup') {
			if ($this.val() === '') {
          label.removeClass('active highlight');
        } else {
          label.addClass('active highlight');
        }
    } else if (e.type === 'blur') {
    	if( $this.val() === '' ) {
    		label.removeClass('active highlight'); 
			} else {
		    label.removeClass('highlight');   
			}   
    } else if (e.type === 'focus') {
      
      if( $this.val() === '' ) {
    		label.removeClass('highlight'); 
			} 
      else if( $this.val() !== '' ) {
		    label.addClass('highlight');
			}
    }

});

$('.tab a').on('click', function (e) {
  
  e.preventDefault();
  
  $(this).parent().addClass('active');
  $(this).parent().siblings().removeClass('active');
  
  target = $(this).attr('href');

  $('.tab-content > div').not(target).hide();
  
  $(target).fadeIn(800);
  
});

$('#submit_btn').on('click', async function (e) {
    e.stopImmediatePropagation();
    e.preventDefault()

    let username = document.getElementById("first_name").value
    let email = document.getElementById("email").value
    let password = document.getElementById("password").value

    let url = "http://"+window.location.host + "/" + "auth/signup"
    let body = {
        first_name: username,
        email: email,
        password: password
    }

    let res = await postData(url, body, null)
    if (res !== null){
        alert("registration successful, please proceed to login")
    }

})

$('#login_btn').on('click', async function (e) {
    e.stopImmediatePropagation();
    e.preventDefault()

    let email = document.getElementById("login-email").value
    let password = document.getElementById("login-password").value

    let url = "http://"+window.location.host + "/" + "auth/login"
    let body = {
        email: email,
        password: password
    }

    let res = await postData(url, body, null)
    if (res !== null){//successfully
        //load chat page
        sessionKey = res.sessionKey.slice(res.sessionKey.search("_")+1)
        chatroomURL =  "http://"+window.location.host + "/" + "chat.html?sessionid="+sessionKey
        location.assign(chatroomURL)
    }

})

async function postData(url, jsonBody, headers) {
    try {

        console.log(url)
        const response = await fetch(url, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: headers ? headers : {},
            redirect: 'follow',
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(jsonBody)
        });

        if (!response.ok) {
            alert("request failed\n" + response.status + "\n" + JSON.stringify(await response.json()))
            return null
        }
        return response.json();
    }
    catch(e){
        console.log(e)
        return null
    }
}