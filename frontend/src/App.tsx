import { useEffect, useState } from "react";
import logo from "./assets/images/logo-universal.png";
import { Send } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

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
    <div className="flex flex-col gap-2 items-center">
      <img height={100} width={100} src={logo} alt="logo" />
      <div className="">
        Message Hash: {resultText}
      </div>
      <div className="flex w-full max-w-sm items-center space-x-2">
        <Input onChange={updateName} autoComplete="off" autoCorrect="off" />
        <Button className="" onClick={sendMessage}>
          Send
        </Button>
      </div>
    </div>
  );
}

export default App;
