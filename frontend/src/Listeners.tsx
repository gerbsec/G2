import { useForm, SubmitHandler } from "react-hook-form";
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
import Popup from "reactjs-popup";
import "reactjs-popup/dist/index.css";
import { useServer } from './ServerContext';



export function Listeners() {
  const { ip, port } = useServer();
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [showTable, setShowTable] = useState(false);

  const {
    data: listeners,
    isLoading,
    error,
  } = useQuery("listeners", () => {
    return fetch(`http://${ip}:${port}/Listeners`)
      .then((res) => res.json())
      .then((json) => z.array(ListenerSchema).parse(json));
  });

  const deleteListener = useMutation(
    (listener: Listener) => {
      return fetch(`http://${ip}:${port}/StopListener/${listener.name}`, {
        method: "DELETE",
      });
    },
    {
      onSuccess(data, variables, context) {
        queryClient.invalidateQueries("listeners");
      },
    }
  );

  const createListener = useMutation(
    (listener: Listener) => {
      return fetch(`http://${ip}:${port}/Listener`, {
        method: "POST",
        body: JSON.stringify(listener),
      });
    },
    {
      onSuccess(data, variables, context) {
        queryClient.invalidateQueries("listeners");
        setDialogOpen(false);
      },
    }
  );

  const {
    register,
    handleSubmit,
    watch,
    reset,
    formState: { errors },
  } = useForm<Listener>();
  const onSubmit: SubmitHandler<Listener> = (data) =>
    createListener.mutate(data, { onSuccess: () => reset() });

  return (
    <>
      <Popup
        modal
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        contentStyle={{ backgroundColor: 'gray-800', color: 'text-white', width: '400px' }}
      >
        <div className="w-full bg-gray-800 text-white rounded p-4">
          <form
            className="shadow-md rounded p-4"
            onSubmit={handleSubmit(onSubmit)}
          >

            <div className="mb-4">
              <label className="block text-white-700 text-sm font-bold mb-2">
                Name
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                {...register("name", { required: true })}
                type="text"
                placeholder="Name"
              />
            </div>
            <div className="mb-6">
              <label className="block text-white-700 text-sm font-bold mb-2">
                Port
              </label>
              <input
                className="shadow appearance-none  border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                {...register("bindPort", { required: true })}
                placeholder="Port"
                type="number"
              />
            </div>
            <div className="flex items-center justify-between">
              <button
                className="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                type="submit"
              >
                Confirm
              </button>
            </div>
          </form>
        </div>
      </Popup>
      <main className="container mx-auto">
        <div className="w-full bg-gray-700 rounded border border-gray-300 p-4 min-h-[100px] min-w-[200px] overflow-auto">
          {(!listeners || listeners.length === 0) ? (
            <>
              <div className="text-center text-white-500">No listeners active</div>
              <div className="text-center mt-4">
                <button
                  className="text-white w-8 bg-blue-500 rounded"
                  onClick={() => {
                    setDialogOpen(true);
                  }}
                >
                  Add
                </button>
              </div>
            </>
          ) : (
            <table className="w-full text-sm text-left text-white-500 dark:text-gray-400">
              <thead className="text-md text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                <tr>
                  <th className="text-md text-white bg-gray-800">Name</th>
                  <th className="text-md text-white bg-gray-800">Port</th>
                  <th className="text-md text-white bg-gray-800"></th>
                </tr>
              </thead>
              <tbody>
                {listeners?.map((listener) => (
                  <tr key={listener.name}>
                    <td>{listener.name}</td>
                    <td>{listener.bindPort}</td>
                    <td>
                      <button
                        className="text-red-600 w-8 bg-white box-border border-2 border-white rounded"
                        onClick={() => deleteListener.mutate(listener)}
                      >
                        X
                      </button>
                    </td>
                  </tr>
                ))}
                <tr>
                  <td></td>
                  <td>
                    <ul>
                      <li>
                        <button
                          className="text-white w-8 bg-blue-500 rounded"
                          onClick={() => {
                            setDialogOpen(true);
                          }}
                        >
                          Add
                        </button>
                      </li>
                    </ul>
                  </td>
                  <td></td>
                </tr>
              </tbody>
            </table>
          )}
        </div>
      </main>
    </>
  );
}
