function Prompt() {
    let toast = function(c) {
    // This is a kind of Destructuring assignment syntex and use Default values syntex.
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;
        
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer)
            toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    let success = function(c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer
        })
    }

    let error = function(c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer
        })
    }

    async function custom(c) {
    
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;
        

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            background: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
            //   const elem = document.getElementById('reservation-dates-modal');
            //   const rangePicker = new DateRangePicker(elem, {
            //     format: "yyyy-mm-dd",
            //     showOnFocus: true,
            //   });
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            didOpen: () => {
                // document.getElementById('start').removeAttribute("disabled");
                // document.getElementById('end').removeAttribute("disabled");
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
        })
        
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel){
                if (result.value !== "") {
                    if (c.callback != undefined) {
                        c.callback(result)
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        }

    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}


function SearchAvailability(room_id) {
    let html = `
    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
        <div class="form-row">
            <div class="col">
                <div class="form-row" id="reservation-dates-modal">
                    <div class="col">
                        <input required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                    </div>
                    <div class="col">
                        <input required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                    </div>

                </div>
            </div>
        </div>
    </form>
    `;
    attention.custom({
        title: 'Choose your dates',
        msg: html,
        didOpen: () => {
            const elem = document.getElementById('reservation-dates-modal');
            const rangePicker = new DateRangePicker(elem, {
                format: "yyyy-mm-dd",
                showOnFocus: true,
                minDate: new Date(),
            });   
        },
        callback: function(result) {
            
            let form = document.getElementById("check-availability-form");

            let formData = new FormData(form);

            formData.append("csrf_token", "{{.CSRFToken}}");
            formData.append("room_id", room_id)
            // js 獨有的 fetch api 與ajax不同
            // 以下範例為得到該請求回傳的json資料
            fetch('/search-availability-json', {
                method : "post",
                body : formData,
            })
            .then(response => {
                return response.json();
            })
            .then(data => {
                if (data.ok) {
                    attention.custom({
                        icon: 'success',
                        msg: '<p>Room is available !</p>'
                        + '<p><a href="/book-room?id=' 
                        + data.room_id
                        + '&s='
                        + data.start_date
                        + '&e='
                        + data.end_date
                        + '" class="btn btn-primary">'
                        + 'Book now !</a></p>',
                        showConfirmButton: false,
                    })
                }else {
                    attention.error({
                        msg: "No availability",
                    })
                }
            })
        },
    });
}