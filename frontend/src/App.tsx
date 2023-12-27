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
  const [newMessageHash, setNewMessageHash] = useState(
    "Please enter your message below 👇"
  );
  const [newMessage, setNewMessage] = useState("");
  const [messages, setMessages] = useState<main.Message[]>([]);

  const updateMessage = (e: any) => setNewMessage(e.target.value);

  useEffect(() => {
    console.log("in init effect");
    EventsOn("newMessage", (msg: main.Message) => {
      setMessages((prev) => [...prev, msg]);
    });
    return () => EventsOff("newMessage");
  }, []);

  const sendMessage = async () => {
    let result = await Send(newMessage);
    setNewMessageHash(result);
    setNewMessage("");
  };

  const formatDate = (timestamp: number) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  return (
    <div className="flex flex-col gap-4 items-center">
      <img height={100} width={100} src={logo} alt="logo" />
      <div className="flex flex-row gap-3">
        <Label className="font-bold">Message Hash: </Label>
        <Label>{newMessageHash}</Label>
      </div>
      <div className="flex w-full max-w-sm items-center space-x-2">
        <Input
          value={newMessage}
          onChange={updateMessage}
          autoComplete="off"
          autoCorrect="off"
        />
        <Button className="" onClick={sendMessage}>
          Send
        </Button>
      </div>

      <h1 className="text-xl font-bold mb-2">Message History</h1>
      <ScrollArea className="h-[300px] w-[550px] rounded-md border p-4 bg-gray-100">
        <ul className="text-sm">
          {messages.map((msg, index) => (
            <li key={index} className="mb-1">
              [{formatDate(msg.timestamp)} {msg.name}] says: {msg.content}
            </li>
          ))}
        </ul>
      </ScrollArea>
    </div>
  );
}

export default App;
