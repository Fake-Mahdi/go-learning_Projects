let btn1= document.getElementById("btn1")
let btn2= document.getElementById("btn2")
let display= document.getElementById("display")
const obj = { name: "Natsu ", lastname: "boukhouima", age : 300 };

const send_data = async () => {
  try {
    const res = await fetch("http://localhost:8080/insertData", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(obj)
    });

    const data = await res.json();
    console.log(data);
  } catch (err) {
    console.error(err);
  }
};
const click_me = () =>{
  try {
    fetch("http://localhost:8080/selectData", {
    method : "GET",
    headers : {"Content-Type" : "application/json"}
  }).then(response => {
    if (!response.ok){
      throw new Error("Network response was not ok")
    }
    return response.json()
  }).then(data => {
    console.log(data)
    let Content = ""
    for(const element of data) {
      Content += `<ul>
            <li>${element.name}</li>
        </ul>`
    }
    display.innerHTML = Content
  }).catch(error => {
    console.log(error)
  })
  } catch (error) {
    console.log(error)
  }

}

btn1.addEventListener("click" , () => click_me())
btn2.addEventListener("click" , () => {send_data()})
