import { z } from "zod";

export const ListenerSchema = z.object({
  name: z.string(),
  bindPort: z.string(),
});

export type Listener = z.infer<typeof ListenerSchema>;
