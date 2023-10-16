import React from 'react';

function ListenerCard({ listener, removeListener }) {
  return (
    <div className="listener-card">
      <h3>Name: {listener.name}</h3>
      <p>Port: {listener.bind_port}</p>
      <button onClick={() => removeListener(listener.name)}>Remove</button>
    </div>
  );
}

export default ListenerCard;
