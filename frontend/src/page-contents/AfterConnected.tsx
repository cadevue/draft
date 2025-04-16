import { useState } from 'react';

interface AfterConnectedProps {
    ip: string;
    port: string;
}

const pingWebClient = async (ip: string, port: string) => {
    try {
        const endpoint = `http://${ip}:8080/ping?ip=${ip}&port=${port}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
};
  
const getWebClient = async (ip: string, port: string, key:string) => {
    try {
        const endpoint = `http://${ip}:8080/get?ip=${ip}&port=${port}&key=${key}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
};
  
const setWebClient = async (ip: string, port: string, key:string, value:string) => {
    try {
        const endpoint = `http://${ip}:8080/set?ip=${ip}&port=${port}&key=${key}&value=${value}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
  }
  
const strlnWebClient = async (ip: string, port: string, key:string) => {
    try {
        const endpoint = `http://${ip}:8080/strln?ip=${ip}&port=${port}&key=${key}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
}
  
const delWebClient = async (ip: string, port: string, key:string) => {
    try {
        const endpoint = `http://${ip}:8080/del?ip=${ip}&port=${port}&key=${key}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
}
  
const appendWebClient = async (ip: string, port: string, key:string, value:string) => {
    try {
        const endpoint = `http://${ip}:8080/append?ip=${ip}&port=${port}&key=${key}&value=${value}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
}

const requestLog = async (ip: string, port: string) => {
    try {
        const endpoint = `http://${ip}:8080/request_logs?ip=${ip}&port=${port}`;
        const response = await fetch(endpoint);
        const text = await response.text();
        if (!response.ok) {
            return [false, text];
        }
        return [true, text];
    } catch (error) {
        console.warn(error);
        const err = error as Error;
        return [false, err.message];
    }
}

const AfterConnected = (props: AfterConnectedProps) => {
    console.log(props);
    const [selectedOption, setSelectedOption] = useState("Select command");
    const [response, setResponse] = useState("Not Yet");
  
    const [key, setKey] = useState("");
    const [value, setValue] = useState("");
  
    function onButtonSubmit() {
        setResponse("Loading...");
        switch (selectedOption) {
            case "ping":
                pingWebClient(props.ip, props.port).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "get":
                getWebClient(props.ip, props.port, key).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "set":
                setWebClient(props.ip, props.port, key, value).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "strln":
                strlnWebClient(props.ip, props.port, key).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "del":
                delWebClient(props.ip, props.port, key).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "append":
                appendWebClient(props.ip, props.port, key, value).then(([isConnected, message]) => {
                    if (isConnected) {
                    setResponse(message.toString());
                    } else {
                    setResponse(message.toString());
                    }
                });
                break;
            case "request_log":
                requestLog(props.ip, props.port).then(([isConnected, message]) => {
                    if (isConnected) {
                        setResponse(message.toString());
                    } else {
                        setResponse(message.toString());
                    }
                });
                break;
            default:
                setResponse("Invalid command");
                break;
        }
    }
  
    return (
        <div>
            <div className="flex flex-col gap-5">
            <div className="flex flex-row items-center gap-10">
                <label className="min-w-16">Command: </label>
                <div className="dropdown dropdown-bordered w-full">
                <div tabIndex={0} className="m-1 btn min-w-52 bg-idle">
                    {selectedOption}
                </div>
                <ul
                    tabIndex={0}
                    className="p-2 shadow menu dropdown-content bg-base-100 rounded-box w-52"
                >
                    <li>
                    <a onClick={() => setSelectedOption("ping")}>ping</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("get")}>get</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("set")}>set</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("strln")}>strln</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("del")}>del</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("append")}>append</a>
                    </li>
                    <li>
                    <a onClick={() => setSelectedOption("request_log")}>request_log</a>
                    </li>
                </ul>
                </div>
            </div>
    
            <div className="flex flex-row items-center gap-10">
                <label className="min-w-16">Key: </label>
                <input
                type="text"
                placeholder="Key"
                className="input input-bordered w-full"
                value={key}
                onChange={(e) => setKey(e.target.value)}
                />
            </div>
    
            <div className="flex flex-row items-center gap-10">
                <label className="min-w-16">Value: </label>
                <input
                type="text"
                placeholder="Value"
                className="input input-bordered w-full"
                value={value}
                onChange={(e) => setValue(e.target.value)}
                />
            </div>
    
            <button
                className="btn hover:bg-teal bg-sky text-crust"
                onClick={onButtonSubmit}
            >
                Submit
            </button>
            <p>Response: 
                <br/>
                {response.split('\n').map((line, i) => (
                    <span key={i}>
                        {line}
                        <br/>
                    </span>
                )) }
            </p>
            </div>
        </div>
    );
  };
  
  export default AfterConnected;