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

export function Listeners() {
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);

  const {
    data: listeners,
    isLoading,
    error,
  } = useQuery("listeners", () => {
    return fetch("http://localhost:8080/Listeners")
      .then((res) => res.json())
      .then((json) => z.array(ListenerSchema).parse(json));
  });

  const deleteListener = useMutation(
    (listener: Listener) => {
      return fetch(`http://localhost:8080/StopListener/${listener.name}`, {
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
      return fetch("http://localhost:8080/Listener", {
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
      <Popup modal open={dialogOpen} onClose={() => setDialogOpen(false)}>
        <div className="w-full">
          <form
            className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
            onSubmit={handleSubmit(onSubmit)}
          >
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2">
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
              <label className="block text-gray-700 text-sm font-bold mb-2">
                Port
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                {...register("bindPort", { required: true })}
                id="password"
                placeholder="Port"
                type="number"
              />
              {/* <p className="text-red-500 text-xs italic">
                Please choose a password.
              </p> */}
            </div>
            <div className="flex items-center justify-between">
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
                type="submit"
              >
                Confirm
              </button>
              {/* <a
                className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
                href="#"
              >
                Forgot Password?
              </a> */}
            </div>
          </form>
        </div>
      </Popup>
      <main className="container mx-auto">
        <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
          <thead className="text-md text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th>Name</th>
              <th>Port</th>
              <th></th>
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
      </main>
    </>
  );
}
