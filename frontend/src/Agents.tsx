import React, { useState, useEffect, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { useServer } from './ServerContext';
import { AgentMetadata } from './models';

export function Agents() {
    const { ip, port } = useServer();
    const queryClient = useQueryClient();
    const [selectedAgent, setSelectedAgent] = useState<AgentMetadata | null>(null);
    const [tasks, setTasks] = useState<string[]>([]);
    const [taskIds, setTaskIds] = useState<string[]>([]);
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const intervalId = setInterval(() => {
            queryClient.invalidateQueries('agents');
        }, 1000);
        return () => clearInterval(intervalId);
    }, [queryClient]);

    useEffect(() => {
        if (selectedAgent && taskIds.length > 0) {
            const intervalId = setInterval(() => {
                Promise.all(
                    taskIds.map((taskId) =>
                        fetch(`http://${ip}:${port}/Agents/${selectedAgent.metadata.id}/tasks/${taskId}`)
                            .then((res) => res.json())
                    )
                ).then((outputs) => {
                    const newTasks = outputs.map((o, i) => `$ ${o.Id}\n${o.Result}`);
                    setTasks(newTasks);
                });
            }, 1000);
            return () => clearInterval(intervalId);
        }
    }, [selectedAgent, ip, port, taskIds]);

    const { data: agents, isLoading, error } = useQuery('agents', () => {
        return fetch(`http://${ip}:${port}/Agents`)
            .then((res) => res.json());
    });

    const removeAgent = useMutation(
        (agentId: string) => {
            return fetch(`http://${ip}:${port}/Agents/${agentId}/RemoveAgent`, {
                method: 'DELETE',
            });
        },
        {
            onSuccess() {
                queryClient.invalidateQueries('agents');
            },
        }
    );

    const handleCommand = async () => {
        const input = inputRef.current?.value || '';
        if (selectedAgent) {
            const res = await fetch(`http://${ip}:${port}/Agents/${selectedAgent.metadata.id}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ command: input, arguments: [] }),
            });
            const data = await res.json();
            const newTaskId = data.task.Id;

            setTaskIds((prevIds) => [...prevIds, newTaskId]);
            if (inputRef.current) inputRef.current.value = '';
        }
    };

    const handleRowClick = (agent: AgentMetadata) => {
        setSelectedAgent(agent);
    };

    return (
        <>
            {selectedAgent && (
                <div className="slide-up-terminal">
                    {tasks.map((task, index) => (
                        <div key={index}>{task}</div>
                    ))}
                    <input
                        type="text"
                        ref={inputRef}
                        className="command-input"
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                handleCommand();
                            }
                        }}
                    />
                    <button onClick={() => setSelectedAgent(null)}>Close Terminal</button>
                </div>
            )}
            <main className="container mx-auto">
                <div className="w-full bg-gray-700 rounded border border-gray-300 p-4 min-h-[100px] min-w-[200px] overflow-auto">
                    {(!agents || agents.length === 0) ? (
                        <div className="text-center text-white-500">No agents available</div>
                    ) : (
                        <table className="w-full text-sm text-left text-white-500 dark:text-gray-400">
                            <thead className="text-md text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                                <tr>
                                    <th className="text-md text-white bg-gray-800">ID</th>
                                    <th className="text-md text-white bg-gray-800">Hostname</th>
                                    <th className="text-md text-white bg-gray-800">Username</th>
                                    <th className="text-md text-white bg-gray-800">IP</th>
                                    <th className="text-md text-white bg-gray-800">ProcessName</th>
                                    <th className="text-md text-white bg-gray-800">ProcessId</th>
                                    <th className="text-md text-white bg-gray-800">Architecture</th>
                                    <th className="text-md text-white bg-gray-800"></th>
                                </tr>
                            </thead>
                            <tbody>
                                {agents?.map((agent) => (
                                    <tr key={agent.metadata.id} onClick={() => handleRowClick(agent)}>
                                        <td>{agent.metadata.id}</td>
                                        <td>{agent.metadata.hostname}</td>
                                        <td>{agent.metadata.username}</td>
                                        <td>{agent.metadata.ip}</td>
                                        <td>{agent.metadata.processName}</td>
                                        <td>{agent.metadata.processId}</td>
                                        <td>{agent.metadata.architecture}</td>
                                        <td>{agent.lastSeen}</td>
                                        <td>
                                            <button
                                                className="text-red-600 w-8 bg-white box-border border-2 border-white rounded"
                                                onClick={() => removeAgent.mutate(agent.metadata.id)}
                                            >
                                                X
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>
            </main>
        </>
    );
}
