import { ReactQueryDevtools } from "react-query/devtools";
import { useState } from "react";
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

const queryClient = new QueryClient();

function Navbar() {
  return (
    <ul className="flex child:mx-2 bg-black border-solid border-white border-2">
      <li>
        <a href="/">G2</a>
      </li>
      <li>
        <a href="/payloads">Payloads</a>
      </li>
      <li>
        <a href="/listeners">Listeners</a>
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

function App() {
  return (
    <>
      <Navbar />
      <Listeners />

      <ReactQueryDevtools initialIsOpen={false} />
    </>
  );
}

function Wrapper() {
  return (
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  );
}

export default Wrapper;
