import { ReactQueryDevtools } from "react-query/devtools";
import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from "react-query";
import { Listener, ListenerSchema } from "./models";
import { z } from "zod";
import { Listeners } from "./Listeners";
import { Agents } from "./Agents";
import React, { useState, useEffect } from "react";
import axios from "axios";
import { ServerProvider } from './ServerContext';
import { useServer } from './ServerContext';
import Popup from "reactjs-popup";




const queryClient = new QueryClient();

//@ts-ignore
function LockPage({ setUnlocked }) {
  const [hoverActive, setHoverActive] = useState(false);
  const [inputData, setInputData] = useState("");
  const [isValid, setIsValid] = useState(false);
  const { setServer } = useServer();

  //@ts-ignore
  const validateServer = async (ip, port) => {
    try {
      const response = await axios.get(`http://${ip}:${port}/HealthCheck`);
      if (response.data === "G2 Teamserver") {
        setIsValid(true);
        setUnlocked(true);
        setServer(ip, port); 
      }
    } catch (error) {
      console.error("Invalid server");
    }
  };

  useEffect(() => {
    //@ts-ignore
    const handleKeyPress = (e) => {
      console.log(e)
      if (hoverActive) {
        if (e.key == "Shift") {

        }
        else if (e.key === 'Enter') {
          const [ip, port] = inputData.split(':');
          validateServer(ip, port);
          setInputData("");
        } else {
          setInputData((prevData) => prevData + e.key);
        }
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => {
      window.removeEventListener('keydown', handleKeyPress);
    };
  }, [hoverActive, inputData]);

  return (
    <div>
      <p>Nothing to see here, just an innocent <span onMouseEnter={() => setHoverActive(true)} onMouseLeave={() => setHoverActive(false)}>page</span>.</p>
      {isValid && <div>Valid Server</div>}
    </div>
  );
}



function Navbar() {
  const [showListeners, setShowListeners] = useState(false);
  const { ip, port } = useServer();
  const [payloadDialog, setPayloadDialog] = useState(false);
  const [selectedPlatform, setSelectedPlatform] = useState("Windows");
  const [selectedArch, setSelectedArch] = useState("amd64");

  const generateAgent = async (platform: string, arch: string) => {
  try {
    const response = await axios.post(`http://${ip}:${port}/GenerateAgent`, { platform, arch }, {
      responseType: 'arraybuffer', 
    });

    const blob = new Blob([response.data], { type: 'application/octet-stream' });
    
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `Agent_${platform}_${arch}.bin`);
    
    document.body.appendChild(link);
    link.click();
    
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Failed to generate agent');
  }
};

  return (
    <ul className="flex child:mx-2 bg-black border-solid border-white border-2 shadow-md">
      <li className="px-4 py-2 hover:bg-gray-700">
        <a href="#">G2</a>
      </li>
      <li className="px-4 py-2 hover:bg-gray-700">
        <a href="#" onClick={() => setPayloadDialog(true)}>Payloads</a>
        {payloadDialog && (
          <Popup modal open={payloadDialog} onClose={() => setPayloadDialog(false)}>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                generateAgent(selectedPlatform, selectedArch);
                setPayloadDialog(false);
              }}
              className="bg-gray-600 text-white p-4 rounded"
            >
              <div className="mb-4">
                <label htmlFor="platform" className="block mb-2">Platform:</label>
                <select
                  id="platform"
                  name="platform"
                  className="bg-gray-700 text-white rounded p-2"
                  onChange={(e) => setSelectedPlatform(e.target.value)}
                >
                  <option value="Windows">Windows</option>
                  <option value="Linux">Linux</option>
                  <option value="Darwin">Darwin</option>
                </select>
              </div>

              <div className="mb-4">
                <label htmlFor="architecture" className="block mb-2">Architecture:</label>
                <select
                  id="architecture"
                  name="architecture"
                  className="bg-gray-700 text-white rounded p-2"
                  onChange={(e) => setSelectedArch(e.target.value)}
                >
                  <option value="amd64">amd64</option>
                  <option value="arm64">arm64</option>
                </select>
              </div>

              <button type="submit" className="bg-gray-500 hover:bg-gray-800 text-white rounded py-2 px-4">
                Generate Agent
              </button>
            </form>
          </Popup>
        )}

      </li>
      <li className="relative px-4 py-2 hover:bg-gray-700">
        <a href="#" onClick={() => setShowListeners(!showListeners)}>Listeners</a>
        {showListeners && (
          <div className="absolute z-10 w-full bg-black rounded border border-gray-300">
            <Listeners />
          </div>
        )}
      </li>
    </ul>
  );
}

function Start() {
  return (
    <main className="container mx-auto">
      <div className="rounded border-2 border-black mt-3 shadow-sm box-border"></div>
    </main>
  );
}

//@ts-ignore
function App({ unlocked, setUnlocked }) {
  return (
    <div>
      {unlocked ? (
        <>
          <Navbar />
          <Agents />
        </>
      ) : (
        <LockPage setUnlocked={setUnlocked} />
      )}
    </div>
  );
}

// Wrapper component
function Wrapper() {
  const [unlocked, setUnlocked] = useState(false);

  return (
    <ServerProvider>
      <QueryClientProvider client={queryClient}>
        <App unlocked={unlocked} setUnlocked={setUnlocked} />
      </QueryClientProvider>
    </ServerProvider>
  );
}

export default Wrapper;