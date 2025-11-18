async function fetchStreamingData() {
  const response = await fetch("http://localhost:8080/GetData", {
    method: "GET",
    headers: {
      "Content-Type": "application/json"
    }
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder("utf-8");
  let fullData = "";

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    fullData += decoder.decode(value);
  }

  // Split by lines and parse JSON
  const users = fullData
    .trim()
    .split("\n")
    .map(line => JSON.parse(line));

  console.log("All users:", users);
}

fetchStreamingData();
