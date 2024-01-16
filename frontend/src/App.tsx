import { params } from "@/../wailsjs/go/models";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useEffect, useState } from "react";
import { CreateUser, GetUser, Send } from "../wailsjs/go/main/App";
import { EventsOff, EventsOn } from "../wailsjs/runtime/runtime";
import logo from "./assets/images/logo-universal.png";
import { toast } from "sonner";

function App() {
  const [newMessageHash, setNewMessageHash] = useState(
    "Please enter your message below 👇"
  );
  const [newMessage, setNewMessage] = useState("");
  const [username, setUsername] = useState("");
  const [usernameInput, setUsernameInput] = useState("");
  const [messages, setMessages] = useState<params.Message[]>([]);

  const updateMessage = (e: any) => setNewMessage(e.target.value);

  useEffect(() => {
    console.log("in init effect");
    EventsOn("newMessage", (msg: params.Message) => {
      setMessages((prev) => [msg, ...prev]);
    });
    EventsOn("isOnline", (isOnline: boolean) => {
      if (isOnline) {
        toast.success("You are online.");
      } else {
        toast.warning("You are offline.");
      }
    });

    const init = async () => {
      const name = await GetUser();
      setUsername(name);
    };
    init();

    return () => {
      EventsOff("newMessage");
      EventsOff("isOnline");
    };
  }, []);

  const sendMessage = async () => {
    if (!username || !newMessage) {
      toast.warning("Username or message is empty.");
      return;
    }
    try {
      let result = await Send(newMessage);
      setNewMessageHash(result);
      setNewMessage("");
    } catch (err) {
      toast.error(`Error happens: ${err}`);
    }
  };

  const createUser = async () => {
    try {
      const name = await CreateUser(usernameInput);
      setUsername(name);
      toast("User has been created.");
    } catch (err) {
      toast.error(`Error happens: ${err}`);
    }
  };

  const formatDate = (timestamp: number) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  return (
    <div className="flex flex-col gap-10 items-center justify-center h-screen">
      <img height={100} width={100} src={logo} alt="logo" />

      {!username && (
        <div className="flex w-full max-w-sm items-center space-x-2">
          <Input
            value={usernameInput}
            onChange={(e) => setUsernameInput(e.target.value)}
            placeholder="Enter your username"
            autoComplete="off"
            autoCorrect="off"
          />
          <Button className="w-32" onClick={createUser}>
            Create
          </Button>
        </div>
      )}

      {username && (
        <div className="flex flex-col gap-10">
          <div className="flex w-full max-w-sm items-center space-x-2">
            <Input
              value={newMessage}
              onChange={updateMessage}
              placeholder="Input your message"
              autoComplete="off"
              autoCorrect="off"
            />
            <Button className="w-32" onClick={sendMessage}>
              Send
            </Button>
          </div>

          <div>
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
        </div>
      )}
    </div>
  );
}

export default App;
