import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./global.css";
import { HeroUIProvider } from "@heroui/react";
import { Login } from "@/pages";
import { RouterProvider, createBrowserRouter } from "react-router";

const router = createBrowserRouter([
  {
    path: "/",
    element: <h1>Hello there</h1>,
    errorElement: <h1>Not found</h1>,
  },
  {
    path: "/login",
    element: <Login />,
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <HeroUIProvider>
      <RouterProvider router={router} />
    </HeroUIProvider>
  </StrictMode>
);
