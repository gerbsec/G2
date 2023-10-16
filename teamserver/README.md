# API Documentation

## Agent Endpoints

### `GET /Agents`
- **Description:** Retrieve a list of all agents.
- **Response:** JSON array of agents.
- **CURL Example:**

```bash
curl -X GET http://localhost:8080/Agents
```


### `GET /Agents/:agentId`
- **Description:** Retrieve a specific agent by its ID.
- **Response:** JSON representation of the agent or "Agent not found" error message.
- **CURL Example:**

```bash
curl -X GET http://localhost:8080/Agents/your-fake-guid-here
```


### `GET /Agents/:agentId/tasks`
- **Description:** Get a list of task results for a specific agent.
- **Response:** JSON array of task results or "Agent not found" error message.
- **CURL Example:**

```bash
curl -X GET http://localhost:8080/Agents/your-fake-guid-here/tasks
```


### `GET /Agents/:agentId/tasks/:taskId`
- **Description:** Get a specific task result for a specific agent.
- **Response:** JSON representation of the task result or "Task not found" error message.
- **CURL Example:**

```bash
curl -X GET http://localhost:8080/Agents/your-fake-guid-here/tasks/your-task-id-here
```


### `DELETE /Agents/:agentId/RemoveAgent`
- **Description:** Remove a specific agent by its ID.
- **CURL Example:**

```bash
curl -X DELETE http://localhost:8080/Agents/your-fake-guid-here/RemoveAgent
```


### `POST /Agents/:agentId`
- **Description:** Add a new task to a specific agent.
- **Request Body:** 
- `command`: Command to be executed.
- `arguments`: List of arguments for the command.
- `file`: Binary file content.
- **Response:** JSON representation of the created task with path to access it.
- **CURL Example:**

```bash
curl -X POST http://localhost:8080/Agents/f32701aa-fcfc-4346-b716-f794071d257b \
     -H "Content-Type: application/json" \
     -d '{
          "command": "your-command-here",
          "arguments": ["arg1", "arg2"],
          "file": "your-binary-data-here"
         }'
```

### `POST /GenerateAgent`

- **Description:** Generate an agent.
- **Request Body:**
- `os``: Operating system (e.g., "windows").
- `arch``: Architecture (e.g., "amd64").
- **Response**: Binary file payload for the agent.
- **CURL Example:**

```bash
curl -X POST http://localhost:8080/GenerateAgent \
     -H "Content-Type: application/json" \
     -d '{
          "os": "windows",
          "arch": "amd64"
         }' \
     -o payload.exe
```

## Listener Endpoints

### `GET /Listeners/:name`
- **Description:** Retrieve listener information by its name.
- **Response:** JSON representation of the listener or error message.
- **CURL Example:**
```bash
curl -X GET http://localhost:8080/Listeners/your-listener-name-here
```

### `GET /Listeners`
- **Description:** Retrieve all listener information.
- **Response:** JSON array of listeners.
- **CURL Example:**
```bash
curl -X GET http://localhost:8080/Listeners
```

### `POST /Listener`
- **Description:** Create a new listener.
- **Request Body:** 
  - `name`: Name of the listener.
  - `bindPort`: Port to bind the listener to.
- **Response:** "Listener created" or error message.
- **CURL Example:**
```bash
curl -X POST http://localhost:8080/Listener \
     -H "Content-Type: application/json" \
     -d '{
          "name": "your-listener-name-here",
          "bindPort": 8081
         }'
```

### `DELETE /StopListener/:name`
- **Description:** Stop a listener by its name.
- **Response:** "Listener stopped" message.
- **CURL Example:**
```bash
curl -X DELETE http://localhost:8080/StopListener/your-listener-name-here
```
