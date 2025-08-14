import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./global.css";
import { HeroUIProvider, ToastProvider } from "@heroui/react";
import { Login, Dashboard, NotFound, Redirect } from "@/pages";
import { RouterProvider, createBrowserRouter } from "react-router";
import UserContextProvider from "@/context/UserContext";

const router = createBrowserRouter([
  {
    path: "/:short",
    element: <Redirect />,
    errorElement: <NotFound />,
  },
  {
    path: "/dashboard",
    element: <Dashboard />,
    errorElement: <NotFound />,
  },
  {
    path: "/login",
    element: <Login />,
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <UserContextProvider>
      <HeroUIProvider>
        <RouterProvider router={router} />
        <ToastProvider />
      </HeroUIProvider>
    </UserContextProvider>
  </StrictMode>
);
