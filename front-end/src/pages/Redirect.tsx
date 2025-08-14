import { SERVER_URL } from "@/utils/constants";
import { addToast } from "@heroui/react";
import axios, { type AxiosError } from "axios";
import { CircleNotch } from "phosphor-react";
import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router";

export default function Redirect() {
  const [isLoading, setIsLoading] = useState(true);
  const params = useParams();

  const redirect = useCallback(async () => {
    try {
      const response = await axios.get(`${SERVER_URL}/urls/${params?.short}`);

      if (response.status === 200) {
        window.location.href = response.data.originalUrl;
      }
    } catch (error) {
      const axiosError = error as AxiosError;
      addToast({
        title: "URL Not found",
        description: (axiosError.response?.data as any).error,
        color: "danger",
      });
    } finally {
      setIsLoading(false);
    }
  }, [params]);

  useEffect(() => {
    redirect();
  }, [redirect]);

  return (
    <>
      {isLoading && (
        <div>
          <CircleNotch size={48} className="animate-spin" />
        </div>
      )}
    </>
  );
}
