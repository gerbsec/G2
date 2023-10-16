import React, { useState, useEffect } from 'react';
import axios from 'axios';
import AgentCard from './AgentCard';
import ListenerCard from './ListenerCard';
import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  const [serverAddress, setServerAddress] = useState('');
  const [agents, setAgents] = useState([]);
  const [command, setCommand] = useState('');
  const [commandArgs, setCommandArgs] = useState('');
  const [listenerName, setListenerName] = useState('');
  const [listenerPort, setListenerPort] = useState('');
  const [listeners, setListeners] = useState([]);
  const [isVerified, setIsVerified] = useState(false);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);

  const toggleModal = () => setShowModal(!showModal);

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        const healthData = await axios.get(`http://${serverAddress}/HealthCheck`);
        if (healthData.data === "G2 Teamserver" && isMounted) {
          setIsVerified(true);
        }

        if (isVerified) {
          const agentData = await axios.get(`http://${serverAddress}/Agents`);
          if (isMounted) {
            setAgents(agentData.data);
          }

          const listenerData = await axios.get(`http://${serverAddress}/Listeners`);
          if (isMounted) {
            setListeners(listenerData.data);
          }
        }

      } catch (err) {
        if (isMounted) {
          setError(err.toString());
        }
      }
    };

    if (serverAddress) {
      fetchData();
      const intervalId = setInterval(fetchData, 5000);

      return () => {
        isMounted = false;
        clearInterval(intervalId);
      };
    }
  }, [serverAddress, isVerified]);

  const generateAgent = () => {
    axios.post(`http://${serverAddress}/GenerateAgent`, {
      os: "windows",
      arch: "amd64"
    }).catch(err => setError(err.toString()));
  };

  const addListener = () => {
    axios.post(`http://${serverAddress}/Listener`, {
      name: listenerName,
      bind_port: listenerPort
    }).catch(err => setError(err.toString()));
  };

  const removeListener = (name) => {
    axios.delete(`http://${serverAddress}/StopListener/${name}`)
      .then(response => {
      })
      .catch(err => setError(err.toString()));
  };

  const removeAgent = (id) => {
    axios.delete(`http://${serverAddress}/Agents/${id}/RemoveAgent`).catch(err => setError(err.toString()));
  };

  const sendCommand = (agentId) => {
    axios.post(`http://${serverAddress}/Agents/${agentId}`, {
      command,
      arguments: commandArgs.split(' ')
    }).catch(err => setError(err.toString()));
  };

  return (
    <div className="App">
      {isVerified ? (
        <>
          <header className="App-header">
            <button onClick={generateAgent}>Generate Agent</button>
            <Button variant="primary" onClick={toggleModal}>
              Listener Management
            </Button>
          </header>

          <Modal show={showModal} onHide={toggleModal}>
            <Modal.Header closeButton>
              <Modal.Title>Listener Management</Modal.Title>
            </Modal.Header>
            <Modal.Body>
              <div className="listeners">
                {Array.isArray(listeners) && listeners.length > 0 ? (
                  listeners.map((listener) => (
                    <ListenerCard
                      key={listener.name}
                      listener={listener}
                      removeListener={() => removeListener(listener.name)}
                    />
                  ))
                ) : (
                  <p>No listeners available.</p>
                )}
              </div>
              <input
                type="text"
                placeholder="Listener Name"
                value={listenerName}
                onChange={(e) => setListenerName(e.target.value)}
              />
              <input
                type="text"
                placeholder="Listener Port"
                value={listenerPort}
                onChange={(e) => setListenerPort(e.target.value)}
              />
            </Modal.Body>
            <Modal.Footer>
              <Button variant="secondary" onClick={toggleModal}>
                Close
              </Button>
              <Button variant="primary" onClick={addListener}>
                Add Listener
              </Button>
            </Modal.Footer>
          </Modal>

          <div className="agents">
            {Array.isArray(agents) && agents.length > 0 ? (
              agents.map((agent) => (
                <AgentCard
                  key={agent.Id}
                  agent={agent}
                  removeAgent={() => removeAgent(agent.Id)}
                  sendCommand={(id) => sendCommand(id)}
                />
              ))
            ) : (
              <p>No agents available.</p>
            )}
          </div>
        </>
      ) : (
        <div>
          Please enter a valid Teamserver IP:Port to access the application.
          <input
            type="text"
            placeholder="Enter Teamserver IP:Port"
            onChange={e => setServerAddress(e.target.value)}
          />
        </div>
      )}
    </div>
  );
}

export default App;