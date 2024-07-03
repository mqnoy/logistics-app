import { RouterProvider } from "react-router-dom";
import router from "./router";
import "bulma/css/bulma.css";
import 'react-toastify/dist/ReactToastify.css';
import { Provider } from "react-redux";
import { store } from "./redux";
import { ToastContainer } from "react-toastify";

function App() {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
      <ToastContainer
        position="bottom-right"
        autoClose={1000}
        hideProgressBar
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
    </Provider>
  )
}

export default App
