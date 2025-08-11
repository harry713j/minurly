import { useEffect, useRef, useState } from "react";
import { SERVER_URL } from "./utils/constants";

export default function App() {
  const [inputUrl, setInputUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [responseUrl, setResponseUrl] = useState("");
  const [isCopied, setIsCopied] = useState(false);
  const shortUrlRef = useRef<HTMLInputElement>(null);
  const [isLoadingRedirect, setIsLoadingRedirect] = useState(false);

  useEffect(() => {
    const path = window.location.pathname.slice(1);

    if (path.length === 8) {
      const fetchOriginalUrl = async () => {
        setIsLoadingRedirect(true);
        try {
          const response = await fetch(`${SERVER_URL}/${path}`);
          const data = await response.json();

          if (response.ok) {
            window.location.href = data.originalUrl;
          } else {
            setErrorMessage(data.error || "Inavalid short url");
          }
        } catch (error) {
          setErrorMessage("Failed to resolve the short url");
        } finally {
          setIsLoadingRedirect(false);
        }
      };

      fetchOriginalUrl();
    }
  }, []);

  const onClickShort = async () => {
    setIsLoading(true);
    setErrorMessage("");
    setResponseUrl("");
    try {
      console.log("server -> ", SERVER_URL);
      const response = await fetch(`${SERVER_URL}/short`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ originalUrl: inputUrl }),
      });
      const data = await response.json();

      console.log("Data: ", data);

      if (response.status === 201) {
        const protocol = window.location.protocol;
        const host = window.location.host;

        const shortUrl = `${protocol}//${host}/${data.shortCode}`;

        setResponseUrl(shortUrl);
      } else {
        setErrorMessage(data.error);
      }
    } catch (error) {
      console.error("Error at creating short url ", error);
      setErrorMessage((error as Error).message);
    } finally {
      setIsLoading(false);
      setInputUrl("");
    }
  };

  const copyToclipboard = () => {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(responseUrl);
      setIsCopied(true);
      setTimeout(() => {
        setIsCopied(false);
      }, 2000);
    } else {
      shortUrlRef.current?.select();
      document.execCommand("copy");
      setIsCopied(true);
      setTimeout(() => {
        setIsCopied(false);
      }, 2000);
    }
  };

  return isLoadingRedirect ? (
    <div className="w-full h-screen flex justify-center">
      <div className="w-20 h-20 rounded-full border-t-4 border-t-orange-400 animate-spin"></div>
    </div>
  ) : (
    <div className="w-full h-screen flex justify-center">
      <div className="w-1/2 h-1/3 px-20 py-12 mt-[12rem] flex flex-col items-start space-y-6 border-3 border-black rounded-md shadow-xl ">
        <h1 className="text-4xl text-orange-500 font-semibold">MinURly</h1>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            onClickShort();
          }}
          className="w-full flex items-center space-x-2"
        >
          <input
            type="text"
            placeholder="Enter your long url here"
            value={inputUrl}
            disabled={isLoading}
            onChange={(e) => setInputUrl(e.target.value)}
            className="w-3/5 px-4 py-2 text-slate-800 rounded border-none outline outline-slate-400 transition delay-300 focus:outline-orange-500 placeholder:italic placeholder:opacity-70 "
          />
          <button
            type="submit"
            disabled={isLoading}
            className="bg-orange-400 px-5 py-2 delay-200 transition font-semibold hover:bg-orange-500 rounded text-white cursor-pointer"
          >
            {isLoading ? "Please wait" : "Short"}
          </button>
        </form>

        {errorMessage && <p className="text-red-500">{errorMessage}</p>}
        {responseUrl && (
          <div className="flex items-center space-x-2 mt-10">
            <input
              ref={shortUrlRef}
              value={responseUrl}
              type="text"
              disabled
              className="bg-gray-200 rounded px-4 py-2  "
            />
            <button
              onClick={copyToclipboard}
              className="bg-green-500 rounded px-5 py-2 text-white cursor-pointer"
            >
              {isCopied ? "Copied" : "Copy"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
