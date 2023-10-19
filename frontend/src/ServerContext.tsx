import React, { createContext, useContext, useState } from 'react';

interface ServerContextProps {
    ip: string;
    port: string;
    setServer: (ip: string, port: string) => void;
}

const ServerContext = createContext<ServerContextProps | null>(null);

type ServerProviderProps = {
    children: React.ReactNode;
};


export const ServerProvider: React.FC<ServerProviderProps> = ({ children }) => {
    const [ip, setIp] = useState('');
    const [port, setPort] = useState('');

    const setServer = (newIp: string, newPort: string) => {
        setIp(newIp);
        setPort(newPort);
    };

    return (
        <ServerContext.Provider value={{ ip, port, setServer }}>
            {children}
        </ServerContext.Provider>
    );
};

export const useServer = () => {
    const context = useContext(ServerContext);
    if (!context) {
        throw new Error('useServer must be used within a ServerProvider');
    }
    return context;
};
