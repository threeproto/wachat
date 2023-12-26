import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import { MailPage } from '@/components/mail';
import { accounts, mails } from "@/data"

function App() {
    // const [resultText, setResultText] = useState("Please enter your name below 👇");
    // const [name, setName] = useState('');
    // const updateName = (e: any) => setName(e.target.value);
    // const updateResultText = (result: string) => setResultText(result);

    // function greet() {
    //     Greet(name).then(updateResultText);
    // }

    // return (
    //     <div id="App">
    //         <img src={logo} id="logo" alt="logo"/>
    //         <div id="result" className="result">{resultText}</div>
    //         <div id="input" className="input-box">
    //             <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
    //             <button className="btn" onClick={greet}>Greet</button>
    //         </div>
    //     </div>
    // )

    // const layout = cookies().get("react-resizable-panels:layout")
//   const collapsed = cookies().get("react-resizable-panels:collapsed")

//   const defaultLayout = layout ? JSON.parse(layout.value) : undefined
//   const defaultCollapsed = collapsed ? JSON.parse(collapsed.value) : undefined

  return (
    <>
      <div className="md:hidden">
        {/* <Image
          src="/examples/mail-dark.png"
          width={1280}
          height={727}
          alt="Mail"
          className="hidden dark:block"
        />
        <Image
          src="/examples/mail-light.png"
          width={1280}
          height={727}
          alt="Mail"
          className="block dark:hidden"
        /> */}
      </div>
      <div className="hidden flex-col md:flex">
        <MailPage
          accounts={accounts}
          mails={mails}
          defaultLayout={undefined}
          defaultCollapsed={undefined}
          navCollapsedSize={4}
        />
      </div>
    </>
  )
}

export default App
