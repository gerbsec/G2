import { z } from "zod";

export const ListenerSchema = z.object({
  name: z.string(),
  bindPort: z.string(),
});

export type Listener = z.infer<typeof ListenerSchema>;

export const AgentMetadataSchema = z.object({
  metadata: z.object({
    id: z.string(),
    hostname: z.string(),
    username: z.string(),
    ip: z.string(),
    processName: z.string(),
    processId: z.number(),
    architecture: z.string(),
  }),
  lastSeen: z.string(),
});

export type AgentMetadata = z.infer<typeof AgentMetadataSchema>;
