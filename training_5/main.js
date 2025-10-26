let btn1= document.getElementById("btn1")
const obj = { name: "Mahdi", lastname: "boukhouima", age : 27 };

const send_data = async () => {
  try {
    const res = await fetch("http://localhost:8080/AddUser", {
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
    fetch("http://localhost:8080/SendUser", {
    method : "GET",
    headers : {"Content-Type" : "application/json"}
  }).then(response => {
    if (!response.ok){
      throw new Error("Network response was not ok")
    }
    return response.json()
  }).then(data => {
    console.log(data)
  }).catch(error => {
    console.log(error)
  })
  } catch (error) {
    console.log(error)
  }

}

btn1.addEventListener("click" , () => click_me())
