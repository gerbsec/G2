import React from 'react';

function AgentCard({ agentId, hostname, username, ip, processName, processId, architecture, lastSeen, removeAgent, sendCommand }) {
  return (
    <div className="agent-card">
      <p>ID: {agentId}</p>
      <p>Hostname: {hostname}</p>
      <p>Username: {username}</p>
      <p>IP Address: {ip}</p>
      <p>Process Name: {processName}</p>
      <p>Process ID: {processId}</p>
      <p>Architecture: {architecture}</p>
      <p>Last Seen: {lastSeen}</p>
      <button onClick={removeAgent}>Remove</button>
      <button onClick={() => sendCommand(agentId)}>Send Command</button>
    </div>
  );
}

export default AgentCard;
