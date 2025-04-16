import { useState } from "react";

const connectToWebClient = async (ip: string, port: string) => {
  try {
    const endpoint = `http://${ip}:8080/connect?ip=${ip}&port=${port}`;
    const response = await fetch(endpoint);
    if (!response.ok) {
      const text = await response.text()
      return [false, text];
    }
    return [true, "Success"];
  } catch (error) {
    console.warn(error)
    const err = error as Error;
    return [false, err.message]
  }
};

interface BeforeConnectedProps {
  onConnected: (ip: string, port: string) => void;
}

const BeforeConnected = (props: BeforeConnectedProps) => {
  const [ip, setIp] = useState("");
  const [port, setPort] = useState("");
  const [status, setStatus] = useState("Not Connected");
  const [buttonDisabled, setButtonDisabled] = useState(false);

  const onButtonConnect = async () => {
    setStatus("Connecting . . .");
    setButtonDisabled(true)
    const [isConnected, message] = await connectToWebClient(ip, port);
    if (isConnected) {
      setStatus("Connected:D");
      props.onConnected(ip, port);
    } else {
      setStatus(message.toString());
    }
    setButtonDisabled(false)
  };

  return (
    <>
      <div className="flex flex-col gap-4">
        <p>Status: {status}</p>
        <div className="flex flex-col gap-4 text-white">
          <div className="flex flex-col gap-2">
            <label className="font-bold">Server IP</label>
            <input
              type="text"
              placeholder="Server IP"
              className="input input-bordered w-full"
              value={ip}
              onChange={(e) => setIp(e.target.value)}
            />
          </div>
          <div className="flex flex-col gap-2">
            <label className="font-bold">Server Port</label>
            <input
              type="text"
              placeholder="Server Port"
              className="input input-bordered w-full"
              value={port}
              onChange={(e) => setPort(e.target.value)}
            />
          </div>
        </div>
        <button
          onClick={onButtonConnect}
          className="btn bg-pink text-crust hover:bg-[#AC709C] mt-4
            disabled:opacity-50 disabled:bg-pink disabled:text-crust
          "
          disabled={buttonDisabled}
        >
          Connect
        </button>
      </div>
    </>
  );
};

export default BeforeConnected;
