import React from 'react';

function AgentCard({ agent, removeAgent, sendCommand }) {
  return (
    <div className="agent-card">
      <h3>{agent.Hostname}</h3>
      <p>ID: {agent.Id}</p>
      <p>Status: {agent.ProcessName}</p>
      <p>Username: {agent.Username}</p>
      <p>IP: {agent.Ip}</p>
      <p>Process ID: {agent.ProcessId}</p>
      <p>Architecture: {agent.Architecture}</p>
      <button onClick={() => removeAgent(agent.Id)}>Remove</button>
      <button onClick={() => sendCommand(agent.Id)}>Send Command</button>
    </div>
  );
}
export default AgentCard;
