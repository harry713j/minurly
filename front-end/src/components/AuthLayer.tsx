import { SERVER_URL } from "@/utils/constants";
import axios from "axios";
import { CircleNotch } from "phosphor-react";
import type React from "react";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";

export function AuthLayer({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const navigate = useNavigate();
  const [haveSession, setHaveSession] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    setIsLoading(true);
    axios
      .get(`${SERVER_URL}/auth/status`, { withCredentials: true })
      .then((response) => {
        setHaveSession(response.data.isAuthorized);
      })
      .catch((err) => {
        setHaveSession(false);
        console.error("Auth error: ", err);
      })
      .finally(() => setIsLoading(false));
  }, [location.pathname]);

  useEffect(() => {
    if (location.pathname === "/") {
      navigate(haveSession ? "/dashboard" : "/login");
    } else if (!haveSession && location.pathname.startsWith("/dashboard")) {
      navigate("/login");
    } else if (haveSession && location.pathname.startsWith("/login")) {
      navigate("/dashboard");
    }
  }, [location.pathname, haveSession, navigate, isLoading]);

  if (isLoading) {
    return (
      <div className="w-full h-screen flex justify-center items-center">
        <CircleNotch size={48} className="animate-spin" />
      </div>
    );
  }

  return <>{children}</>;
}
