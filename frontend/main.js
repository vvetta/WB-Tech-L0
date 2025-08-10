const api_url = "http://localhost:8081/order/";

const input_form = document.getElementById("input-form");
const message_container = document.getElementById("message-container");
const textarea = document.getElementById("textarea-json");
const input_id_container = document.getElementById("input-id-container");
const textarea_container = document.getElementById("view-order-container");
const reload_button = document.getElementById("reload-button-container");

function showResponse(json_text) {
  input_id_container.style.display = "none";
  textarea_container.style.display = "flex"; 
  reload_button.style.display = "flex";
  textarea.textContent = json_text;
}

function hideTextAreaContainer() {
  textarea_container.style.display = "none";
  input_id_container.style.display = "flex";
  reload_button.style.display = "none"; 
  input_form.reset();
}

function showMessage(text, message_type) {
  // Может быть два тип сообщений: green и red.

  switch (message_type) {
    case "green":
      message_container.style.backgroundColor = message_type;
      break; 
    case "red":
      message_container.style.backgroundColor = message_type;
      break; 
    default:
      message_container.style.backgroundColor = "red";
  }

  message_container.style.top = "0px";
  message_container.textContent = text;

  setTimeout(() => {
    hideMessage();
  }, 5000);
}

function hideMessage() {
  if (message_container.style.top === "0px") {
    message_container.style.top = "-100px";
  }
}

function serializeForm(formNode) {
  const data = new FormData(formNode); 
  return data;
}

async function sendData(order_id) {
  return await fetch(api_url+order_id, {
    method: "POST"
  })
};

async function handleFormSubmit(event) {
  event.preventDefault();

  const data = serializeForm(input_form);
  const order_id = data.get("order_id");
 
  if (order_id === "") {
    showMessage("Поле не должно быть пустым!");
    return; 
  }

  try {
    const response = await sendData(order_id);
    const json = await response.json();

    if (response.status != 200) {
      showMessage(json.message);
    }

    showResponse(JSON.stringify(json));
    showMessage("Запрос выполнен!", "green");

  } catch (error) {
    // Произошла ошибка при отправке формы. 
    showMessage(error);
    input_form.reset(); 
  }
};

input_form.addEventListener("submit", handleFormSubmit);
reload_button.addEventListener("click", hideTextAreaContainer);
