import { SERVER_URL } from "@/utils/constants";
import axios, { AxiosError } from "axios";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import type { UserContextType, User } from "@/types/types";

const UserContext = createContext<UserContextType>({
  user: null,
  isFetching: true,
  error: "",
  setUser: () => {},
  refetchUser: async () => {},
});

export default function UserContextProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState("");

  const getUser = useCallback(async () => {
    setError("");
    try {
      const response = await axios.get(`${SERVER_URL}/users`, {
        withCredentials: true,
      });

      if (response.status === 200) {
        setUser(response.data.user);
      } else {
        setError("Failed to get user details");
      }
    } catch (error) {
      setError(((error as AxiosError).response?.data as any).error);
    } finally {
      setIsLoading(false);
    }
  }, [setUser, user]);

  useEffect(() => {
    getUser();
  }, []);

  return (
    <UserContext.Provider
      value={{
        user: user,
        setUser: setUser,
        isFetching: isLoading,
        error,
        refetchUser: getUser,
      }}
    >
      {children}
    </UserContext.Provider>
  );
}

export const useUser = () => useContext(UserContext);
