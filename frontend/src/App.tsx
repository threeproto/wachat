import { useEffect, useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { Send } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

function App() {
  const [resultText, setResultText] = useState(
    "Please enter your name below 👇"
  );
  const [name, setName] = useState("");
  const [messages, setMessages] = useState<string[]>([]);

  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  useEffect(() => {
    console.log("in init effect");
    EventsOn("newMessage", (message: string) => {
      console.log("new message:", message);
      messages.push(message);
      console.log("received:", messages);
    });
  }, []);

  function sendMessage() {
    Send(name).then(updateResultText);
  }

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo" />
      <div id="result" className="result">
        {resultText}
      </div>
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={updateName}
          autoComplete="off"
          name="input"
          type="text"
        />
        <button className="btn" onClick={sendMessage}>
          Greet
        </button>
      </div>
    </div>
  );
}

export default App;
