import { useState } from "react";
import AfterConnected from "./page-contents/AfterConnected"
import BeforeConnected from "./page-contents/BeforeConnected"

import Maya from "./assets/maya.png"

function App() {
  const [isConnected, setIsConnected] = useState(false)
  const [ip, setIp] = useState("");
  const [port, setPort] = useState("");

  const OnConnected = (ip: string, port: string) => {
    setIp(ip);
    setPort(port);
    setIsConnected(true);
  }

  return (
    <>
    <div className="flex justify-center items-end gap-4 mt-12">
      <img src={Maya} className="w-16"/>
      <h1 className='text-6xl text-red-600 text-center font-bold'>DRaft</h1>
    </div>
    <div className="flex flex-col w-2/5 h-full justify-center gap-8">
      {isConnected ? <AfterConnected ip={ip} port={port}/> : <BeforeConnected onConnected={OnConnected}/>}
    </div>
    </>
  )
}

export default App
