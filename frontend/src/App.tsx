import { main } from "@/../wailsjs/go/models";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useEffect, useState } from "react";
import { Send } from "../wailsjs/go/main/App";
import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import logo from "./assets/images/logo-universal.png";

function App() {
  const [resultText, setResultText] = useState(
    "Please enter your message below 👇"
  );
  const [name, setName] = useState("");
  const [messages, setMessages] = useState<main.Message[]>([]);

  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  useEffect(() => {
    console.log("in init effect");
    EventsOn("newMessage", (msg: main.Message) => {
      console.log("new message:", msg);
      messages.push(msg);
      console.log("after messages:", messages);
    });
    return () => EventsOff("newMessage");
  }, []);

  function sendMessage() {
    Send(name).then(updateResultText);
  }

  return (
    <div className="flex flex-col gap-4 items-center">
      <img height={100} width={100} src={logo} alt="logo" />
      <div className="flex flex-row gap-3">
        <Label className="font-bold">Message Hash: </Label>
        <Label>{resultText}</Label>
      </div>
      <div className="flex w-full max-w-sm items-center space-x-2">
        <Input onChange={updateName} autoComplete="off" autoCorrect="off" />
        <Button className="" onClick={sendMessage}>
          Send
        </Button>
      </div>

      <h1 className="text-xl font-bold mb-2">Message History</h1>
      <ScrollArea className="h-[300px] w-[550px] rounded-md border p-4 bg-gray-100">
        <ul className="text-sm">
          {messages.map((msg, index) => (
            <li key={index} className="mb-1">
              {msg.timestamp} {msg.name} says: {msg.content}
            </li>
          ))}
        </ul>
      </ScrollArea>
    </div>
  );
}

export default App;
